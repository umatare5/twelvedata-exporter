# Collectors

The quote series, published once per resolved symbol per scrape. Nothing here is switchable: the exporter has no `--collector.*` flag, so every scrape publishes all five.

## Metrics

| Metric                            | Type  | Description                                 |
| :-------------------------------- | :---- | :------------------------------------------ |
| `twelvedata_price`                | Gauge | Real-time or the latest available price     |
| `twelvedata_previous_close_price` | Gauge | Closing price of the previous day           |
| `twelvedata_change_price`         | Gauge | Change since the previous close             |
| `twelvedata_change_percent`       | Gauge | Change since the previous close, in percent |
| `twelvedata_volume`               | Gauge | Trading volume during the bar               |

## Labels

All five carry the same four labels, and every value is the upstream reply echoed back.

| Label      | Holds                                       |
| :--------- | :------------------------------------------ |
| `symbol`   | The symbol as the reply spelled it          |
| `name`     | The instrument name the quote carried       |
| `exchange` | The venue name, not its MIC                 |
| `currency` | The currency the row's prices are quoted in |

**`symbol`**

it is the reply's own spelling rather than the one the scrape URL used, so `?symbols=aapl` labels the series `AAPL` and a rule written against the requested case misses it.

**`exchange`**

it carries the venue name where the reply also offers a MIC, and the MIC never becomes a label. A live `AAPL` reply labels `NASDAQ` while its `mic_code` reads `XNGS`, against the `XNAS` the upstream documentation shows.

**`name`**

it is whatever the live reply spells, and the documentation and the live service disagree. The documented `AAPL` example gives `Apple Inc` where a live reply gives `Apple Inc.`, and a changed value renames every series carrying it.

- A cross-listing quotes in its own currency, so ranking across symbols means little.
- A bare symbol resolves to one listing, which `exchange`, `mic_code` and `country` narrow.

**`EUR/USD` and `BTC/USD`**

an instrument outside equities drops the fields the labels need — `currency`, `mic_code`, `volume` and `average_volume` — so either reply leaves `currency` empty and reads a parsed `0` for volume. Both markets also keep `is_market_open` true around the clock.

## Specifications

Each entry carries what the HELP text and the shared [Absence](README.md#absence) rules do not.

**`twelvedata_price`**

it is computed as `previous_close + change` rather than read from the reply's `close`, so it agrees with the two series beside it at every instant and carries the float artifact the addition leaves.

- A live `AAPL` reply whose `close` read `326.57001` published `326.57000999999997`.
- An unparsed previous close collapses this series onto `twelvedata_change_price`.

**`twelvedata_previous_close_price`**

it means the previous bar, and the bar is whatever `interval` names. The exporter sends no `interval`, so it takes the upstream default of one day and the series is the prior session's close.

- At `interval=1min` a live reply gave `previous_close` `319.82999` against `close` `319.98999`, so a different default would silently change what `twelvedata_change_price` measures.

**`twelvedata_volume`**

it counts only what traded inside the current bar, which is why it is a gauge: a counter's reset at each new bar would read as a rollover.

**`twelvedata_change_percent`**

it moves on a corporate action as it moves on a trade. `/time_series` takes an `adjust` parameter defaulting to `splits` where `/quote` takes none, so an unadjusted four-for-one split reads as a −75% move and no threshold on magnitude separates one from a crash.

**`dp`**

it sets the precision of every price field, five places by default, and rounds to decimal places rather than to significant figures, so the fields arrive as strings to be parsed.

- Five places suit a US equity and are coarse for an FX pair, where `dp=2` rounds `change` to `-0.00`.
- Asking for more surfaces the float32 the API stores: `AAPL` at `dp=11` returns `319.97000122070`.
- The `rolling_*` and `extended_*` fields arrive only when the plan and the request ask for them.

**`is_market_open` and `datetime`**

they reach no series. The reply carries both and the exporter publishes neither, while Prometheus stamps each sample with scrape time, so a repeated out-of-session bar is indistinguishable from a live one.

**`prometheus.rules.sample.yml`**

it derives its RSI from scrape samples rather than from bars, and a session repeats one daily bar across every scrape in it, so the period the rule measures is wall-clock. It demonstrates the mechanism rather than carrying a trading signal.

> [!IMPORTANT]
> These names, types and labels are the contract a Prometheus configuration is written against, so a change to any of them is SemVer-signalled and ships with its own [CHANGELOG](../CHANGELOG.md) entry.
