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

### The API

- **Credits are spent per symbol, against a per-minute allowance.** `/quote` costs one credit per symbol, so folding symbols into one request would not lower the spend, and `/api_usage` reports `plan_limit` as requests per minute.
- **No response header reports the remaining budget.** A reply carries no rate-limit, credit or `Retry-After` header, so the allowance left is visible only through `/api_usage`, which costs a credit of its own.
- **An error arrives as a `code`/`message`/`status` object that decodes into the quote struct without failing.** The HTTP status mirrors `code` — a rejected key answers `401` and a missing `symbol` answers `404` — so the status is what separates an error from a quote. A body-only check sees empty strings instead.
- **Every price field is a JSON string whose precision `dp` sets, five places by default.** Parse rather than assume, and expect absence. `volume`, `average_volume`, `mic_code` and `currency` are documented as unavailable for some instrument types. The `rolling_*` and `extended_*` fields appear only when the plan and the request parameters ask for them.
- **The key belongs in the `Authorization: apikey <key>` header the documentation recommends.** The `?apikey=` form the client sends today puts the credential in every proxy and access log along the path. Never log it, and never vendor a response that carries account-level detail.
- **`apikey=demo` answers for a small set of symbols alone.** Everything else, a comma-separated list included, comes back `401`, so a demo request proves reachability rather than behaviour.

### The market

- **A bar is a running aggregate, not a settled value.** `/quote` defaults to `interval=1day`, so `close` moves with every scrape of a live session and `volume` accumulates from the open and resets at the next one. A gauge is the only honest type for either.
- **`previous_close` means the previous bar, and `interval` decides what a bar is.** At `1day` it is the prior session's close, which is what the metric help claims; at `1min` a live reply gives `previous_close` `319.82999` against `close` `319.98999`. Changing `interval` silently changes what `twelvedata_change_price` measures.
- **Outside the session the same bar repeats.** `is_market_open` goes false while `datetime` stays on the last session, and Prometheus stamps each sample with scrape time, so a stale close is indistinguishable from a live one. Nothing the exporter publishes carries the market state that would separate them.
- **A corporate action moves the price without a trade.** `/time_series` takes `adjust` and defaults it to `splits`; `/quote` takes no such parameter, so its treatment of a split is undocumented. An unadjusted 4-for-1 reads as a −75% move, which is why no alert should fire on the magnitude of `twelvedata_change_percent` alone.
- **The venue is chosen for you, and it fixes the currency.** A bare `symbol` resolves to one listing, which `exchange`, `mic_code` and `country` narrow. Cross-listings quote in their own currency, so ranking `twelvedata_price` across symbols only means something within one.
- **The strings that resolution returns are label values, and they drift.** The documented `AAPL` example gives `mic_code` `XNAS` and `name` `Apple Inc`, where a live reply gives `XNGS` and `Apple Inc.` — a changed value renames every series carrying it.
- **Instruments outside equities drop the fields the labels need.** A live `EUR/USD` or `BTC/USD` reply carries no `currency`, `mic_code`, `volume` or `average_volume`, so the `currency` label goes empty and the volume gauge reads a parsed `0`. Both markets also keep `is_market_open` true around the clock.
- **`dp` rounds to decimal places, not to significant figures.** Five is generous for a US equity and coarse for an FX pair, where `EUR/USD` at `dp=2` returns `change` `-0.00`. Asking for more surfaces the float32 the API stores: `AAPL` at `dp=11` returns `319.97000122070`.
- **An indicator computed in Prometheus is not the indicator a chart shows.** The RSI in [`prometheus.rules.sample.yml`](prometheus.rules.sample.yml) runs over scrape samples, which repeat one daily bar for the length of a session, so its period is wall-clock rather than bars. It demonstrates the mechanism; it is not a trading signal.
