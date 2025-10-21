# Deploy

## Binary
```bash
CGO_ENABLED=0 make build
./bin/rpcv2-hist

##Docker
```docker compose up -d

##Kubernetes
```kubectl apply -f kubernetes/

#$Terraform

See

```terraform/

(community maintained).

##Helm
```helm install rpcv2-hist ./helm

##Migration from BigTable
```go run scripts/migrate-from-bigtable.go \
  -bt-project my-gcp-proj \
  -bt-instance solana \
  -bt-table mainnet \
  -ch-addr clickhouse:9000

##Tuning

Run

```scripts/clickhouse-tuning.sh

on each ClickHouse node.