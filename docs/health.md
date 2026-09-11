# Exporter Health

This is the whole set of series the exporter publishes about itself, and none of them takes a flag. The quote series and their labels are in [Collectors](collectors.md).

## Metrics

| Metric                              | Type    | Description                                         |
| :---------------------------------- | :------ | :-------------------------------------------------- |
| `twelvedata_queries_total`          | Counter | Count of completed queries                          |
| `twelvedata_failed_queries_total`   | Counter | Count of failed queries                             |
| `twelvedata_query_duration_seconds` | Summary | Duration of queries to the upstream API             |
| `twelvedata_http_requests_total`    | Counter | Declared with the quote labels, never given a child |

## Labels

None of these series carries a label. `twelvedata_http_requests_total` declares the four quote labels and never receives a value for them.

## Specifications

Each entry carries what the HELP text and the shared [Counter Semantics](README.md#counter-semantics) rules do not.

**`twelvedata_queries_total`**

it increments once where a scrape enters the collector rather than once per upstream call, so the request rate the account is billed for is this rate multiplied by the symbol count.

**`twelvedata_failed_queries_total`**

it is declared and registered but never incremented, so it stays `0` whatever the upstream answers and any alert or dashboard built on it is dead.

- A failed symbol is visible only as the absence of its quote series.

**`twelvedata_query_duration_seconds`**

it observes an instant against itself rather than an elapsed interval, so `_sum` accumulates nanosecond-scale values instead of latency while `_count` follows the quotes that parsed.

- The observation sits after the response is parsed, so a request that failed or returned a nameless quote reaches neither `_count` nor `_sum`.
- Read `_count` as a success count and take latency from the scrape duration instead.

**The registry `/price` serves**

it is built fresh for each scrape and carries the quote series alongside these four, so they are a part of it rather than the whole of it.

- The process and Go collectors go on a registry no handler serves, so no `go_` series appear.
- `twelvedata_http_requests_total` gets no child, so it publishes nothing rather than zero.
