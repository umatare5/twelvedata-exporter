# Security Policy

## Supported Versions

No older tag gets a patch branch, so only the most recent tagged release carries fixes — reproduce a finding against it before reporting.

## Reporting a Vulnerability

Report privately through [GitHub Security Advisories](https://github.com/umatare5/twelvedata-exporter/security/advisories/new). **Please do not report a vulnerability through a public GitHub issue or a pull request.**

One maintainer works on this in their own time, so the response is best effort with no promised window. The advisory goes out once the fix ships, because publishing it earlier discloses the flaw while every deployment is still exposed. It carries a CVE request and credits the reporter unless they ask otherwise.

## What to Include

**Redact these first.** Each one names an account or a position rather than a defect, so none of them belongs in a report.

- The Twelve Data API key, which a log line, a process listing or a scrape URL can carry
- The account the key resolves to, and any credit or billing figure tied to that account
- A private scrape URL, whose `symbols` list discloses the instruments being watched

Then include the following:

- **Affected versions** (required): The `twelvedata-exporter` release, and the image tag if you ran the container
- **Reproduction steps** (required): The flags and environment variables, and the scrape URL with the key removed
- **Output** (required): The exposition or the log lines, with every value above removed
- **Impact assessment** (required): The exploit scenario, and what it reaches
- **Suggested fix** (optional): Proposed remediation, if any
- **Disclosure status** (required): Whether it is shared elsewhere, and your plan for sharing it

## Scope

In scope:

- The API key reaching a log line, the landing page, or the body of a scrape response
- The API key reaching an upstream request other than as the documented `apikey` query parameter
- A `symbols` value escaping the URL it arrived in, into a label name or a metric name
- A scrape reaching an endpoint other than `/quote`, or a host other than `api.twelvedata.com`
- The published container image

Out of scope:

- The unauthenticated metrics endpoint, because the operator controls its exposure with a network path
- A dependency advisory with no path reachable from `./cmd` — show the path, or a `govulncheck` finding
- A defect in the Twelve Data API itself, which belongs to its operator rather than to this client
- Credit exhaustion, since the symbol count and the scrape interval that cause it are the operator's own
