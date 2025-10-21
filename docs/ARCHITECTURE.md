# Architecture

## Fractal Scaling
Root → N shards → each shard is a full storage backend (ClickHouse, Postgres, Parquet).  
Shard selection via `slotFromSignature(signature) % shards`.  
Hot paths keep shards in-memory; cold paths spill to S3-parquet.

## Security
- mTLS between every service (opt-out via flag)
- JWT or API-key on JSON-RPC surface
- Static binary with no CGO → distroless scratch
- 90-day responsible disclosure policy

## Storage
- ClickHouse default: 30-day hot NVMe, then TTL to S3-parquet
- Partitioned by `intDiv(slot, 864000)` (~100k slots)
- Replacing-merge for `signatures_latest` view

## Cost
BigTable: ~70k USD/mo  
rpcv2-hist: ~4k USD/mo (same data, 3×16-core AMD, 256 GB RAM, 2×2 TB NVMe, 300 TB egress)

## Benchmarks
See `internal/bench/bench_test.go` – 2 T-row `getSignaturesForAddress` 