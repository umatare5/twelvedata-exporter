# Contributing

The [shared contribution guide](https://github.com/umatare5/.github/blob/main/CONTRIBUTING.md) covers what every exporter shares. This page carries the rest.

## Development

CI runs Format and Lint, Test and Build, Coverage, Prometheus Rules and CodeQL on every pull request, and govulncheck, actionlint, markdownlint and Link Check when the paths each one watches change.

- **The coverage threshold is zero** — no package carries a test yet, so the job gates on nothing.
- **The job publishes a figure** — it tightens with the first test that lands.
- **The image declares port 10016** — `EXPOSE` publishes nothing, so `docker run -p` does.

Two commands reproduce the `Prometheus Rules` job locally.

```bash
promtool check rules --lint all --lint-fatal prometheus.rules.sample.yml
promtool check config --lint all --lint-fatal prometheus.sample.yml
```

Both carry `--lint-fatal` because `promtool` otherwise prints a lint finding and still exits 0, so the job would pass over a rule it had just faulted.

## Testing

- **The tree carries no test** — `make test-unit` finds no `*_test.go` anywhere.
- **The first one sets conventions** — it has none to follow, so its package takes its shape.
- **The upstream is reachable only from inside** — `baseURL` is unexported.
- **A test in `internal` redirects it** — it points at an `httptest` server, which one outside cannot.
- **`apikey=demo` proves reachability, not behaviour** — a few symbols answer and the rest `401`.
- **A demo request takes one symbol** — a comma-separated list comes back `401` whatever it names.
- **A captured reply stays out of the tree** — it may carry account-level detail.

## Code Style

The metric names, HELP strings, types and labels are the contract a Prometheus configuration is written against, so changing one breaks the alerts and dashboards built on it. A change to any of them is SemVer-signalled and ships with its own CHANGELOG entry.

## Documentation

Every fact has one page that owns it, and the other pages link to it rather than restating it.

| Page                 | Owns                                    |
| :------------------- | :-------------------------------------- |
| `README.md`          | What it is, how to run and scrape it    |
| `docs/README.md`     | The scrape path, endpoints and absence  |
| `docs/collectors.md` | The quote series and their labels       |
| `docs/health.md`     | The exporter's own health series        |
| `docs/help.md`       | The verbatim `--help` transcript        |
| `SECURITY.md`        | The credential, the egress and the cost |
| `AGENTS.md`          | The invariants a change has to hold     |

A sentence about what the API returns is written only after a live reply showed it, because the documentation and the live service disagree. See [Labels](docs/collectors.md#labels) for the drift that has been measured.
