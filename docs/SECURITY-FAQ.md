# Security FAQ

## How are secrets managed?
All secrets are injected via env vars or mounted files; nothing is baked into images.

## Is mTLS mandatory?
No, but enabled by default in Kubernetes via Linkerd or Istio sidecars.

## What crypto is used?
- Blake3 for deterministic sharding
- TLS 1.3 for transport
- JWT HS256 for auth tokens

## Do you store private keys?
Never. The service is stateless and holds no keys.

## How do I report vulnerabilities?
Email security@faithful-rpc.org with PoC and severity. 90-day disclosure window.