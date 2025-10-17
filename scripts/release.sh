#!/usr/bin/env bash
# AGPL-3.0
set -euo pipefail

VERSION="${1:-}"
if [[ -z "$VERSION" ]]; then
  echo "usage: ./release.sh <version>"
  exit 1
fi

echo "=> Building static binary"
CGO_ENABLED=0 make build

echo "=> Building OCI image"
docker build -t "ghcr.io/faithful-rpc/rpcv2-hist:${VERSION}" .

echo "=> Pushing image"
docker push "ghcr.io/faithful-rpc/rpcv2-hist:${VERSION}"

echo "=> Generating SBOM"
syft "ghcr.io/faithful-rpc/rpcv2-hist:${VERSION}" -o spdx-json > sbom.spdx.json

echo "=> Signing"
cosign sign --yes "ghcr.io/faithful-rpc/rpcv2-hist:${VERSION}"

echo "=> Attest SBOM"
cosign attest --yes --predicate sbom.spdx.json --type spdx "ghcr.io/faithful-rpc/rpcv2-hist:${VERSION}"

echo "✅ Released ${VERSION}"