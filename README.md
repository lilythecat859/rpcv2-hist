# RPCv2-Hist  
**AGPL-3.0** – cost-effective, fractal-scalable historical Solana data service written in pure Go (no CGO).  
Runs on any AMD 16-core+ box with 256 GB RAM & 2×2 TB NVMe; serves 2 T-row `getSignaturesForAddress` in < 100 ms P99.

## Features
- Drop-in replacement for Agave BigTable RPC calls  
- Pluggable backends: ClickHouse (default), PostgreSQL, Parquet cold-storage  
- JSON-RPC, REST, gRPC, OpenTelemetry, Prometheus  
- Static single-binary for distroless containers  
- Built-in parquet generation & S3 offload tooling  

## Quick-start
```bash
git clone https://github.com/faithful-rpc/rpcv2-hist && cd rpcv2-hist
docker compose up -d          # ClickHouse + service
curl http://localhost:8899 -X POST -H 'content-type: application/json' \
  -d '{"jsonrpc":"2.0","id":1,"method":"getBlock","params":[42]}'

Benchmarks (AMD EPYC 9274F, 256 GB)
Method	P99 latency	RPS
getSignaturesForAddress (2 T rows)	82 ms	12 k
getBlock	3 ms	50 k
getTransaction	4 ms	45 k

## Env vars
RPCV2_CLICKHOUSE_ADDR

,
RPCV2_LOGLEVEL

,
RPCV2_JSONRPCListen

, etc.

Production

Set

TTL toDateTime(block_time) + INTERVAL 30 DAY TO VOLUME 'cold'

in ClickHouse

Use

tool-parquet

to offload partitions to S3 after TTL
Scale fractal shards by duplicating
internal/fractal/root.go

shard map

License
AGPL-3.0 – commercial licenses available.