# Repository Instructions

> [!IMPORTANT]
> Read [`README.md`](README.md) for product overview, metric surface, flags, and operator usage.

## Tech Stack

- Go 1.27+ (see [`go.mod`](go.mod))
- [`prometheus/client_golang`](https://github.com/prometheus/client_golang) — metric registration and exposition
- [`urfave/cli/v3`](https://github.com/urfave/cli) — CLI flags and application lifecycle
- [`sirupsen/logrus`](https://github.com/sirupsen/logrus) — structured logging
- [`jinzhu/configor`](https://github.com/jinzhu/configor) — configuration loading
- [`goreleaser`](https://goreleaser.com/) v2 — cross-platform release builds (see [`.goreleaser.yml`](.goreleaser.yml))

## Repository Structure

- `cmd/` — Entry point (`main.go`); calls `cli.Start()`
- `cli/` — CLI flag definitions and app wiring (urfave/cli/v3)
- `config/` — Flag and environment variable parsing, defaults, and validation
- `internal/` — Collector, upstream Twelve Data client, and HTTP server
- `log/` — logrus setup and logging helpers
- `scripts/` — Hook helpers invoked by pre-commit, not by the build
- `docs/` — Reference pages the README delegates to, and the logo assets
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

## Code Style

- `golangci-lint` v2 with the `gci`, `gofumpt`, `goimports`, and `golines` formatters is the single source of truth (see [`.golangci.yml`](.golangci.yml)).
- Keep metric names, help strings, types, and labels stable unless a SemVer-signaled breaking change is intentional.
- Comments record only what the code cannot say, and never address the reader.

## Testing

- Run `make build` and `make test-unit` before committing.
- Place tests next to code under test (`*_test.go`). The repository has no unit tests yet.

## Commits and PRs

- Use [Conventional Commits](https://www.conventionalcommits.org/) (`feat:`, `fix:`, `chore(deps):`, etc.).
- Sign off commits with `Signed-off-by:` (DCO).
- Open PRs against `main`. CI runs lint, tests, CodeQL, govulncheck, promtool, markdownlint and lychee.
- Record any change to the metric surface under `## [Unreleased]` in [`CHANGELOG.md`](CHANGELOG.md).
- Call out any metrics, flags, or default scrape-path changes explicitly because they affect Prometheus configs and alerts.

## Domain Knowledge

These are the upstream behaviours every change to the client or the collector has to hold.

- **Credits are spent per symbol, against a per-minute allowance.** `/quote` costs one credit per symbol, so folding symbols into one request would not lower the spend, and `/api_usage` reports `plan_limit` as requests per minute.
- **No response header reports the remaining budget.** A reply carries no rate-limit, credit or `Retry-After` header, so the allowance left is visible only through `/api_usage`, which costs a credit of its own.
- **An error arrives as a `code`/`message`/`status` object that decodes into the quote struct without failing.** The HTTP status mirrors `code` — a rejected key answers `401` and a missing `symbol` answers `404` — so the status is what separates an error from a quote. A body-only check sees empty strings instead.
- **Every price field is a JSON string whose precision `dp` sets, five places by default.** Parse rather than assume, and expect absence. `volume`, `average_volume`, `mic_code` and `currency` are documented as unavailable for some instrument types. The `rolling_*` and `extended_*` fields appear only when the plan and the request parameters ask for them.
- **`/quote` describes a bar, not a tick.** Its default `interval` is `1day`, and `change` and `percent_change` are taken against `previous_close`. The same bar repeats while `is_market_open` is false, so an overnight scrape holds a value rather than going absent.
- **The strings the labels carry are the API's resolution, not the operator's input.** A bare `symbol` is ambiguous across venues, which `exchange`, `mic_code` and `country` narrow, and the reply echoes what it chose. Those strings drift, and a changed value renames every series carrying it. The documented `AAPL` example gives `mic_code` `XNAS` and `name` `Apple Inc`, where a live reply gives `XNGS` and `Apple Inc.`
- **The key belongs in the `Authorization: apikey <key>` header the documentation recommends.** The `?apikey=` form the client sends today puts the credential in every proxy and access log along the path. Never log it, and never vendor a response that carries account-level detail.
- **`apikey=demo` answers for a small set of symbols alone.** Everything else, a comma-separated list included, comes back `401`, so a demo request proves reachability rather than behaviour.
