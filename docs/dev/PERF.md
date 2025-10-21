# Performance Tuning

## OS
Run `scripts/perf-tune.sh` on bare-metal.

## ClickHouse
Run `scripts/clickhouse-tuning.sh`.

## Go
- Use `fasthttp` if JSON-RPC becomes CPU-bound (future opt-in)
- Keep `CGO_ENABLED=0` for static binaries
- Use `sync.Pool` for hot buffers

## Benchmark
```bash
go test -bench=. -benchmem ./...

Profile
scripts/bench-pprof.sh