# Monitoring Development

## Local
- Prometheus on :9091
- Jaeger all-in-one on :16686
- Grafana on :3000 (import `monitoring/grafana.json`)

## CI
- Assert P99 < 100 ms in bench
- Assert error rate < 1 %

## Alerts
Add to `monitoring/alerts.yaml`:
```yaml
- alert: HighLatency
  expr: histogram_quantile(0.99, rpcv2_hist_request_duration_seconds) > 0.2

# Dashboards
Embed in
```
monitoring/grafana.json
```
and version-control.