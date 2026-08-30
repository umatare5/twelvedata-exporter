# Repository Instructions

> [!IMPORTANT]
> Read [README.md](README.md) for product overview, metric surface, flags, and operator usage.

## Tech Stack

- Go 1.27+ (see [go.mod](go.mod))
- [`prometheus/client_golang`](https://github.com/prometheus/client_golang) — metric registration and exposition
- [`urfave/cli/v3`](https://github.com/urfave/cli) — CLI flags and application lifecycle
- [`sirupsen/logrus`](https://github.com/sirupsen/logrus) — structured logging
- [`jinzhu/configor`](https://github.com/jinzhu/configor) — configuration loading
- [`goreleaser`](https://goreleaser.com/) v2 — cross-platform release builds (see [.goreleaser.yml](.goreleaser.yml))

## Repository Structure

- `cmd/` — Entry point (`main.go`); calls `cli.Start()`
- `cli/` — CLI flag definitions and app wiring (urfave/cli/v3)
- `config/` — Flag and environment variable parsing, defaults, and validation
- `internal/` — Collector, upstream Twelve Data client, and HTTP server
- `log/` — logrus setup and logging helpers
- `docs/` — Project assets (logos)
- `prometheus.sample.yml` / `prometheus.rules.sample.yml` — Prometheus scrape and rule examples

## Setup and Commands

- `go build ./...` — Build all packages
- `go run ./cmd` — Run the exporter locally
- `golangci-lint run` — Lint (see [.golangci.yml](.golangci.yml)); install per <https://golangci-lint.run/docs/welcome/install/local/>
- `make image` — Build the Docker image (see [Makefile](Makefile))

## Code Style

- `golangci-lint` v2 with the `gci`, `gofumpt`, `goimports`, and `golines` formatters is the single source of truth (see [.golangci.yml](.golangci.yml)).
- Keep metric names, help strings, types, and labels stable unless a SemVer-signaled breaking change is intentional.
- Comments record only what the code cannot say, and never address the reader.

## Testing Instructions

- Run `go build ./...` and `go test ./...` before committing.
- Place tests next to code under test (`*_test.go`). The repository has no unit tests yet.

## Commits and PRs

- Use [Conventional Commits](https://www.conventionalcommits.org/) (`feat:`, `fix:`, `chore(deps):`, etc.).
- Sign off commits with `Signed-off-by:` (DCO).
- Open PRs against `main`. CI runs lint, tests, and CodeQL.
- Call out any metrics, flags, or default scrape-path changes explicitly because they affect Prometheus configs and alerts.

## Domain Knowledge

- This repository ships a Prometheus exporter binary and container image for Twelve Data, not a reusable Go SDK.
- The exporter listens on `0.0.0.0:10016` and serves metrics on `/price` by default; symbols are passed per scrape via `?symbols=`.
- `TWELVEDATA_API_KEY` (or `--twelvedata.api-key`) is the primary runtime credential. Never log it or vendor responses that may expose account-sensitive details.
- Favor low-cardinality labels and stable metric contracts over exposing every upstream field.
