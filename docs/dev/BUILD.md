# Build

## Requirements
- Go 1.22+
- Docker (optional)
- Make

## Static binary
```bash
CGO_ENABLED=0 make build
```

## Cross-compile
```
GOOS=linux GOARCH=arm64 make build
```

Docker image
```
make image
```

Generate protos
```
buf generate
```

Run tests
```
make test
```

Run benches
```
make bench
```