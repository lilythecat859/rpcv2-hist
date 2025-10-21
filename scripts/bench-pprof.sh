#!/usr/bin/env bash
# AGPL-3.0
set -euo pipefail

go test -bench=BenchmarkGetSignaturesForAddress -benchmem -cpuprofile=cpu.prof -memprofile=mem.prof ./internal/bench/...
go tool pprof -http=:8080 cpu.prof