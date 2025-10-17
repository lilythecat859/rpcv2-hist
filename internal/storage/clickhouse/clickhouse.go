package clickhouse

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"go.uber.org/zap"

	"github.com/faithful-rpc/rpcv2-hist/internal/model"
	"github.com/faithful-rpc/rpcv2-hist/internal/storage"
)

type DB struct {
	conn   driver.Conn
	logger *zap.Logger
	cfg    Config
}

type Config struct {
	Addr         string
	Database     string
	User         string
	Password     string
	AsyncInsert  bool
	MaxOpenConns int
	MaxIdleConns int
	ConnMaxLifetime time.Duration
}

type Option func(*DB)

func WithLogger(l *zap.Logger) Option {
	return func(d *DB) { d.logger = l }
}

func New(ctx context.Context, cfg Config, opts ...Option) (*DB, error) {
	db := &DB{cfg: cfg}
	for _, o := range opts {
		o(db)
	}
	optsCH := &clickhouse.Options{
		Addr: []string{cfg.Addr},
		Auth: clickhouse.Auth{
			Database: cfg.Database,
			Username: cfg.User,
			Password: cfg.Password,
		},
		Settings: clickhouse.Settings{
			"async_insert": cfg.AsyncInsert,
		},
		MaxOpenConns:    cfg.MaxOpenConns,
		MaxIdleConns:    cfg.MaxIdleConns,
		ConnMaxLifetime: cfg.ConnMaxLifetime,
	}
	conn, err := clickhouse.Open(optsCH)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}
	db.conn = conn
	return db, nil
}

func (d *DB) Ping(ctx context.Context) error {
	return d.conn.Ping(ctx)
}

func (d *DB) Close() error {
	return d.conn.Close()
}

func (d *DB) GetBlock(ctx context.Context, slot uint64, commitment storage.Commitment) (*model.Block, error) {
	row := d.conn.QueryRow(ctx, `
		SELECT slot, blockhash, parent_slot, block_time, height, raw
		FROM blocks
		WHERE slot = ? AND commitment = ?
	`, slot, string(commitment))
	var b model.Block
	if err := row.Scan(&b.Slot, &b.Blockhash, &b.ParentSlot, &b.BlockTime, &b.Height, &b.Raw); err != nil {
		return nil, fmt.Errorf("scan block: %w", err)
	}
	return &b, nil
}

func (d *DB) GetBlocksWithLimit(ctx context.Context, start, limit uint64, commitment storage.Commitment) ([]uint64, error) {
	rows, err := d.conn.Query(ctx, `
		SELECT slot
		FROM blocks
		WHERE slot >= ? AND commitment = ?
		ORDER BY slot
		LIMIT ?
	`, start, string(commitment), limit)
	if err != nil {
		return nil, fmt.Errorf("query blocks: %w", err)
	}
	defer rows.Close()
	var slots []uint64
	for rows.Next() {
		var s uint64
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		slots = append(slots, s)
	}
	return slots, rows.Err()
}

func (d *DB) GetBlockTime(ctx context.Context, slot uint64) (*time.Time, error) {
	row := d.conn.QueryRow(ctx, `SELECT block_time FROM blocks WHERE slot = ?`, slot)
	var t int64
	if err := row.Scan(&t); err != nil {
		return nil, err
	}
	pt := time.Unix(t, 0)
	return &pt, nil
}

func (d *DB) GetTransaction(ctx context.Context, signature string, commitment storage.Commitment) (*model.Transaction, error) {
	row := d.conn.QueryRow(ctx, `
		SELECT signature, slot, tx_idx, block_time, signer, fee, compute_units, err, raw
		FROM transactions
		WHERE signature = ? AND commitment = ?
	`, signature, string(commitment))
	var tx model.Transaction
	if err := row.Scan(&tx.Signature, &tx.Slot, &tx.Index, &tx.BlockTime, &tx.Signer, &tx.Fee, &tx.ComputeUnits, &tx.Err, &tx.Raw); err != nil {
		return nil, fmt.Errorf("scan tx: %w", err)
	}
	return &tx, nil
}

func (d *DB) GetSignaturesForAddress(ctx context.Context, addr string, opts storage.SignatureOpts) ([]model.SignatureInfo, error) {
	q := `
		SELECT signature, slot, err, memo, block_time
		FROM signatures
		WHERE address = ? AND commitment = ?
	`
	args := []interface{}{addr, string(opts.Commitment)}
	if opts.Before != nil {
		q += ` AND slot < (SELECT slot FROM signatures WHERE signature = ? LIMIT 1)`
		args = append(args, *opts.Before)
	}
	if opts.Until != nil {
		q += ` AND slot > (SELECT slot FROM signatures WHERE signature = ? LIMIT 1)`
		args = append(args, *opts.Until)
	}
	q += ` ORDER BY slot DESC LIMIT ?`
	args = append(args, opts.Limit)

	rows, err := d.conn.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("query sigs: %w", err)
	}
	defer rows.Close()
	var out []model.SignatureInfo
	for rows.Next() {
		var si model.SignatureInfo
		var bt int64
		if err := rows.Scan(&si.Signature, &si.Slot, &si.Err, &si.Memo, &bt); err != nil {
			return nil, err
		}
		si.BlockTime = time.Unix(bt, 0)
		out = append(out, si)
	}
	return out, rows.Err()
}