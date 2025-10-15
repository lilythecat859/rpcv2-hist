//go:build !cgo
// +build !cgo

package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lilythecat859/rpcv2-hist/internal/api/jsonrpc"
	"github.com/lilythecat859/rpcv2-hist/internal/api/rest"
	"github.com/lilythecat859/rpcv2-hist/internal/config"
	"github.com/lilythecat859/rpcv2-hist/internal/fractal"
	"github.com/lilythecat859/rpcv2-hist/internal/ingest"
	"github.com/lilythecat859/rpcv2-hist/internal/storage/clickhouse"
	"github.com/lilythecat859/rpcv2-hist/internal/telemetry"
	"github.com/oklog/run"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	if err := runMain(); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}
}

func runMain() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logger, err := telemetry.NewLogger(cfg.LogLevel)
	if err != nil {
		return fmt.Errorf("new logger: %w", err)
	}
	defer func() { _ = logger.Sync() }()

	tp, err := telemetry.NewTraceProvider(ctx, cfg)
	if err != nil {
		return fmt.Errorf("new trace provider: %w", err)
	}
	defer func() { _ = tp.Shutdown(ctx) }()
	otel.SetTracerProvider(tp)

	mp := telemetry.NewMetricProvider(cfg)
	defer func() { _ = mp.Shutdown(ctx) }()

	db, err := clickhouse.New(ctx, cfg.ClickHouse, clickhouse.WithLogger(logger))
	if err != nil {
		return fmt.Errorf("open clickhouse: %w", err)
	}
	defer func() { _ = db.Close() }()

	ing, err := ingest.New(db, ingest.WithLogger(logger))
	if err != nil {
		return fmt.Errorf("new ingest: %w", err)
	}
  
	github.com/spf13/cobra v1.8.0
	github.com/spf13/viper v1.18.2
	github.com/stretchr/testify v1.9.0
	github.com/tidwall/gjson v1.17.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.49.0
	go.opentelemetry.io/otel v1.25.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.25.0
	go.opentelemetry.io/otel/sdk v1.25.0
	go.uber.org/zap v1.27.0
	golang.org/x/sync v0.7.0
	google.golang.org/grpc v1.63.2
	google.golang.org/protobuf v1.33.0
	gopkg.in/yaml.v3 v3.0.1
)
