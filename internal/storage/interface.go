package storage

import (
	"context"
	"time"

	"github.com/lilythecat859/rpcv2-hist/internal/model"
)

// HistoricalStore is the minimal pluggable interface any backend must satisfy.
type HistoricalStore interface {
	// Ping returns nil if the backend is reachable.
	Ping(context.Context) error
	// Close releases resources.
	Close() error

	// Block methods
	GetBlock(ctx context.Context, slot uint64, commitment Commitment) (*model.Block, error)
	GetBlocksWithLimit(ctx context.Context, start, limit uint64, commitment Commitment) ([]uint64, error)
	GetBlockTime(ctx context.Context, slot uint64) (*time.Time, error)

	// Transaction methods
	GetTransaction(ctx context.Context, signature string, commitment Commitment) (*model.Transaction, error)

	// Signature methods
	GetSignaturesForAddress(ctx context.Context, addr string, opts SignatureOpts) ([]model.SignatureInfo, error)
}

// Commitment level alias to avoid importing Solana SDK here.
type Commitment string

const (
	CommitmentProcessed Commitment = "processed"
	CommitmentConfirmed Commitment = "confirmed"
	CommitmentFinalized Commitment = "finalized"
)

// SignatureOpts bundles pagination and filtering.
type SignatureOpts struct {
	Limit      uint64
	Before     *string // signature
	Until      *string // signature
	Commitment Commitment
}

// StoreKind identifies the driver for factory usage.
type StoreKind string

const (
	StoreClickHouse StoreKind = "clickhouse"
	StorePostgres   StoreKind = "postgres"
	StoreParquet    StoreKind = "parquet"
)