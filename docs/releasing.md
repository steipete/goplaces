# Releasing

Quick, repeatable release checklist. Mirrors gifgrep cadence.

## Before

- Update `CHANGELOG.md` for the new version.
- Run gate: `./scripts/check-coverage.sh` + `golangci-lint run ./...`.
- Ensure `main` is clean and pushed.
- Ensure `HOMEBREW_TAP_GITHUB_TOKEN` secret is set (pushes formula to `steipete/homebrew-tap`).

## Tag + Release

```bash
git tag vX.Y.Z
git push origin vX.Y.Z
```

GitHub Actions runs GoReleaser on tag push (`.github/workflows/release.yml`).

## Notes

- CLI version set via ldflags in `.goreleaser.yml`:
  `-X github.com/steipete/goplaces/internal/cli.Version={{.Version}}`
- Local smoke build: `make goplaces`
