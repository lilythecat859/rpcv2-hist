package bench

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/faithful-rpc/rpcv2-hist/internal/model"
	"github.com/faithful-rpc/rpcv2-hist/internal/storage"
	"github.com/faithful-rpc/rpcv2-hist/internal/storage/clickhouse"
)

func getStore(tb testing.TB) storage.HistoricalStore {
	cfg := clickhouse.Config{
		Addr:         getenv("CH_ADDR", "127.0.0.1:9000"),
		Database:     getenv("CH_DB", "solana"),
		User:         "default",
		Password:     "",
		AsyncInsert:  true,
		MaxOpenConns: 64,
	}
	db, err := clickhouse.New(context.Background(), cfg)
	require.NoError(tb, err)
	tb.Cleanup(func() { _ = db.Close() })
	return db
}

func BenchmarkGetSignaturesForAddress(b *testing.B) {
	ctx := context.Background()
	store := getStore(b)
	addr := "11111111111111111111111111111111"
	opts := storage.SignatureOpts{
		Limit:      1000,
		Commitment: storage.CommitmentConfirmed,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := store.GetSignaturesForAddress(ctx, addr, opts)
		require.NoError(b, err)
	}
}

func BenchmarkGetBlock(b *testing.B) {
	ctx := context.Background()
	store := getStore(b)
	slot := uint64(rand.Int63n(200_000_000))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := store.GetBlock(ctx, slot, storage.CommitmentFinalized)
		require.NoError(b, err)
	}
}

func BenchmarkGetTransaction(b *testing.B) {
	ctx := context.Background()
	store := getStore(b)
	sig := generateRandomSig()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := store.GetTransaction(ctx, sig, storage.CommitmentFinalized)
		require.NoError(b, err)
	}
}

func generateRandomSig() string {
	b := make([]byte, 64)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}