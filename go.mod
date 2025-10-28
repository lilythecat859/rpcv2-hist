module github.com/lilythecat859/rpcv2-hist

go 1.22

require (
	// ---------- Database ----------
	github.com/ClickHouse/clickhouse-go/v2 v2.23.0
	github.com/lib/pq v1.10.9

	// ---------- Data formats ----------
	github.com/apache/arrow/go/v15 v15.0.0
	github.com/klauspost/compress v1.17.8
	github.com/pierrec/lz4/v4 v4.1.21

	// ---------- Ethereum ----------
	github.com/ethereum/go-ethereum v1.14.0

	// ---------- HTTP / RPC ----------
	github.com/gorilla/mux v1.8.1
	google.golang.org/grpc v1.63.2
	google.golang.org/protobuf v1.33.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.0

	// ---------- OpenTelemetry ----------
	go.opentelemetry.io/otel v1.25.0
	go.opentelemetry.io/otel/sdk v1.25.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.49.0
	// **THIS IS THE CORRECT LINE – no colon, just space between path and version**
	go.opentelemetry.io/otel/semconv v1.24.0

	// ---------- Logging ----------
	go.uber.org/zap v1.27.0

	// ---------- Metrics / Monitoring ----------
	github.com/prometheus/client_golang v1.19.0

	// ---------- CLI / Config ----------
	github.com/spf13/cobra v1.8.0
	github.com/spf13/viper v1.18.2

	// ---------- Misc utilities ----------
	github.com/oklog/run v1.1.0
	github.com/stretchr/testify v1.9.0
	github.com/tidwall/gjson v1.17.0
	golang.org/x/sync v0.7.0
	gopkg.in/yaml.v3 v3.0.1
)
