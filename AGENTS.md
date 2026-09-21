# Repository Instructions

> [!IMPORTANT]
> Read [`README.md`](README.md) for the project overview.

## Tech Stack

- Go 1.27+ (see [`go.mod`](go.mod))
- [`prometheus/client_golang`](https://github.com/prometheus/client_golang) v1.24+ – metric registration and HTTP handler
- [`urfave/cli/v3`](https://github.com/urfave/cli) v3.11+ – CLI flags and application lifecycle
- [`sirupsen/logrus`](https://github.com/sirupsen/logrus) – structured logging
- [`goreleaser`](https://goreleaser.com/) v2 – cross-platform release builds (see [`.goreleaser.yml`](.goreleaser.yml))

## Repository Structure

Read from `cmd/main.go`. Each package is named for what it owns.

- [`cmd/main.go`](cmd/main.go) – application entry point
- [`cli/`](cli) – command-line flags, defaults and app wiring
- [`config/`](config) – flag reads, defaults and API key validation
- [`internal/collector.go`](internal/collector.go) – metric descriptions and collection logic
- [`internal/twelvedata.go`](internal/twelvedata.go) – upstream API client and data structures
- [`internal/server.go`](internal/server.go) – HTTP server configuration and routing
- [`log/`](log) – logrus level and formatter setup
- [`docs/`](docs) – reference pages behind the README
- [`scripts/`](scripts) – helper scripts the pre-commit hooks run

## Setup and Commands

Run `make pre-commit-install` first.

- Read [`Makefile`](Makefile) which lists all available make targets and their descriptions.
- Read [`CONTRIBUTING.md`](CONTRIBUTING.md) which provides guidelines for contributing to the project.

## Code Style

Follow [Effective Go](https://go.dev/doc/effective_go) conventions and the software development principles DRY/YAGNI/SRP.

- Keep code simple and readable, avoiding clever tricks that obscure intent.
- Keep minimal for all changes, coding, testing, commenting, and documentation.
- Write simple comments that explain the reasoning behind the code, not just what it does.

## Testing

Follow [`CONTRIBUTING.md`](CONTRIBUTING.md).

- Run `make test-unit` before creating a commit.

## Commits and PRs

Follow [Conventional Commits](https://www.conventionalcommits.org/) (`feat:`, `fix:`, `chore(deps):`, etc.).

- Run pre-commit and ensure all hooks pass before committing.
- Must Sign off all commits with `Signed-off-by:` (DCO).
- Open PRs against `main`. Create Draft PR as default.

## Domain Knowledge

Learn the constraints outside the exporter, because they decide what a metric can mean.

### About the API

An error decodes into the quote struct, so absence is the only signal. See [Absence](docs/architecture.md#absence).

- A credit is spent per symbol, not per request, so batching would not lower the cost.
- The client never reads the status, so a rejected key and a bad symbol look alike.
- A price field is a string whose precision the upstream sets, so parse rather than assume.
- A field may be missing entirely, so a failed parse publishes `0` in its place.
- An `apikey` parameter beats the header, so the symbol is encoded before reaching the query.

### About the market

A bar is a running aggregate, so every quote series is a gauge. See [Metrics](README.md#collector-metrics).

- The exporter sends no `interval`, so the upstream default of one day sets what a bar means.
- A corporate action moves the price without a trade, because `/quote` cannot adjust for one.
- Resolution picks the venue, so a drifting exchange string renames every series.
- Nothing published carries market state, so a stale bar reads as a live one.
