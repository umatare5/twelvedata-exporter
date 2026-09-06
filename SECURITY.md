# Security Policy

The [shared policy](https://github.com/umatare5/.github/blob/main/SECURITY.md) carries the supported versions, the reporting channel, what a report must contain and the out-of-scope list. This page carries what is specific to a Twelve Data client.

## What to Include

Redact these before reporting, in addition to the credentials the shared policy names.

- The Twelve Data API key, from a flag, an environment variable, a process listing or a log line
- The account it resolves to, and any credit balance or billing figure tied to that account
- The `symbols` list of a private scrape URL, which names the instruments being watched

A scrape URL carries `symbols` alone, so reproduction needs it together with the flags in force.

## Exposure

This exporter holds one credential, and `--twelvedata.api-key` and `TWELVEDATA_API_KEY` are its only sources, so nothing arriving in a scrape can set it. The scrape path answers a `symbols` query and reads no other parameter.

- **The variable is safer** — the flag reaches the process table, which every account on the host reads, while `/proc/<pid>/environ` opens to the owner and root alone.
- **The log records the URI** — every request to the scrape path is logged whole, so the symbol list lands in the exporter's own log alongside whatever else the query carried.
- **The labels name the instrument** — `symbol`, `name`, `exchange` and `currency` reach an unauthenticated endpoint, so one scrape discloses the whole watch list.
- **The landing page is fixed** — `/` prints example symbols compiled into the binary, so it discloses the listen address and the scrape path and nothing an operator chose.

> [!IMPORTANT]
> The key reaching a log line, the landing page or a scrape response body is a vulnerability, as is a `symbols` value escaping the URL it arrived in into a metric name or a label name.
>
> So is a request reaching a host other than `api.twelvedata.com`, or an endpoint other than `/quote`, since both are compiled in and no flag or query moves them.

## Egress

Each scrape reaches `https://api.twelvedata.com/quote` once per symbol under a ten second timeout, so a stalled upstream ends the request rather than holding the scrape open.

- **The key travels in the query string** — the client appends `?apikey=` to every request, so the credential sits in the URL rather than in a header.
- **The header is the documented form** — [`AGENTS.md`](AGENTS.md) records `Authorization: apikey <key>` as the form Twelve Data documents, and the `?apikey=` form the client sends as a defect this repository still carries.
- **A URL outlives a header** — a query string is written to the access log of every proxy and gateway on the path, so a key sent that way survives in logs the operator does not own.
- **A scrape spends money** — `/quote` costs one credit per symbol, so the symbol count and the scrape interval together set the spend against a per-minute plan allowance.

## Out of Scope

- A defect in the Twelve Data service belongs to Twelve Data rather than to this client.
- Credit or billing exhaustion, which the operator's own symbol count and scrape interval decide.
