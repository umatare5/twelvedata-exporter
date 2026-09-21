# Security Policy

Please follow **[the shared security policy](https://github.com/umatare5/.github/blob/main/SECURITY.md)**, which covers:

- **Supported Versions** – only the latest release carries fixes, so reproduce against it.
- **Reporting a Vulnerability** – the private advisory path and what the response promises.
- **What to Include** – the credentials and addresses to redact, and the fields to send.
- **Exposure** – the unauthenticated surface and the container the image ships.
- **Out of Scope** – findings that belong to the monitored system or to an operator's own configuration.

This page specifies what is particular to this one.

## What to Include

Reproduction needs **the collector flags** in force and **the symbols**.

Redact these before reporting, in addition to the credentials the shared policy names.

- **Tokens** – The Twelve Data API key, from a flag, an environment variable, a process listing or a log line
- **Account** – The account the key resolves to, and any credit balance or billing figure tied to it

## Exposure

This exporter holds one credential, **the Twelve Data API key**.

- **Variables** – The environment is narrower than flags, because `/proc/<pid>/environ` opens to the owner and root alone.
- **Flags** – The flag reaches a process table every account on the host reads, making the variable safer.
- **Logs** – Every request to the scrape path is logged whole by the exporter, recording the URI.
- **Queries** – The symbol list lands in the log, alongside whatever else the query carried.
- **Labels** – The labels `symbol`, `name`, `exchange` and `currency` carry the instrument names in plaintext.

## Endpoints

No route authenticates, so the network path the listener sits on is the whole access control. See also [Endpoints](docs/architecture.md#endpoints).

- **Listener** – `--web.listen-address` defaults to `0.0.0.0`, which answers on every interface.
- **Disclosure** – One scrape discloses the whole watch list for that URL.
- **Landing page** – The `/` route prints example symbols compiled into the binary.
- **Isolation** – No operator choice reaches the landing page; it names the listen address and the scrape path alone.

## Ingress Paths

The exporter exposes a listening socket for incoming HTTP scrapes.

- **Reach** – The exporter listens on all interfaces by default, which accepts every host that routes to it.
- **No allowlist** – The exporter filters no sender, which leaves the packet filter or authenticating proxy to enforce it.
- **Restriction** – `--web.listen-address` narrows the bind, and anything finer needs a packet filter or a proxy.

## Egress Paths

The exporter opens outbound connections only to the Twelve Data API.

- **Host** – Every API call strictly egresses to `https://api.twelvedata.com/quote` under a ten-second timeout.
- **Credential** – The client strictly sends the token in the `Authorization` header, keeping it out of URLs and proxy logs.
- **Cost** – The API cost scales inherently with the requested symbols, because `/quote` costs one credit per symbol.
- **Scaling** – The exporter repeats the quote call for each requested symbol, scaling the total outbound request count.

## Out of Scope

- **Origin** – A defect in the Twelve Data service belongs to Twelve Data rather than to this client.
- **Credit** – Credit or billing exhaustion, which the operator's requirements dictate.
