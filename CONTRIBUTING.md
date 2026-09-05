# Contributing

Thank you for considering a contribution.

## Commands

The following `make` commands are available for development and testing:

| Command                     | Description                                     |
| :-------------------------- | :---------------------------------------------- |
| `make help`                 | Display available targets and requirements      |
| `make build`                | Build the binary to `./tmp/twelvedata-exporter` |
| `make lint`                 | Run golangci-lint and tidy go.mod               |
| `make test-unit`            | Run unit tests with coverage using gotestsum    |
| `make test-unit-coverage`   | Generate HTML coverage report                   |
| `make clean`                | Remove build artifacts and backup files         |
| `make image`                | Build Docker image                              |
| `make pre-commit-install`   | Install the pre-commit hooks                    |
| `make pre-commit-test`      | Run every hook across the tree                  |
| `make pre-commit-uninstall` | Remove the pre-commit hooks                     |

Markdown style is enforced by the `markdownlint-cli2` hook that `make pre-commit-install` wires in, and again in CI. Links are checked in CI only, because that run reaches third-party hosts, so `lychee .` is what reproduces a link failure locally.

The Prometheus samples are checked in CI too. `promtool check rules` and `promtool check config` run against [`prometheus.rules.sample.yml`](prometheus.rules.sample.yml) and [`prometheus.sample.yml`](prometheus.sample.yml). Both run with `--lint-fatal`, because promtool otherwise prints a lint finding and still exits 0.

## Build

The repository includes a ready to use `Dockerfile`. To build a new Docker image:

```bash
make image
```

This cross-compiles a Linux binary into `./tmp/image/linux/<arch>`, then builds from `./tmp/image` rather than the repository root. The `Dockerfile` expects the GoReleaser context layout, `linux/<arch>/twelvedata-exporter` beside `LICENSE`, which the root does not carry.

The image is tagged `$USER/twelvedata-exporter` and declares port 10016 without publishing it, so publish it with `docker run -p`. Released images are pushed to `ghcr.io/umatare5/twelvedata-exporter` by GoReleaser instead.

## Release

To release a new version, follow these steps:

1. Rename the `## [Unreleased]` section in `CHANGELOG.md` to `## [vX.Y.Z]`, and add that version's release link at the foot of the file.
2. Update the version in the `VERSION` file to match.
3. Submit a pull request with both files.

Merging that pull request starts the release. A push to `main` touching `VERSION` runs the [release workflow](https://github.com/umatare5/twelvedata-exporter/actions/workflows/go-release.yml), which tags the commit and publishes the release in the same run.

## Pull requests

1. [Fork](https://github.com/umatare5/twelvedata-exporter/fork) the repository
2. Create a feature branch
3. Commit your changes
4. Record any change to the metric surface under the `## [Unreleased]` section in `CHANGELOG.md`, adding the section if it is not there yet
5. Rebase your local changes against the `main` branch
6. Create a new Pull Request
