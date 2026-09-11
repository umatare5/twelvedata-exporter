<div align="center">

  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="./docs/assets/logo_dark.png" width="115px" />
    <source media="(prefers-color-scheme: light)" srcset="./docs/assets/logo.png" width="115px" />
    <img alt="twelvedata-exporter" src="./docs/assets/logo.png" width="115px" />
  </picture>

  <h1>twelvedata-exporter</h1>

  <p>A third-party Prometheus Exporter for Twelve Data.</p>

  <p>
    <img alt="GitHub Tag" src="https://img.shields.io/github/v/tag/umatare5/twelvedata-exporter?label=Latest%20version" />
    <a href="https://github.com/umatare5/twelvedata-exporter/actions/workflows/go-test-build.yml"><img alt="Test and Build" src="https://github.com/umatare5/twelvedata-exporter/actions/workflows/go-test-build.yml/badge.svg?branch=main" /></a>
    <a href="https://github.com/umatare5/twelvedata-exporter/actions/workflows/go-vulncheck.yml"><img alt="govulncheck" src="https://github.com/umatare5/twelvedata-exporter/actions/workflows/go-vulncheck.yml/badge.svg?branch=main" /></a><br>
    <a href="https://pkg.go.dev/github.com/umatare5/twelvedata-exporter@main"><img alt="Go Reference" src="https://pkg.go.dev/badge/umatare5/twelvedata-exporter.svg" /></a>
    <a href="./LICENSE"><img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg" /></a>
  </p>

</div>

## Overview

This exporter fetches quotes from [Twelve Data](https://twelvedata.com/) and serves them as Prometheus metrics.

- 💹 **Quote Surface**: Price, previous close, change, percent change and volume per symbol
- 🔎 **Per-Scrape Symbols**: The symbol list travels in the scrape URL, not in a config file
- ⏱️ **Bounded Upstream Calls**: Each request is capped at ten seconds by the client
- 📐 **Indicator Examples**: RSI written as a recording rule, in [`prometheus.rules.sample.yml`](./prometheus.rules.sample.yml)

> [!IMPORTANT]
> The exporter refuses to start without an API key — generate one from [Getting Started](https://twelvedata.com/docs#authentication).

> [!NOTE]
> The exporter spends one credit per symbol per scrape, so the symbol count and the scrape interval together decide whether a plan holds. See [Pricing](https://twelvedata.com/pricing), [Credits](https://support.twelvedata.com/en/articles/5615854-credits) and [Stock Exchanges](https://support.twelvedata.com/en/collections/2787973-stock-exchanges).

## Quick Start

### 1. Set the API token

```bash
export TWELVEDATA_API_KEY="your-twelvedata-api-token"
```

### 2. Run the exporter with Docker

```bash
docker run -p 10016:10016 -e TWELVEDATA_API_KEY ghcr.io/umatare5/twelvedata-exporter
```

> [!WARNING]
> The per-architecture tags — `latest-amd64`, `latest-arm64` and their `vX-`, `vX.Y-` and `vX.Y.Z-` counterparts — are **deprecated and no longer published**. They stopped receiving updates after v1.1.0, so `latest-amd64` and `v1-amd64` still resolve to v1.1.0 and never move again. Pull `latest`, `vX`, `vX.Y` or `vX.Y.Z` instead, each of which serves both `linux/amd64` and `linux/arm64` from v1.0.2 on.

> [!TIP]
> If you prefer using binaries, download them from the [Release](https://github.com/umatare5/twelvedata-exporter/releases).
>
> **Supported Platform:** `linux_amd64`, `linux_arm64`, `darwin_amd64`, `darwin_arm64` and `windows_amd64`

### 3. Scrape it

See [Prometheus Configuration](#prometheus-configuration) for the job and the recording rules.

## Flags

`twelvedata-exporter --help` prints every flag, and [`docs/help.md`](docs/help.md) carries it with notes.

- `--twelvedata.api-key` is the one required flag, and `TWELVEDATA_API_KEY` fills it.
- `--web.listen-address`, `--web.listen-port` and `--web.scrape-path` place the endpoint.
- No `--collector.*` flag exists, so every scrape publishes the whole quote surface.

## Endpoints

The exporter serves two endpoints:

- `/` — landing page, which prints the query format when reached at <http://localhost:10016/>
- `/price` — metrics endpoint, configurable via `--web.scrape-path`

See [Endpoints](docs/README.md#endpoints) for the method and status each keeps, and [Scrape Path](docs/README.md#scrape-path) for how `symbols` is read.

## Metrics

The exporter publishes one quote surface and its own health, catalogued in `docs/`:

| Page                                  | Covers                                          |
| :------------------------------------ | :---------------------------------------------- |
| **[Collectors](docs/collectors.md)**  | The five quote series, their labels and meaning |
| **[Exporter health](docs/health.md)** | The exporter's own series and how to read them  |

Every quote series is a gauge, and one scrape publishes all five for each symbol it resolved:

| Metric                            | Type  | Description                                 |
| :-------------------------------- | :---- | :------------------------------------------ |
| `twelvedata_price`                | Gauge | Real-time or the latest available price     |
| `twelvedata_previous_close_price` | Gauge | Closing price of the previous day           |
| `twelvedata_change_price`         | Gauge | Change since the previous close             |
| `twelvedata_change_percent`       | Gauge | Change since the previous close, in percent |
| `twelvedata_volume`               | Gauge | Trading volume during the bar               |

> [!NOTE]
> See [`docs/README.md`](docs/README.md) for the absence, scrape-path and counter rules every series shares.

> [!IMPORTANT]
> A symbol whose request fails is skipped rather than published as zero, and the scrape still answers 200. Alert on `absent(twelvedata_price)` rather than on `up`, which stays 1 through every upstream failure.

### Exporter Health Metrics

These series describe the exporter itself rather than the quotes it fetches. They carry no labels, and [`docs/health.md`](docs/health.md) carries the whole set with the reading each one needs.

| Metric                              | Type    | Description                             |
| :---------------------------------- | :------ | :-------------------------------------- |
| `twelvedata_queries_total`          | Counter | Count of completed queries              |
| `twelvedata_failed_queries_total`   | Counter | Count of failed queries                 |
| `twelvedata_query_duration_seconds` | Summary | Duration of queries to the upstream API |

> [!NOTE]
> Read none of the three without [Specifications](docs/health.md#specifications): one never increments and one measures nothing.

## Examples

### Command Lines

No symbol is named until a scrape arrives, so the exporter starts with the key alone.

```bash
$ TWELVEDATA_API_KEY="foobarbaz" ./twelvedata-exporter
INFO[0000] Starting the Twelvedata exporter on 0.0.0.0:10016
```

Open <http://localhost:10016/> for the query format and the example URLs it prints.

### Prometheus Configuration

#### Job Configuration Example

Add the job from [`prometheus.sample.yml`](./prometheus.sample.yml) to your Prometheus configuration.

#### Recording Rules Configuration Example

Add the rules from [`prometheus.rules.sample.yml`](./prometheus.rules.sample.yml) to your Prometheus configuration.

## Contributing

See [`CONTRIBUTING.md`](CONTRIBUTING.md) for the development setup, the tests and the conventions.

## License

MIT. This exporter builds on [quotes-exporter](https://github.com/marcopaganini/quotes-exporter) and [yquotes-exporter](https://github.com/tcolgate/yquotes_exporter), which came before it. The binary statically links Apache-2.0, MIT and BSD 3-Clause dependencies, whose notices are reproduced in [`NOTICE`](NOTICE) and shipped alongside [`LICENSE`](LICENSE) in every release archive and container image.
