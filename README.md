<div align="center">

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="docs/assets/logo_dark.png" width="200px">
  <img src="docs/assets/logo.png" width="200px">
</picture>

  <h1>twelvedata-exporter</h1>

  <p>A third-party Prometheus Exporter for Twelvedata.</p>

  <p>
    <img alt="GitHub Tag" src="https://img.shields.io/github/v/tag/umatare5/twelvedata-exporter?label=Latest%20version" />
    <a href="https://github.com/umatare5/twelvedata-exporter/actions/workflows/go-test-build.yml"><img alt="Test and Build" src="https://github.com/umatare5/twelvedata-exporter/actions/workflows/go-test-build.yml/badge.svg?branch=main" /></a>
    <a href="https://goreportcard.com/badge/github.com/umatare5/twelvedata-exporter"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/umatare5/twelvedata-exporter" /></a>
    <a href="https://pkg.go.dev/github.com/umatare5/twelvedata-exporter@main"><img alt="Go Reference" src="https://pkg.go.dev/badge/umatare5/twelvedata-exporter.svg" /></a>
    <a href="./LICENSE"><img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg" /></a>
  </p>

</div>

## Overview

This exporter allows a prometheus instance to monitor prices of stocks, ETFs, and mutual funds.

> [!Important]
>
> To access the Twelvedata API, you need an access token. Please register with [Twelvedata](https://twelvedata.com/) in advance and generate an access token by referring to [the official document: Getting Started - Authentication](https://twelvedata.com/docs#authentication).

> [!Note]
>
> The Twelvedata API has some limitations based on the license. For example, API limit, accessible market and others. For the limitations, please refer to [twelvedata - Pricing](https://twelvedata.com/pricing) with following documents:
>
> - [Twelvedata Support - Credits](https://support.twelvedata.com/en/articles/5615854-credits)
> - [Twelvedata Support - Stock Exchanges](https://support.twelvedata.com/en/collections/2787973-stock-exchanges)

## Quick Start

```bash
docker run -p 10016:10016 -e TWELVEDATA_API_KEY ghcr.io/umatare5/twelvedata-exporter
```

- `-p`: Publish a container's port `10016/tcp`, to the host `10016/tcp`.
- `-e`: Forward environment variable `TWELVEDATA_API_KEY` into a container.

> [!Tip]
> If you would like to use binaries, please download them from [release page](https://github.com/umatare5/twelvedata-exporter/releases).
>
> - `linux_amd64`, `linux_arm64`, `darwin_amd64`, `darwin_arm64` and `windows_amd64` are supported.

## Syntax

```bash
NAME:
   Fetch quotes from Twelvedata API - twelvedata-exporter

USAGE:
   twelvedata-exporter COMMAND [options...]

VERSION:
   1.0.1

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --web.listen-address value, -I value     Set IP address (default: "0.0.0.0")
   --web.listen-port value, -P value        Set port number (default: 10016)
   --web.scrape-path value, -p value        Set the path to expose metrics (default: "/price")
   --twelvedata.api-key value, -a value     Set key to use twelvedata API [$TWELVEDATA_API_KEY]
   --help, -h                               show help
   --version, -v                            print the version
```

## Configuration

This exporter supports following environment variables:

| Environment Variable | Description                          |
| :------------------- | ------------------------------------ |
| `TWELVEDATA_API_KEY` | The API Key to be used for requests. |

## Metrics

This exporter returns following metrics:

| Metric Name                       | Description                              | Type  | Example Value   |
| --------------------------------- | ---------------------------------------- | ----- | --------------- |
| `twelvedata_change_percent`       | Changed percent since last close price.  | Gauge | `1.00975`       |
| `twelvedata_change_price`         | Changed price since last close price.    | Gauge | `1.72`          |
| `twelvedata_price`                | Real-time or the latest available price. | Gauge | `172.06`        |
| `twelvedata_previous_close_price` | Closing price of the previous day.       | Gauge | `170.34`        |
| `twelvedata_volume`               | Trading volume during the bar.           | Gauge | `1.5206856e+07` |

<details><summary><u>Click to show full metrics</u></summary><p>

```plain
# HELP twelvedata_change_percent Changed percent since last close price.
# TYPE twelvedata_change_percent gauge
twelvedata_change_percent{currency="USD",exchange="NASDAQ",name="Alphabet Inc",symbol="GOOGL"} 1.00975
# HELP twelvedata_change_price Changed price since last close price.
# TYPE twelvedata_change_price gauge
twelvedata_change_price{currency="USD",exchange="NASDAQ",name="Alphabet Inc",symbol="GOOGL"} 1.72
# HELP twelvedata_failed_queries_total Count of failed queries
# TYPE twelvedata_failed_queries_total counter
twelvedata_failed_queries_total 0
# HELP twelvedata_previous_close_price Closing price of the previous day.
# TYPE twelvedata_previous_close_price gauge
twelvedata_previous_close_price{currency="USD",exchange="NASDAQ",name="Alphabet Inc",symbol="GOOGL"} 170.34
# HELP twelvedata_price Real-time or the latest available price.
# TYPE twelvedata_price gauge
twelvedata_price{currency="USD",exchange="NASDAQ",name="Alphabet Inc",symbol="GOOGL"} 172.06
# HELP twelvedata_queries_total Count of completed queries
# TYPE twelvedata_queries_total counter
twelvedata_queries_total 1
# HELP twelvedata_query_duration_seconds Duration of queries to the upstream API
# TYPE twelvedata_query_duration_seconds summary
twelvedata_query_duration_seconds_sum 0
twelvedata_query_duration_seconds_count 0
# HELP twelvedata_volume Trading volume during the bar.
# TYPE twelvedata_volume gauge
twelvedata_volume{currency="USD",exchange="NASDAQ",name="Alphabet Inc",symbol="GOOGL"} 1.5206856e+07
```

</p></details>

## Usage

### Exporter

To refer to the usage, please access http://localhost:10016/ after starting the exporter.

```bash
$ TWELVEDATA_API_KEY="foobarbaz"
$ docker run -p 10016:10016 -e TWELVEDATA_API_KEY ghcr.io/umatare5/twelvedata-exporter
INFO[0000] Listening on port 0.0.0.0:10016
```

or using a binary:

```bash
$ TWELVEDATA_API_KEY="foobarbaz"
$ ./twelvedata-exporter
INFO[0000] Listening on port 0.0.0.0:10016
```

### Prometheus

Please refer to [prometheus.sample.yml#L27-L42](./prometheus.sample.yml#L27-L42).

- To know how to write technical indicators as PromQL, please refer to [prometheus.rules.sample.yml](./prometheus.rules.sample.yml).

## Development

### Build

The repository includes a ready to use `Dockerfile`. Run the following command to build a new image:

```bash
make image
```

The new image is named as `$USER/twelvedata-exporter` and exports `10016/tcp` to your host.

### Release

I'm releasing this exporter manually.

```shell
git tag vX.Y.Z && git push --tags
```

Run the release workflow.

- [GitHub Actions: release workflow](https://github.com/umatare5/twelvedata-exporter/actions/workflows/release.yaml)

## Contribution

1. Fork ([https://github.com/umatare5/twelvedata-exporter/fork](https://github.com/umatare5/twelvedata-exporter/fork))
2. Create a feature branch
3. Commit your changes
4. Rebase your local changes against the master branch
5. Create a new Pull Request

## Licence

[MIT](LICENSE)

## Author

[umatare5](https://github.com/umatare5)

## Acknowledgements

I used to use [Marco Paganini](https://github.com/marcopaganini)'s [quotes-exporter](https://github.com/marcopaganini/quotes-exporter) before. However, due to changes in the external endpoint, that exporter was broken and archived.
Now, I built this exporter taking Marco's exporter as a reference. My thanks to Marco the predecessor, and [Tristan Colgate-McFarlane](https://github.com/tcolgate) the creator of [yquotes-exporter](https://github.com/tcolgate/yquotes_exporter) who preceded Marco.
