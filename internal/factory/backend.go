package factory

import (
	"context"
	"fmt"

	"github.com/faithful-rpc/rpcv2-hist/internal/storage"
	"github.com/faithful-rpc/rpcv2-hist/internal/storage/clickhouse"
)

func NewBackend(ctx context.Context, kind storage.StoreKind, cfg any) (storage.HistoricalStore, error) {
	switch kind {
	case storage.StoreClickHouse:
		c, ok := cfg.(clickhouse.Config)
		if !ok {
			return nil, fmt.Errorf("invalid clickhouse config")
		}
		return clickhouse.New(ctx, c)
	default:
		return nil, fmt.Errorf("unknown backend %q", kind)
	}
}