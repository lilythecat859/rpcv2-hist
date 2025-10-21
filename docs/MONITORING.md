# Monitoring & Alerting

## Metrics Endpoint
Prometheus scrape at `:9091/metrics`

## Key Alerts
- `rpcv2_hist_request_duration_seconds` P99 > 200 ms
- `rpcv2_hist_requests_total` error rate > 1 %
- ClickHouse disk > 85 %

## Dashboards
Import Grafana JSON from `monitoring/grafana.json`

## Tracing
OTLP gRPC on `:4317` → Jaeger

## Logs
Structured JSON to stdout; ship with Vector/Loki