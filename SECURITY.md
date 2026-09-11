# Security Policy

The [shared security policy](https://github.com/umatare5/.github/blob/main/SECURITY.md) covers what every exporter here shares. This page carries the rest.

## What to Include

Redact these before reporting, in addition to the credentials the shared policy names.

- The Twelve Data API key, from a flag, an environment variable, a process listing or a log line
- The account it resolves to, and any credit balance or billing figure tied to that account
- The `symbols` list of a private scrape URL, which names the instruments being watched

A scrape URL carries `symbols` alone, so reproduction needs it together with the flags in force.

## Exposure

This exporter holds one credential, and `--twelvedata.api-key` and `TWELVEDATA_API_KEY` are its only sources, so nothing arriving in a scrape can set it. The scrape path answers a `symbols` query and reads no other parameter.

- **The variable is safer** — the flag reaches a process table every account on the host reads.
- **The environment is narrower** — `/proc/<pid>/environ` opens to the owner and root alone.
- **The log records the URI** — every request to the scrape path is logged whole by the exporter.
- **The symbol list lands there** — so does whatever else the query carried.
- **The labels name the instrument** — `symbol`, `name`, `exchange` and `currency` carry it.
- **The endpoint is unauthenticated** — one scrape discloses the whole watch list.
- **The landing page is fixed** — `/` prints example symbols compiled into the binary.
- **No operator choice reaches it** — the page names the listen address and the scrape path alone.

> [!IMPORTANT]
> The key reaching a log line, the landing page or a scrape response body is a vulnerability, as is a `symbols` value escaping the URL it arrived in into a metric name or a label name.
>
> So is a request reaching a host other than `api.twelvedata.com`, or an endpoint other than `/quote`, since both are compiled in and no flag or query moves them.

## Egress Paths

Each scrape reaches `https://api.twelvedata.com/quote` once per symbol under a ten second timeout. See [Scrape Path](docs/README.md#scrape-path) for what that bounds and what it does not.

### Credential

- **The key travels in the query string** — the client appends `?apikey=` to every request.
- **No header carries it** — the credential sits in the URL rather than in one.
- **The documented form is a header** — `Authorization: apikey <key>` is what Twelve Data documents.
- **The client sends the query form** — [`AGENTS.md`](AGENTS.md) records it as a defect this repository carries.
- **A URL outlives a header** — each proxy and gateway on the path logs the query string.
- **Those logs are not the operator's** — a key sent that way survives in every one of them.

### Cost

- **A scrape spends money** — `/quote` costs one credit per symbol.
- **The plan allowance is per minute** — the symbol count and the scrape interval set the spend.

## Out of Scope

- A defect in the Twelve Data service belongs to Twelve Data rather than to this client.
- Credit or billing exhaustion, which the operator's own symbol count and scrape interval decide.
