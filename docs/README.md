# Documentation

Reference pages for twelvedata-exporter.

The [README](../README.md) covers getting the exporter running and scraped. These pages carry the catalogues, the rules every series obeys and the flags that reach them.

| Page                        | Focus                                          |
| :-------------------------- | :--------------------------------------------- |
| [Collectors](collectors.md) | The quote series, their labels and meaning     |
| [Health](health.md)         | The exporter's own series and how to read them |
| [Help](help.md)             | Flags and defaults, as `--help` prints         |

## Technical Information

### Scrape Path

- **Symbols travel in the URL** — a scrape names them, and the exporter keeps no list of its own.
- **Repeats concatenate** — `?symbols=A&symbols=B,C` queries all three, in the order given.
- **A repeat is still fetched** — nothing deduplicates the list, so the second request is billed and only its series is dropped at exposition.
- **No cache exists** — every scrape re-fetches every symbol. See [Cost](../SECURITY.md#cost).
- **No reply reports the budget left** — none carries a rate-limit, credit or `Retry-After` header, so the remaining allowance is visible only through `/api_usage`, which spends a credit of its own.
- **Requests run one at a time** — each is capped at ten seconds by the client.
- **A slow upstream can outlast the scrape** — ten seconds apiece and no cap on the whole means seven stalled symbols exceed the server's 60-second write timeout, and the response is cut short.

> [!NOTE]
> See [Credits](https://support.twelvedata.com/en/articles/5615854-credits) for what a plan allows.

### Endpoints

- **`/` answers everything unmatched** — it is a catch-all, so a mistyped path returns the help page.
- **No path returns 404** — only the scrape path is routed, and every other path reaches that page.
- **Any method is accepted** — neither handler restricts one, so `POST` scrapes exactly as `GET` does.
- **Naming no symbol returns 200 and an empty body** — the scrape counts as successful.

### Absence

- **A failed symbol goes absent** — its five series are omitted and the scrape still answers 200.
- **A field that fails to parse reads `0`** — no label separates that from a genuine zero.
- **Alert on `absent(twelvedata_price)`** — never on `up`, which stays 1 through every failure.
- **An upstream error looks like an empty quote** — it arrives as a `code`/`message`/`status` object that decodes into the same struct with every field empty.
- **The status carries the discriminator the client ignores** — a rejected key answers `401` and an unknown symbol `404`, yet the rejection turns on the missing name.

> [!IMPORTANT]
> [Specifications](collectors.md#specifications) names the fields each instrument type drops and what a zero then means.

### Counter Semantics

- **`twelvedata_queries_total` counts scrapes** — it rises once per scrape, not once per symbol.
- **`twelvedata_query_duration_seconds` lags by one** — its observation lands after the gather has read it, so `_count` reports the scrape before.
- **One of the three never moves at all** — [Health](health.md#specifications) names it and says why.
