# Releasing

Release builds are intended to use GoReleaser.

Before tagging:

```bash
go test ./...
go vet ./...
./scripts/validate-local.sh
```

Dry run:

```bash
goreleaser release --snapshot --clean
```

Tagged release:

```bash
git tag v0.1.0
goreleaser release --clean
```

The Homebrew formula target is `davemorin/homebrew-tap`.
GitHub Actions needs a `TAP_GITHUB_TOKEN` repository secret with write access to
that tap repository.

Users install from the tap with:

```bash
brew install davemorin/tap/supacrawl
```
