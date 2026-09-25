# Contributing

Thank you for your interest in contributing to the twelvedata-exporter.

Please follow **[the shared contribution guide](https://github.com/umatare5/.github/blob/main/CONTRIBUTING.md)**, which covers:

- **Development** – the tools to install and the order the pre-commit hooks run in.
- **Command** – the `make` targets and what each one does.
- **Testing** – test placement, mutation checks and what a fixture must carry.
- **Documentation** – page ownership, pinned headings and the verbatim `--help` transcript.
- **Release** – the three files a release touches and what a push to `main` triggers.
- **Pull Requests** – the branch, commit and changelog steps, and what never enters a commit.

This page specifies what is particular to this one.

## Development

These points are where this repository departs from the shared defaults.

- **Do not assume every check runs.** Four are path-filtered: govulncheck, markdownlint, Link Check and actionlint.
- **Do not lean on the coverage gate.** No test exists yet, so CI passes at 0 percent coverage until the first one lands.
- **Copy a flag change into the README.** The `--help` block lives there, and nothing regenerates it.

## Testing

No Go test exists yet, so `promtool` lints the sample rules and nothing asserts behavior.

- **Check the rules with `promtool` locally.** `--lint-fatal` is required, because a lint finding otherwise still exits 0.
- **Redirect `baseURL` to an `httptest` server.** It is unexported, so only a test inside `internal` can rewrite it.
- **Probe reachability with the `demo` key.** A few symbols answer and the rest return `401`, so it proves the path alone.
- **Name one symbol per demo request.** A comma-separated list comes back `401` whatever instruments it names.

The `Prometheus Rules` job runs these two commands.

```bash
promtool check config --lint all --lint-fatal prometheus.sample.yml
promtool check rules --lint all --lint-fatal prometheus.rules.sample.yml
```
