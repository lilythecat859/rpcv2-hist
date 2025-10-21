# Contributing Code

1. Fork and branch (`git checkout -b feat/my-feature`)
2. Write tests + bench
3. Run `make lint test bench`
4. Commit with DCO sign-off (`git commit -s`)
5. Push and open PR

## Style
- Standard gofmt / goimports
- No naked returns
- Context as first param
- zap for logging

## Testing
- Unit tests in `*_test.go`
- Integration tests in `integration_test.go`
- Benchmarks in `bench_test.go`
- Use testcontainers for ClickHouse

## CI
GitHub Actions runs lint, test, bench on every PR.