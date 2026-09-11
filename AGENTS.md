# Repository Instructions

> [!IMPORTANT]
> Read [`README.md`](README.md) for product overview, metric surface, flags, and operator usage.

## Tech Stack

- Go 1.27+ (see [`go.mod`](go.mod))
- [`prometheus/client_golang`](https://github.com/prometheus/client_golang) — metric registration and exposition
- [`urfave/cli/v3`](https://github.com/urfave/cli) — CLI flags and application lifecycle
- [`sirupsen/logrus`](https://github.com/sirupsen/logrus) — structured logging
- [`goreleaser`](https://goreleaser.com/) v2 — cross-platform release builds (see [`.goreleaser.yml`](.goreleaser.yml))

## Repository Structure

- `cmd/` — Entry point (`main.go`); calls `cli.Start()`
- `cli/` — CLI flag definitions and app wiring (urfave/cli/v3)
- `config/` — Flag and environment variable parsing, defaults, and validation
- `internal/` — Collector, upstream Twelve Data client, and HTTP server
- `log/` — logrus setup and logging helpers
- `scripts/` — Hook helpers invoked by pre-commit, not by the build
- `docs/` — Reference pages and the logo assets; [Documentation](CONTRIBUTING.md#documentation) names each owner
- `prometheus.sample.yml` / `prometheus.rules.sample.yml` — Prometheus scrape and rule examples

## Setup and Commands

Install required tools (one-time):

- `go install gotest.tools/gotestsum@latest`
- `golangci-lint` — See <https://golangci-lint.run/docs/welcome/install/local/>
- `make pre-commit-install` wires `no-commit-to-main`, `golangci-lint`, `actionlint`, `gitleaks` and `markdownlint-cli2` (see [`.pre-commit-config.yaml`](.pre-commit-config.yaml))

Make targets ([`Makefile`](Makefile)):

- `make build` — Build the binary into `tmp/twelvedata-exporter`
- `make lint` — `golangci-lint run` + `go mod tidy`
- `make test-unit` — Run unit tests via `gotestsum` with coverage
- `make test-unit-coverage` — Generate HTML report at `coverage/report.html`
- `make clean` — Remove build artifacts and `.bak*` files
- `make image` — Build the Docker image (`$USER/twelvedata-exporter`)
- `make pre-commit-install` / `pre-commit-test` / `pre-commit-uninstall` — Manage the pre-commit hooks
- `make help` lists every target, and the `markdownlint-cli2` hook rewrites files with `--fix`

## Code Style

- `golangci-lint` v2 is the single source of truth for format and lint (see [`.golangci.yml`](.golangci.yml)).
- The metric surface is a contract; see [Code Style](CONTRIBUTING.md#code-style) for what changing it costs.
- Say family for one metric name with the HELP and TYPE its series share.
- Comments record only what the code cannot say, and never address the reader.

## Testing

- Run `make build` and `make test-unit` before committing.
- Place tests next to code under test (`*_test.go`). The repository has no unit tests yet.
- See [Testing](CONTRIBUTING.md#testing) for what the first test fixes and what it may not commit.

## Commits and PRs

- Use [Conventional Commits](https://www.conventionalcommits.org/) (`feat:`, `fix:`, `chore(deps):`, etc.).
- Sign off commits with `Signed-off-by:` (DCO).
- Open PRs against `main`. See [Development](CONTRIBUTING.md#development) for what CI runs.
- Record any change to the metric surface under `## [Unreleased]` in [`CHANGELOG.md`](CHANGELOG.md).
- Call out any change to a metric, a flag or the scrape path, which alerts are written against.

## Domain Knowledge

These are the invariants a change has to hold, and the upstream behaviours behind them.

### The API

- **A credit is spent per symbol, not per request**, so batching would not lower it. See [Scrape Path](docs/README.md#scrape-path).
- **An error decodes into the quote struct without failing.** The status mirrors the error code and the client never reads it, so absence is the only signal. See [Absence](docs/README.md#absence).
- **Every price field is a string whose precision the upstream sets**, so parse rather than assume, and expect the field to be missing. See [Specifications](docs/collectors.md#specifications).
- **The key travels in the query string today**, which the header form would avoid, so it reaches every proxy and access log along the path. See [Credential](SECURITY.md#credential).

### The Market

- **A bar is a running aggregate rather than a settled value**, which is why every quote series is a gauge rather than a counter. See [Specifications](docs/collectors.md#specifications).
- **What `previous_close` means depends on the interval**, and the exporter asks for none, so the upstream default decides what the change series measures.
- **A corporate action moves the price without a trade**, and `/quote` cannot adjust for one, so no threshold on the magnitude of a change separates one from a crash.
- **Resolution picks the venue, and its strings drift**, renaming series. See [Labels](docs/collectors.md#labels).
- **Nothing published carries market state**, so a stale bar reads as a live one.

> [!IMPORTANT]
> Write a sentence here only after a live reply showed the behaviour. See [Documentation](CONTRIBUTING.md#documentation) for why, and the page that owns each fact for the measured values.
