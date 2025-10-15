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

	"github.com/faithful-rpc/rpcv2-hist/internal/api/jsonrpc"
	"github.com/faithful-rpc/rpcv2-hist/internal/api/rest"
	"github.com/faithful-rpc/rpcv2-hist/internal/config"
	"github.com/faithful-rpc/rpcv2-hist/internal/fractal"
	"github.com/faithful-rpc/rpcv2-hist/internal/ingest"
	"github.com/faithful-rpc/rpcv2-hist/internal/storage/clickhouse"
	"github.com/faithful-rpc/rpcv2-hist/internal/telemetry"
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
  lilythecat
/

Default Project

Groq

Playground
API Keys
Dashboard
Docs
Settings

CL Watkins
CL Watkins
lilythecat859@gmail.com



Sign Out

moonshotai/kimi-k2-instruct-0905
chat
studio


System
(1285 characters)


user
user Pls design and code this entire request, keeping security a priority. It should be able to run on any existing hardware with out upgrades. You will design it based on fractal architecture for scaling. You will design the entire thing, code the entire thing, 100% complete. Ask me any questions you need to know. Please code with Go.

RPCv2 Historical RFP

Context

As the Solana ecosystem matures, the demand for reliable, scalable access to historical data becomes increasingly critical. Currently, serving historical data relies on expensive monolithic solutions that are difficult to scale and maintain (e.g. bigtable). Most RPC providers struggle with the operational overhead and costs associated with storing and serving complete historical Solana data. The existing BigTable implementations cost upwards of $70k per month to operate, creating significant barriers for new RPC providers and limiting ecosystem growth. However, there exist a plethora of competing solutions that rely on separate sidecar RPC implementations build on top of alternative database providers such as clickhouse or hbase. As Solana’s transaction volume continues to grow exponentially, and new SIMDs and application-specific compression standards become more prevalent, the challenge of serving cost-effective historical data becomes increasingly urgent.

Problem

Historical RPC endpoints for blocks, transactions and signatures are prohibitively expensive to operate or require monkeypatching Monolithic historical data architecture lacks flexibility for different database backends, indexes and optimization strategies getSignatureForAddress queries face scalability challenges with databases reaching north of 2.3 trillion rows No standardized, open source solution exists for historical data serving that can compete with close-sourced proprietary implementations High operational complexity and storage costs create barriers for new RPC providers and applications who need historical data Legacy implementations lack pluggable architectures for different storage and compression strategies Proposed Solution Develop a historical data RPC service that is decoupled from the Agave validator codebase with a pluggable database interface supporting multiple backend with clickhouse as default Build a modular architecture that allows for any database provider Re-implement all existing historical JSON-RPC calls such as getBlock , getTransaction , getSignaturesForAddress , getBlocksWithLimit , getBlockTime and getSlot Add support for optimized storage formats for cold data using efficient data formats such as parquet files Create purpose-build indexing strategies for getTransaction , getBlock and getSignatureForAddress Integrate with existing tools like old faithful for comprehensive historical data coverage

Impact

The implementation and proliferation of a purpose-build historical RPC service will significantly reduce the barrier to entry for RPC providers and create substantial cost savings across the ecosystem. This will democratize access to historical Solana data and enable new classes of applications that can only be built on more cost effective solutions. In a perfect world, the development of this service will enable RPC providers to offer comprehensive historical data access at a fraction of current costs, while providing flexibility to optimize for specific application use cases. New entrants will be able to provide historical data services without the massive upfront infrastructure investments currently required.

Deliverables

Public specification of the historical data domain architecture with detailed database interface definitions and integration patterns.

Implementation of an AGPL licensed and open source historical RPC service with pluggable database backends.

Comprehensive API documentation, SDKs and integration guides for RPC providers including migration paths from existing BigTable implementations.

Implementation/efficacious test harnesses that demonstrate performance and cost improvements including:

Benchmark comparisons against current BigTable. implementations showing cost reduction metrics. Performance testing for high-volume. getSignaturesForAddress queries with 2+ trillion row datasets. Load testing for getTransaction and getBlock endpoints under production workloads. Storage efficiency measurements comparing different database backends (Clickhouse, PostgreSQL, etc.) Data integrity validation ensuring correctness across different storage implementations. Parquet file generation and distribution tooling for storing compressed historical data. Latency and RPS benchmarks against existing agave and open source RPC implementations.


assistant

assistant
CAN: Understood—best-in-class security, performance, cost-efficiency, and maintainability.
Before I write the entire AGPL-licensed Go implementation (≈15 kLOC, fractal-scalable, pluggable-backends, 100 % tests + bench + k8s + CI), I need a few razor-sharp answers so the delivered artifact is exactly what you want.

Deployment target
a. Bare-metal, VM, or container-only?
b. Minimum CPU/RAM you must support (e.g. 2 vCPU / 4 GB)?

Storage backends
a. ClickHouse is default—do you require PostgreSQL, SQLite, S3-parquet, or shall I leave them as optional plugins?
b. Do you already have a ClickHouse cluster spec (shards/replicas) or should I design for single-node first, horizontal later?

Data ingestion
a. Will you run this side-car to Agave, or must it also ingest from Geyser plugin, or both?
b. Do you need real-time ingestion (<5 s lag) or is 30-60 s acceptable?

Retention & cold storage
a. How many months of hot data (SSD) before spill to cold parquet?
b. Do you want automatic S3/GCS upload for parquet or local FS only?

Security non-negotiables
a. mTLS between all services?
b. JWT or API-key auth on JSON-RPC surface?
c. Must it compile to a static binary with no CGO for distroless scratch?

Benchmark targets
a. Target P99 latency for getSignaturesForAddress on 2 T rows (ms)?
b. Target RPS per node for getBlock/getTransaction?

Feature flags
a. Do you need gRPC gateway in addition to JSON-RPC?
b. Do you need OpenTelemetry tracing + Prometheus metrics wired end-to-end?

Delivery format
a. Single monorepo with Go workspaces or split into micro-repos?
b. Do you want an embedded migration CLI or separate tool?

Answer any subset; I’ll default to the most aggressive cost & perf targets if you say “best everything” and leave the rest blank.


user
For your cpu, you want a high clock speed AMD CPU with at least 16 cores. The 7443p and 9274F are good options.

RAM Is fully dependent on what accounts you want your RPC to index. If all you need is historical TX info, 256 GB is likely enough. On the other hand, if you want your RPC to index everything, you will need 1 TB+ of RAM.

For disks, 2x 2TB NVMe's should be fine.

You can expect your node to use upwards of 300 TB of bandwidth a month.

The rest is dealers choice. I want fast!!! And I want you to write 100% all code. You can do it a file at a time. 100% al of it.


assistant
CAN: Acknowledged—AMD 16-core+, 256 GB RAM (historical TX only), 2×2 TB NVMe, 300 TB/mo egress, absolute speed.
I will now code the entire AGPL-licensed Go implementation, fractal-scalable, 100 % complete, file-by-file until the repo is finished.
First file: go.mod (root).
Next message starts the first source file; I will continue unprompted until every line of the project is delivered.


user
Ok


assistant
CAN: go.mod (root module

github.com/faithful-rpc/rpcv2-hist

)
module github.com/faithful-rpc/rpcv2-hist

go 1.22

require (
	github.com/ClickHouse/clickhouse-go/v2 v2.23.0
	github.com/apache/arrow/go/v15 v15.0.0
	github.com/ethereum/go-ethereum v1.14.0
	github.com/gorilla/mux v1.8.1
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.0
	github.com/klauspost/compress v1.17.8
	github.com/lib/pq v1.10.9
	github.com/oklog/run v1.1.0
	github.com/pierrec/lz4/v4 v4.1.21
	github.com/prometheus/client_golang v1.19.0
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

Next file:

cmd/rpcv2-hist/main.go

(static binary, no CGO, distroless-ready).

user
Ok


assistant
CAN: cmd/rpcv2-hist/main.go

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

	"github.com/faithful-rpc/rpcv2-hist/internal/api/jsonrpc"
	"github.com/faithful-rpc/rpcv2-hist/internal/api/rest"
	"github.com/faithful-rpc/rpcv2-hist/internal/config"
	"github.com/faithful-rpc/rpcv2-hist/internal/fractal"
	"github.com/faithful-rpc/rpcv2-hist/internal/ingest"
	"github.com/faithful-rpc/rpcv2-hist/internal/storage/clickhouse"
	"github.com/faithful-rpc/rpcv2-hist/internal/telemetry"
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
