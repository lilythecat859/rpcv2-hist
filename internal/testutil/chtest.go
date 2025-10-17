package testutil

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/ClickHouse/clickhouse-go/v2"
)

// ClickHouseTestContainer starts a disposable ClickHouse for integration tests.
func ClickHouseTestContainer(t testing.TB) clickhouse.Conn {
	t.Helper()
	addr := os.Getenv("CLICKHOUSE_TEST_ADDR")
	if addr == "" {
		t.Skip("CLICKHOUSE_TEST_ADDR not set")
	}
	opts := &clickhouse.Options{
		Addr: []string{addr},
		Auth: clickhouse.Auth{
			Database: "solana_test",
			Username: "default",
		},
	}
	conn, err := clickhouse.Open(opts)
	if err != nil {
		t.Fatalf("open test clickhouse: %v", err)
	}
	if err := conn.Ping(context.Background()); err != nil {
		t.Fatalf("ping test clickhouse: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	return conn
}