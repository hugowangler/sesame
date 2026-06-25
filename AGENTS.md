# AGENTS.md

## Repo Shape
- Go 1.19 CLI module: `github.com/hugowangler/sesame`.
- Entrypoint is `main.go`, which delegates to Cobra commands in `cmd/`.
- Git parsing and repository discovery live in `internal/git`; config read/write logic lives in `internal/config`.
- There are currently no test files, Makefile, lint config, or pre-commit config.

## Commands
- Build/install locally: `go install`.
- Run the CLI without installing: `go run . --help`.
- Verify compile/tests across packages: `go test ./...`.
- Keep dependencies tidy before release-impacting changes: `go mod tidy`.

## Config And Manual Testing Gotchas
- Default config path is `$HOME/.config/.sesame.yaml`.
- Commands that store repositories can write config; use `--config` with an existing temp file when testing write paths to avoid changing the real user config.
- `sesame open` may launch a browser through `github.com/pkg/browser`; prefer `--help`, unit tests, or isolated config tests unless browser behavior is intentional.
- Repository names, owners, and hosts parsed from remotes are lowercased before storage.

## Release
- Releases are driven by `.github/workflows/goreleaser.yaml` on published GitHub releases.
- GoReleaser runs `go mod tidy`, builds with `CGO_ENABLED=0`, and targets linux/windows/darwin.
