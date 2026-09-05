# Exporter health

This is the whole set of series the exporter publishes about itself, and none takes a label. The quote series and their labels are in the [README](../README.md#metrics).

## Metrics

| Metric                              | Type    | Description                             |
| :---------------------------------- | :------ | :-------------------------------------- |
| `twelvedata_queries_total`          | Counter | Count of completed queries              |
| `twelvedata_failed_queries_total`   | Counter | Count of failed queries                 |
| `twelvedata_query_duration_seconds` | Summary | Duration of queries to the upstream API |

## Specifications

Each entry carries what the series' HELP text does not.

**`twelvedata_queries_total`**

it increments once where a scrape enters the collector rather than once per upstream call. The request rate the account is billed for is therefore that rate multiplied by the symbol count, not the counter itself.

**`twelvedata_failed_queries_total`**

it is declared and registered but never incremented, so it stays `0` whatever the upstream answers. A failed symbol is visible only as the absence of its quote series.

**`twelvedata_query_duration_seconds`**

it observes `time.Since(time.Now())` rather than an elapsed interval, so `_sum` stays `0` while `_count` follows the quotes that parsed. Read `_count` as a success count and take latency from the scrape duration instead.

- The observation sits after the response is parsed, so a request that failed or returned a nameless quote reaches neither `_count` nor `_sum`.

**the registry `/price` serves**

the three series above are the whole of it. The process and Go runtime collectors are registered on a second registry no handler serves, so no `go_` or `process_` series reach a scrape.

- `twelvedata_http_requests_total` is declared with the quote labels but never given a child, so it publishes nothing either.
