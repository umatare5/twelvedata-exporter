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
- ⏱️ **Bounded Upstream Calls**: Each request is capped at ten seconds, so a stall cannot hold a scrape open
- 📐 **Indicator Examples**: RSI written as a recording rule, in [`prometheus.rules.sample.yml`](./prometheus.rules.sample.yml)

> [!IMPORTANT]
> The exporter refuses to start without an API key, so generate one first — see [Getting Started — Authentication](https://twelvedata.com/docs#authentication).

> [!NOTE]
> The exporter spends one credit per symbol per scrape, so the symbol count and the scrape interval together decide whether a plan holds. See [Pricing](https://twelvedata.com/pricing), [Credits](https://support.twelvedata.com/en/articles/5615854-credits) and [Stock Exchanges](https://support.twelvedata.com/en/collections/2787973-stock-exchanges).

## Quick Start

### 1. Run the exporter with Docker

```bash
docker run -p 10016:10016 -e TWELVEDATA_API_KEY ghcr.io/umatare5/twelvedata-exporter
```

The image declares `10016/tcp` without publishing it, so `-p` is what makes the exporter reachable, and `-e` forwards the key from the shell rather than baking it into the image.

The published tags are `latest`, `vX`, `vX.Y` and `vX.Y.Z`. Each one is a multi-platform image covering `linux/amd64` and `linux/arm64`, so Docker selects the architecture of the host.

> [!TIP]
> If you prefer using binaries, download them from the [Release](https://github.com/umatare5/twelvedata-exporter/releases).
>
> **Supported Platform:** `linux_amd64`, `linux_arm64`, `darwin_amd64`, `darwin_arm64` and `windows_amd64`

### 2. Scrape it

See [Prometheus Configuration](#prometheus-configuration) for the job and the recording rules.

## Syntax

`twelvedata-exporter --help` prints every flag, and [`docs/help.md`](docs/help.md) carries the same list.

## Configuration

The API key is the one required setting. `TWELVEDATA_API_KEY` and `--twelvedata.api-key` carry the same value, but the flag reaches the process table where every account on the host reads it, so prefer the environment variable. The exporter exits at start-up when neither is set.

## Endpoints

The exporter serves two endpoints:

- `/` — landing page, which prints the query format when reached at <http://localhost:10016/>
- `/price` — metrics endpoint, configurable via `--web.scrape-path`

> [!IMPORTANT]
> The `symbols` parameter takes a comma-separated list and may repeat, and every occurrence is concatenated into one list. A request carrying none returns 200 with an empty body, so Prometheus counts that scrape as successful while every series is absent — alert on the absence of `twelvedata_price`, not on `up`.

## Metrics

Every quote series is a gauge, and one scrape publishes all five for each symbol it resolved:

| Metric                            | Type  | Description                                 |
| :-------------------------------- | :---- | :------------------------------------------ |
| `twelvedata_price`                | Gauge | Real-time or the latest available price     |
| `twelvedata_previous_close_price` | Gauge | Closing price of the previous day           |
| `twelvedata_change_price`         | Gauge | Change since the previous close             |
| `twelvedata_change_percent`       | Gauge | Change since the previous close, in percent |
| `twelvedata_volume`               | Gauge | Trading volume during the bar               |

The same four labels are attached to all five, and only `symbol` is chosen by the operator:

| Label      | Holds                                       |
| :--------- | :------------------------------------------ |
| `symbol`   | The symbol as the scrape URL spelled it     |
| `name`     | The instrument name the quote carried       |
| `exchange` | The exchange the quote was taken from       |
| `currency` | The currency the row's prices are quoted in |

> [!IMPORTANT]
> `twelvedata_price` is computed as `previous_close + change` rather than read from the quote's `close` field, so it agrees with the two series beside it at every instant. A field that fails to parse becomes `0`, which no label distinguishes from a genuine zero. A symbol whose request fails is skipped, so its series are absent rather than zero.

### Exporter Health Metrics

These series describe the exporter itself rather than the quotes it fetches. They carry no labels, and [`docs/health.md`](docs/health.md) carries the whole set with the reading each one needs.

| Metric                              | Type    | Description                             |
| :---------------------------------- | :------ | :-------------------------------------- |
| `twelvedata_queries_total`          | Counter | Count of completed queries              |
| `twelvedata_failed_queries_total`   | Counter | Count of failed queries                 |
| `twelvedata_query_duration_seconds` | Summary | Duration of queries to the upstream API |

> [!IMPORTANT]
> Read none of the three as a health signal without [`docs/health.md`](docs/health.md). `twelvedata_failed_queries_total` has no increment path and stays `0`, and `twelvedata_query_duration_seconds` observes an instant against itself, so its `_sum` stays `0` too.

## Use Cases

### Basic Usage

No symbol is named until a scrape arrives, so the exporter starts with the key alone.

```bash
$ TWELVEDATA_API_KEY="foobarbaz" ./twelvedata-exporter
INFO[0000] Starting the Twelvedata exporter on 0.0.0.0:10016
```

Open <http://localhost:10016/> for the query format and the example URLs it prints.

### Prometheus Configuration

#### Job Configuration Example

Add the job from [`prometheus.sample.yml`](./prometheus.sample.yml) to your Prometheus configuration. The exporter reads its symbols from `params.symbols`, so that list and the job's `scrape_interval` set the credit spend.

#### Recording Rules Configuration Example

Add the rules from [`prometheus.rules.sample.yml`](./prometheus.rules.sample.yml) to your configuration. They derive the indicators once per evaluation rather than in every dashboard query.

## Contributing

See [`CONTRIBUTING.md`](CONTRIBUTING.md) for the `make` targets, the Docker build and the release process.

## Acknowledgement

I ran [Marco Paganini](https://github.com/marcopaganini)'s [quotes-exporter](https://github.com/marcopaganini/quotes-exporter) until an upstream endpoint change broke it and it was archived. This one builds on his, with thanks to him and to [Tristan Colgate-McFarlane](https://github.com/tcolgate), whose [yquotes-exporter](https://github.com/tcolgate/yquotes_exporter) came first.

## Licence

MIT. The binary statically links Apache-2.0, MIT and BSD 3-Clause dependencies, whose notices are reproduced in [`NOTICE`](NOTICE) and shipped alongside [`LICENSE`](LICENSE) in every release archive and container image.
