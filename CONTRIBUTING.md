# Contributing

The [shared conventions](https://github.com/umatare5/.github/blob/main/CONTRIBUTING.md) carry the tool setup, the `make` targets, the hook order, the release procedure and the pull request rules. This page carries what is specific to this exporter.

## Development

CI runs Format and Lint, Test and Build, Coverage, Prometheus Rules and CodeQL on every pull request, and govulncheck, actionlint, markdownlint and Link Check when the paths each one watches change.

- **The coverage threshold is zero** — no package carries a test yet, so the job publishes a figure rather than gating on one, and it tightens with the first test that lands.
- **The image declares port 10016** — `EXPOSE` publishes nothing, so `docker run -p` does.

Two commands reproduce the `Prometheus Rules` job locally.

```bash
promtool check rules --lint all --lint-fatal prometheus.rules.sample.yml
promtool check config --lint all --lint-fatal prometheus.sample.yml
```

Both carry `--lint-fatal` because `promtool` otherwise prints a lint finding and still exits 0, so the job would pass over a rule it had just faulted.

## Testing

- **The tree carries no test** — `make test-unit` finds no `*_test.go` anywhere, so the first one added sets its package's conventions rather than following them.
- **The upstream is reachable only from inside** — `baseURL` is unexported, so a test in `internal` points it at an `httptest` server while one outside the package cannot.

## Code Style

The metric names, HELP strings, types and labels are the contract a Prometheus configuration is written against, so changing one breaks the alerts and dashboards built on it. A change to any of them is SemVer-signalled and ships with its own CHANGELOG entry.

## Documentation

Every fact has one page that owns it, and the other pages link to it rather than restating it.

| Page             | Owns                                   |
| :--------------- | :------------------------------------- |
| `README.md`      | What it is, how to run and scrape it   |
| `docs/health.md` | The exporter's own health series       |
| `docs/help.md`   | The verbatim `--help` transcript       |
| `AGENTS.md`      | The upstream behaviours a change holds |

A sentence in `AGENTS.md` about what the API returns is written only after a live reply showed it, because the documentation and the live service disagree. The documented `AAPL` example gives `mic_code` `XNAS` where a live reply gives `XNGS`, and a changed value renames every series that carries it.
