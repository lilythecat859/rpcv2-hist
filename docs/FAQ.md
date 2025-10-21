# FAQ

## Can I run this on ARM?
Yes, the static binary is `GOARCH=arm64` compatible.

## Does it support PostgreSQL?
Yes, implement `storage.HistoricalStore` and register in `factory.go`.

## How do I back-fill?
Use `tool-parquet` + `migrate-from-bigtable.go`.

## Is re-sharding online?
Yes, fractal root reshards based on slot range; no downtime.

## Can I disable REST?
Set `RESTListen=""` in env.

## Where is the config file?
There is none; only environment variables.

## How big is the binary?
~28 MB compressed; runs from scratch.

## Is there a UI?
No, only JSON-RPC/REST/gRPC. Use any Solana explorer.

## Do you support gRPC reflection?
Yes, enable via `--reflection` flag.

## How do I rotate logs?
Stdout is JSON; use Vector/Loki or Fluent Bit sidecars.