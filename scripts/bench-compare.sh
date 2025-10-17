#!/usr/bin/env bash
# AGPL-3.0
set -euo pipefail

echo "=> BigTable baseline (70k USD/mo) vs rpcv2-hist (ClickHouse)"

# Requires existing BT credentials
go test -bench=BenchmarkGetSignaturesForAddress \
  -benchmem -count=10 -benchtime=10s ./internal/bench/... | tee bt.txt

# Start local ClickHouse via docker compose
docker compose up -d clickhouse
sleep 10
go test -bench=BenchmarkGetSignaturesForAddress \
  -benchmem -count=10 -benchtime=10s ./internal/bench/... |tee ch.txt

echo "=> Results"
benchstat bt.txt ch.txt