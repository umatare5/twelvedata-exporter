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

`--twelvedata.api-key` and `TWELVEDATA_API_KEY` carry the same value, but the flag reaches the process table where every account on the host reads it. The exporter exits at start-up when neither is set.

`--web.scrape-path` moves the metrics endpoint without changing what it needs, so a relocated path still answers only to a `symbols` query.
