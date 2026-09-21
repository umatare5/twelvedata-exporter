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
    <a href="https://github.com/umatare5/twelvedata-exporter/actions/workflows/go-vulncheck.yml"><img alt="govulncheck" src="https://github.com/umatare5/twelvedata-exporter/actions/workflows/go-vulncheck.yml/badge.svg?branch=main" /></a>
    <a href="./LICENSE"><img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg" /></a>
  </p>

</div>

## Overview

This exporter allows a Prometheus instance to monitor prices of stocks, ETFs, and mutual funds on [Twelve Data](https://twelvedata.com/).

- 💹 **Major Metrics**: Expose price, previous close, change, percent change and volume per symbol.
- 📊 **Multiple Symbols**: Fetch multiple symbols in a single scrape, each spending one API credit.
- 📈 **Technical Analysis**: Allow to evaluate stock prices utilizing PromQL and other tools in the ecosystem.
- 🔔 **Signal-based Alerting**: Allow to monitor technical analysis indicators utilizing the alerting rules.

> [!NOTE]
>
> The Twelvedata API has some limitations based on the license. For example, API limit, accessible market and others. For the limitations, please refer to [twelvedata - Pricing](https://twelvedata.com/pricing) and [Twelvedata Support - Credits](https://support.twelvedata.com/en/articles/5615854-credits).

## Installation

This exporter supports both container images and OS-specific binaries installations.

**A. Using Container**

```bash
docker pull ghcr.io/umatare5/twelvedata-exporter
```

**B. Using OS-Specific binaries**

Download from [Releases](https://github.com/umatare5/twelvedata-exporter/releases). `linux_(amd64|arm64)`, `darwin_(amd64|arm64)` and `windows_amd64` are supported.

## Quick Start

This exporter needs an API key. See the **[Twelve Data Developer Docs](https://twelvedata.com/docs/introduction/overview)** to get an API key first.

**1. Set the API key**

```bash
export TWELVEDATA_API_KEY="your-twelvedata-api-token"
```

**2. Run the exporter with Docker**

```bash
docker run -p 10016:10016 -e TWELVEDATA_API_KEY ghcr.io/umatare5/twelvedata-exporter:v1.3.0
```

**3. Scrape it**

```bash
curl http://localhost:10016/price?symbols=SPY
```

> [!TIP]
>
> See [Metrics](#metrics) for available metrics, and [Prometheus Configuration](#prometheus-configuration) for the job and the alerting rules.

## Flags

The exporter supports the following command-line flags:

```text
NAME:
   twelvedata-exporter - Fetch quotes from Twelvedata API

USAGE:
   twelvedata-exporter COMMAND [options...]

VERSION:
   1.3.0

GLOBAL OPTIONS:
   --web.listen-address string, -I string  Set IP address (default: "0.0.0.0")
   --web.listen-port int, -P int           Set port number (default: 10016)
   --web.scrape-path string, -p string     Set the path to expose metrics (default: "/price")
   --twelvedata.api-key string, -a string  Set key to use twelvedata API [$TWELVEDATA_API_KEY]
   --help, -h                              show help
   --version, -v                           print the version
```

## Endpoints

The exporter serves two endpoints. See [Endpoints](docs/architecture.md#endpoints) for the details.

| Path     | Detail                                                                    |
| :------- | :------------------------------------------------------------------------ |
| `/`      | Landing page – confirming the exporter is up at <http://localhost:10016/> |
| `/price` | Metrics endpoint – the scrape path that executes API calls                |

## Metrics

This exporter exposes metrics for upstream quotes and internal operational state.

### Collector Metrics

The following table lists the metrics this exporter publishes. See **Appendix 1** below the table for more details.

| Metric                            | Type  | Description                                 |
| :-------------------------------- | :---- | :------------------------------------------ |
| `twelvedata_price`                | Gauge | Real-time or the latest available price     |
| `twelvedata_previous_close_price` | Gauge | Closing price of the previous day           |
| `twelvedata_change_price`         | Gauge | Change since the previous close             |
| `twelvedata_change_percent`       | Gauge | Change since the previous close, in percent |
| `twelvedata_volume`               | Gauge | Trading volume during the bar               |

<details><summary><b>Appendix 1 - Collector Metrics Details</b></summary><p>

#### About the Metrics

**`twelvedata_price`**: The price is computed dynamically as `previous_close + change` rather than read directly from the reply's `close` field. This ensures the value mathematically agrees with the two adjacent series at every instant. An unparsed previous close collapses this metric entirely onto `twelvedata_change_price`.

**`twelvedata_previous_close_price`**: This metric reflects the close of the previous bar, whose dimension the upstream `interval` parameter defines. The exporter sends no `interval`, forcing the upstream default of one day. Consequently, this series strictly represents the prior session's close.

**`twelvedata_change_price`**: The change is read verbatim from the reply's `change` field rather than computed. `twelvedata_price` derives from it as `previous_close + change`, so this metric is the input to that sum rather than its result. A `change` that fails to parse reads `0`, collapsing `twelvedata_price` onto the previous close.

**`twelvedata_change_percent`**: The percentage change shifts identically on a corporate action as it does on a standard trade. The upstream `/quote` endpoint accepts no adjustment parameters. An unadjusted four-for-one stock split will read as a −75% market move, and no internal threshold separates this from an actual market crash.

**`twelvedata_volume`**: The volume counts only the trades executed inside the current bar. This reset behavior at the start of each new bar dictates its Gauge type. Modeling this as a Counter would cause Prometheus to misinterpret each daily reset as a counter rollover.

#### About the Labels

**`symbol`**: The symbol label uses the exact spelling returned by the upstream reply, ignoring the casing provided in the scrape URL. A query for `?symbols=aapl` labels the resulting series `AAPL`. A recording rule written against the lower-cased request name will subsequently miss the generated series.

**`exchange`**: The exchange label carries the venue name exactly as the reply provides it. While the reply also offers a MIC (Market Identifier Code), the MIC never becomes a label. A live `AAPL` reply labels `NASDAQ` while its `mic_code` reads `XNGS`, against the `XNAS` the upstream documentation shows.

**`name`**: The name label reflects whatever the live reply spells for the instrument. The upstream documentation and the live service frequently disagree on formatting. A changed upstream value silently renames every series carrying it.

**`currency`**: The currency label is fixed by the specific venue the request resolved to, rather than by the symbol asked for. A cross-listing quotes in its own currency, rendering raw value comparisons across symbols meaningless without external conversion.

</p></details>

### Exporter Health Metrics

The following table lists the metrics this exporter publishes. See **Appendix 2** below the table for more details.

| Metric                              | Type    | Description                                         |
| :---------------------------------- | :------ | :-------------------------------------------------- |
| `twelvedata_queries_total`          | Counter | Count of completed queries                          |
| `twelvedata_failed_queries_total`   | Counter | Count of failed queries                             |
| `twelvedata_query_duration_seconds` | Summary | Duration of queries to the upstream API             |
| `twelvedata_http_requests_total`    | Counter | Declared with the quote labels, never given a child |

<details><summary><b>Appendix 2 - Exporter Health Metrics Details</b></summary><p>

#### About the Metrics

**`twelvedata_queries_total`**: This counter increments exactly once when a scrape enters the internal collector, regardless of how many symbols that scrape requests. Consequently, the actual upstream request rate billed to the account equals this counter's rate multiplied by the symbol count per scrape.

**`twelvedata_failed_queries_total`**: This counter is declared and registered within the code but is never programmatically incremented, so it perpetually stays at `0` irrespective of upstream answers or internal failures. Any alert or dashboard built relying on this metric will remain dead. Since `twelvedata_failed_queries_total` never increments, a failed symbol is exclusively visible through the sudden absence of its corresponding quote series.

**`twelvedata_query_duration_seconds`**: This summary times the upstream call from the request through the parsed response, so `_sum` divided by `_count` yields the mean latency of a query that produced a quote. The observation occurs strictly after the response parses successfully. A request that fails entirely or returns a nameless quote reaches neither `_count` nor `_sum`, leaving `_count` a count of successes rather than attempts and keeping the ten-second timeout out of `_sum`.

**`twelvedata_http_requests_total`**: Because `twelvedata_http_requests_total` never receives a child increment, it publishes absolutely nothing instead of exposing a zero value.

</p></details>

## Examples

### Exporter Configuration

```bash
$ TWELVEDATA_API_KEY="your-twelvedata-api-token" ./twelvedata-exporter
INFO[0000] Starting the Twelvedata exporter on 0.0.0.0:10016
```

### Prometheus Configuration

There are several Prometheus configuration examples provided below:

- **Example Job**: Add from [`prometheus.sample.yml`](./prometheus.sample.yml) to your Prometheus.
- **Example Recording Rules**: Add from [`prometheus.rules.sample.yml`](./prometheus.rules.sample.yml) to your Prometheus.

## Documentation

The reference page under [`docs/`](docs/) carries the behaviour behind the metrics above.

- **[Architecture](docs/architecture.md)** – the scrape path, the absence rules and others.

## Contributing

See [`CONTRIBUTING.md`](CONTRIBUTING.md) for the development setup, test conventions and others.

## License

MIT. The binary statically links Apache-2.0, MIT and BSD 3-Clause dependencies, whose notices are reproduced in [`NOTICE`](NOTICE) and shipped alongside [`LICENSE`](LICENSE) in every release archive and container image.
