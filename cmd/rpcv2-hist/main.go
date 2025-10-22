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

fractalRoot := fractal.NewRoot(db, logger)

	var g run.Group
	// JSON-RPC server
	{
		rpcSrv := jsonrpc.NewServer(fractalRoot, logger)
		mux := http.NewServeMux()
		mux.Handle("/", rpcSrv)
		srv := &http.Server{
			Addr:    cfg.JSONRPCListen,
			Handler: telemetry.HTTPMiddleware(mux),
		}
		g.Add(func() error {
			logger.Info("starting json-rpc", zap.String("addr", cfg.JSONRPCListen))
			return srv.ListenAndServe()
		}, func(err error) {
			shutdownCtx, done := context.WithTimeout(context.Background(), 5*time.Second)
			_ = srv.Shutdown(shutdownCtx)
			done()
		})
	}
	// REST gateway
	{
		restSrv := rest.NewServer(fractalRoot, logger)
		srv := &http.Server{
			Addr:    cfg.RESTListen,
			Handler: telemetry.HTTPMiddleware(restSrv),
		}
		g.Add(func() error {
			logger.Info("starting rest", zap.String("addr", cfg.RESTListen))
			return srv.ListenAndServe()
		}, func(err error) {
			shutdownCtx, done := context.WithTimeout(context.Background(), 5*time.Second)
			_ = srv.Shutdown(shutdownCtx)
			done()
		})
	}
	// gRPC server
	{
		ln, err := net.Listen("tcp", cfg.GRPCListen)
		if err != nil {
			return fmt.Errorf("grpc listen: %w", err)
		}
		srv := grpc.NewServer(
			grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
			grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
		)
		// register services here
		g.Add(func() error {
			logger.Info("starting grpc", zap.String("addr", cfg.GRPCListen))
			return srv.Serve(ln)
		}, func(err error) {
			srv.GracefulStop()
		})
	}
	// Ingester
	{
		g.Add(func() error {
			return ing.Run(ctx)
		}, func(err error) {
			cancel()
		})
	}
	// Signal handler
	{
		g.Add(func() error {
			<-ctx.Done()
			return errors.New("terminated")
		}, func(err error) {
			cancel()
		})
	}

	return g.Run()
}