# GitHub Copilot Agent Mode – Repo Instructions (umatare5/twelvedata-exporter)

### Scope & Metadata

- **Last Updated**: 2025-08-10

- **Precedence**: Highest in this repository (see §2)

- **Goals**:

  - **Primary Goal**: Contribute **only** to the Prometheus **exporter** in this repo (CLI binary + HTTP metrics endpoint).
  - Keep flags, environment variables, and exported metrics **clear, small, and stable**; maximize **readability/maintainability**.
  - Prefer **minimal diffs** (avoid unnecessary churn).
  - Defaults are secure (no secret logging; HTTP exposure is the operator’s choice and should be safe by default).

- **Non‑Goals**:

  - Creating unrelated services, dashboards, or multi-binary layouts in this repo.
  - Adding scraping targets beyond Twelve Data’s API scope.
  - Emitting or persisting secrets/test credentials.

## 0. Normative Keywords

- NK-001 (**MUST**) Interpret **MUST / MUST NOT / SHOULD / SHOULD NOT / MAY** per RFC 2119/8174.

## 1. Repository Purpose & Scope

- GP-001 (**MUST**) Treat this repository as a **Go Prometheus exporter** for **Twelve Data** price/quote metrics, exposing HTTP metrics for Prometheus to scrape. Baseline behavior, flags, env, and metrics are defined in `README.md`.
- GP-002 (**MUST NOT**) Assume library/SDK usage; this project ships a runnable **binary** (and container image) for Prometheus scraping. Quick start uses Docker (`-p 10016:10016 -e TWELVEDATA_API_KEY …`) and prebuilt binaries.
- GP-003 (**SHOULD**) Align with Twelve Data API basics (auth via API key; plan/limits). Do not hardcode vendor assumptions beyond what `README.md` documents.

## 2. Precedence & Applicability

- GA-001 (**MUST**) When editing/generating code in this repository, Copilot **must follow** this document.
- GA-002 (**MUST**) In this repository, **this file (`copilot-instructions.md`) has the highest precedence** over any other instruction set. **On conflict, always prioritize this file**.
- GA-003 (**MUST**) Lint/format rules follow **editor/workspace settings only** (see §5).
- GA-004 (**MUST**) There is **no separate review instruction**. Review behavior is defined by this file as well.

## 3. Expert Personas (for AI edits/reviews)

- EP-001 (**MUST**) Act as a **Go 1.24 expert**.
- EP-002 (**MUST**) Act as a **Prometheus exporter** expert (exposition format, help/type/labels, stability).
- EP-003 (**SHOULD**) Be familiar with **Twelve Data API** (auth, endpoints, plan limits) to the extent referenced by this repo.

## 4. Security & Privacy

- SP-001 (**MUST NOT**) Log tokens or credentials. **MUST** mask secrets (e.g., `${TOKEN:0:6}...`). Never print `TWELVEDATA_API_KEY` (even on error).
- SP-002 (**MUST**) Handle upstream failures gracefully (429/5xx) without leaking sensitive payloads; use bounded retries/backoff only when compatible with vendor limits.

## 5. Editor‑Driven Tooling (single source of truth)

- ED-001 (**MUST**) Lint/format/type checks follow repository settings (e.g., `.golangci.yml`, `.husky.yaml`, `.markdownlint.json`, `.goimportsignore`).
- ED-002 (**MUST NOT**) Add ad‑hoc flags/rules or inline disables that are not configured.
- ED-003 (**SHOULD**) When rules and reality conflict, propose a **minimal settings PR** instead of local overrides.

## 6. Coding Principles (Basics)

- GC-001 (**MUST**) Apply **KISS/DRY** and keep code quality high.
- GC-002 (**MUST**) Avoid magic numbers; **use named constants** proactively (e.g., default port and path).
- GC-003 (**MUST**) Use **predicate helpers** (e.g., `is*`/`has*`) to improve readability.

## 7. Coding Principles (Conditionals)

- CF-001 (**MUST**) Prefer predicate helpers in conditions.
- CF-002 (**MUST**) Prefer **early returns** to keep logic simple and fast.

## 8. Coding Principles (Loops)

- LP-001 (**MUST**) Prefer **early exits** (`return`/`break`/`continue`) to avoid deep nesting.

## 9. Working Directory / Temp Files

- WD-001 (**MUST**) Place all temporary artifacts (work files, coverage, test binaries, etc.) **under `./tmp`**.
- WD-002 (**MUST**) Before completion, delete **zero‑byte files** (**exception**: keep `.keep`).

## 10. Model‑Aware Execution Workflow (when shell execution is available)

- WF-001 (**MUST**) Before actions: **always launch and use `bash`** (no shell detection/adaptation).
- WF-002 (**MUST**) After editing Go code: run **`go build ./...`**. If a `Makefile` target exists (e.g., `make image` for Docker), prefer it when packaging is touched.
- WF-003 (**MUST**) After editing Go code: run **`go test ./...`**. If no tests exist yet (WIP repo), **skip safely** and do not fabricate tests unless explicitly requested.
- WF-004 (**SHOULD**) If integration env is configured, run documented integration targets; otherwise **skip safely** (do not block).
- WF-005 (**MUST**) After editing **shell scripts** (if any): execute them under `./scripts/` with documented options and ensure dependent `make` targets succeed.
- WF-006 (**MUST**) On completion: summarize actions/results into `./.copilot_reports/<prompt_title>_<YYYY-MM-DD_HH-mm-ss>.md`.

## 11. Tests / Quality Gate (for AI reviewers)

- QG-001 (**MUST**) Keep CI green. Do not merge code that violates configured lint or build. (Tests may be absent while repo is WIP—don’t invent brittle scaffolding.)
- QG-002 (**SHOULD**) Ensure metrics exposition remains stable (names/types/labels) or is **SemVer**‑signaled in releases (see §13).

## 12. Change Scope & Tone (for AI reviewers)

- CS-001 (**MUST**) Focus on the **diff**; propose wide refactors only with explicit request/label (e.g., `allow-wide`).
- CS-002 (**SHOULD**) Tag comments with **\[BLOCKER] / \[MAJOR] / \[MINOR (Nit)] / \[QUESTION] / \[PRAISE]**.
- CS-003 (**SHOULD**) Structure comments as **“TL;DR → Evidence (rule/proof) → Minimal‑diff proposal”**.

## 13. Quick Checklist (before completion)

- QC-001 (**MUST**, **v1.0.0+**) Any **metrics changes** (name/type/labels), CLI flags, or default listen path/port are SemVer‑appropriate and documented in `README.md`.
- QC-002 (**MUST**) Follow **`./README.md`** for baseline requirements (API key, port/path, metrics list, usage).
- QC-003 (**MUST**) If/when `./docs/SECURITY.md` is added, follow it; otherwise follow §4. (At present, no explicit security policy file.)
- QC-004 (**MUST**) Run required build targets (Go build, and **`make image`** if Docker packaging changed).
- QC-005 (**MUST**) Lint/format are clean per editor settings (no ad‑hoc flags/inline disables).
- QC-006 (**MUST**) Temp artifacts are under `./tmp`, zero‑byte files removed, and the report is written to `./.copilot_reports/`.

---

### Appendix A — Repository Baseline (for agents)

- **Binary role**: Prometheus **exporter** for quotes via Twelve Data. Exposes metrics such as:

  - `twelvedata_price` (gauge)
  - `twelvedata_change_price` (gauge)
  - `twelvedata_change_percent` (gauge)
  - `twelvedata_previous_close_price` (gauge)
  - `twelvedata_volume` (gauge)
  - internal telemetry: `twelvedata_queries_total` (counter), `twelvedata_failed_queries_total` (counter), `twelvedata_query_duration_seconds` (summary)
  - **Do not rename** without SemVer bump and docs updates.

- **Listen & scrape** (defaults):

  - Address: `0.0.0.0`
  - Port: `10016`
  - Path: `/price`
  - Flags (and env):

    - `--web.listen-address` (`-I`)
    - `--web.listen-port` (`-P`)
    - `--web.scrape-path` (`-p`)
    - `--twelvedata.api-key` (`-a`) — **or** `TWELVEDATA_API_KEY`

- **Usage**:

  - Docker: `docker run -p 10016:10016 -e TWELVEDATA_API_KEY ghcr.io/umatare5/twelvedata-exporter`
  - Local binary: runs on port `10016`; refer to `http://localhost:10016/` for usage output.

- **Samples**:

  - Prometheus job: see `prometheus.sample.yml`.
  - Rule examples: see `prometheus.rules.sample.yml` for indicator‑style PromQL.

- **Development**:

  - Container packaging via `make image` (see `README.md` / `Makefile`).
  - Releases are tagged and published via GitHub Actions workflow.
