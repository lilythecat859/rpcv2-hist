# Release Process

1. Bump version in `Makefile`, `helm/Chart.yaml`, `sdk/*/Cargo.toml`, `sdk/js/package.json`, `sdk/python/pyproject.toml`
2. Update CHANGELOG.md
3. Run `scripts/release.sh <version>`
4. Git tag `v<version>` and push
5. Create GitHub release with binaries + checksums
6. Announce on Discord/Twitter