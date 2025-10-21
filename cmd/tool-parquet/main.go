//go:build !cgo
// +build !cgo

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"

	"github.com/lilythecat859/rpcv2-hist/internal/parquet"
)

var (
	inDir  = flag.String("in", "", "input directory with JSON block files")
	outDir = flag.String("out", "", "output directory for .parquet")
	day    = flag.String("day", "", "YYYY-MM-DD to process")
)

func main() {
	flag.Parse()
	if *inDir == "" || *outDir == "" || *day == "" {
		flag.Usage()
		os.Exit(1)
	}
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	t, err := time.Parse("2006-01-02", *day)
	if err != nil {
		log.Fatalf("bad day: %v", err)
	}
	start := t.Unix()
	end := t.Add(24 * time.Hour).Unix()

	pattern := filepath.Join(*inDir, fmt.Sprintf("%d_*.json", start))
	files, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatalf("glob: %v", err)
	}
	if len(files) == 0 {
		log.Fatalf("no files match %s", pattern)
	}

var blocks []model.Block
	for _, f := range files {
		if err := loadBlocks(f, &blocks); err != nil {
			logger.Error("load", zap.String("file", f), zap.Error(err))
		}
	}
	out := filepath.Join(*outDir, fmt.Sprintf("blocks_%s.parquet", *day))
	w := parquet.NewWriter(out, parquet.WithLogger(logger))
	if err := w.WriteBlocks(blocks); err != nil {
		log.Fatalf("write: %v", err)
	}
	fmt.Println(out)
}

func loadBlocks(path string, dst *[]model.Block) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(dst)
}