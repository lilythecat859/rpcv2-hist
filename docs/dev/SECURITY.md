# Security Development

## Threat Model
- Malicious JSON-RPC input → sanitize & validate
- Sybil on gRPC → mTLS + rate-limit
- Data tampering → checksum raw blobs
- DoS → rate-limit & K8s HPA

## Checklist
- [ ] gosec scan in CI
- [ ] govulncheck on deps weekly
- [ ] fuzz JSON-RPC handlers
- [ ] mTLS e2e test
- [ ] S3 bucket encryption
- [ ] RBAC least-privilege

## Tools
```bash
gosec ./...
govulncheck ./...