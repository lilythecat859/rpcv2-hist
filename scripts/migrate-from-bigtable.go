//go:build ignore
// +build ignore

// migrate-from-bigtable streams rows into ClickHouse
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/bigtable"
	"github.com/ClickHouse/clickhouse-go/v2"
)

var (
	btProject  = flag.String("bt-project", "", "BigTable project")
	btInstance = flag.String("bt-instance", "", "BigTable instance")
	btTable    = flag.String("bt-table", "", "BigTable table")
	chAddr     = flag.String("ch-addr", "127.0.0.1:9000", "ClickHouse addr")
	chDB       = flag.String("ch-db", "solana", "ClickHouse database")
	workers    = flag.Int("workers", 16, "parallel workers")
)

func main() {
	flag.Parse()
	if *btProject == "" || *btInstance == "" || *btTable == "" {
		flag.Usage()
		os.Exit(1)
	}
	ctx := context.Background()

	btClient, err := bigtable.NewClient(ctx, *btProject, *btInstance)
	if err != nil {
		log.Fatalf("bt client: %v", err)
	}
	tbl := btClient.Open(*btTable)

	chConn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{*chAddr},
		Auth: clickhouse.Auth{Database: *chDB},
	})
	if err != nil {
		log.Fatalf("ch open: %v", err)
	}
defer chConn.Close()

	batch := make(chan bigtableRow, 10000)
	var wg sync.WaitGroup
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			inserter(chConn, batch)
		}()
	}

	tbl.ReadRows(ctx, bigtable.PrefixRange(""), func(r bigtable.Row) bool {
		batch <- bigtableRow{key: r.Key(), row: r}
		return true
	}, bigtable.RowFilter(bigtable.StripValueFilter()))
	close(batch)
	wg.Wait()
	fmt.Println("migration complete")
}

type bigtableRow struct {
	key string
	row bigtable.Row
}

func inserter(conn clickhouse.Conn, rows <-chan bigtableRow) {
	ctx := context.Background()
	stmt, err := conn.PrepareBatch(ctx, `
		INSERT INTO blocks (slot, blockhash, parent_slot, block_time, height, commitment, raw)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		log.Fatalf("prepare: %v", err)
	}
	for r := range rows {
		var blk block
		if err := json.Unmarshal(r.row["block"][0].Value, &blk); err != nil {
			log.Printf("skip bad row %s: %v", r.key, err)
			continue
		}
		if err := stmt.Append(
			blk.Slot,
			blk.Blockhash,
			blk.ParentSlot,
			blk.BlockTime,
			blk.Height,
			"finalized",
			r.row["block"][0].Value,
		); err != nil {
log.Printf("append: %v", err)
		}
		if stmt.Rows() >= 10000 {
			if err := stmt.Send(); err != nil {
				log.Printf("send: %v", err)
			}
		}
	}
	if err := stmt.Send(); err != nil {
		log.Printf("final send: %v", err)
	}
}

type block struct {
	Slot       uint64 `json:"slot"`
	Blockhash  string `json:"blockhash"`
	ParentSlot uint64 `json:"parentSlot"`
	BlockTime  int64  `json:"blockTime"`
	Height     uint64 `json:"blockHeight"`
}