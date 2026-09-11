# Help

The `twelvedata-exporter --help` text, transcribed from the binary.

```text
NAME:
   twelvedata-exporter - Fetch quotes from Twelvedata API

USAGE:
   twelvedata-exporter COMMAND [options...]

VERSION:
   1.2.1

GLOBAL OPTIONS:
   --web.listen-address string, -I string  Set IP address (default: "0.0.0.0")
   --web.listen-port int, -P int           Set port number (default: 10016)
   --web.scrape-path string, -p string     Set the path to expose metrics (default: "/price")
   --twelvedata.api-key string, -a string  Set key to use twelvedata API [$TWELVEDATA_API_KEY]
   --help, -h                              show help
   --version, -v                           print the version
```

## Notes

- `TWELVEDATA_API_KEY` fills `--twelvedata.api-key` and is the only variable the binary reads.
- Either form carries the same value, but the flag reaches the process table where every account on the host can read it, so prefer the variable.
- An unset key stops start-up as a required flag; an empty one stops it in validation.
- `--web.scrape-path` moves the metrics endpoint without changing what it needs, so a relocated path still answers only to a `symbols` query.
- The path is unvalidated: a value with no leading `/`, or `/` itself, panics at start-up.
- `--web.listen-address` is printed verbatim in the landing page's example links, so the `0.0.0.0` default renders links that are not usable as shown.
