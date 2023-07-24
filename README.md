# twelvedata-exporter

![](https://github.com/umatare5/twelvedata-exporter/workflows/Go/badge.svg)

This is a simple stock and funds quotes exporter for
[prometheus](http://prometheus.io). This exporter allows a prometheus instance
to monitor prices of stocks, ETFs, and mutual funds, possibly alerting the user
on any desirable condition (note: prometheus configuration not covered here.)

## Data Provider Setup

This project uses the [stonks page](https://stonks.scd31.com) to fetch stock
price information. This method **does not** support Mutual Funds, but avoids
the hassle of having to create an API key and quota issues of most financial
API providers.

The program is smart enough to "memoize" calls to the financial data provider
and by default caches quotes for 10m. This should reduce the load on the
finance servers, as prometheus tends to scrape exporters on short time
intervals.

## Building the exporter

To build the exporter, you need a relatively recent version of the [Go
compiler](http://golang.org). Download and install the Go compiler and type the
following commands to download, compile, and install the twelvedata-exporter binary
to `/usr/local/bin`:

```bash
OLDGOPATH="$GOPATH"
export GOPATH="/tmp/tempgo"
go get -u -t -v github.com/umatare5/twelvedata-exporter
sudo mv $GOPATH/bin/twelvedata-exporter /usr/local/bin
export GOPATH=$OLDGOPATH
rm -rf /tmp/tempgo
```

## Docker image

The repository includes a ready to use `Dockerfile`. To build a new image, type:

```bash
make image
```

Run `docker images` to see the list of images. The new image is named as
$USER/twelvedata-exporter and exports port 9341 to your host.

## Running the exporter

To run the exporter, just type:

```base
twelvedata-exporter
```

The exporter listens on port 9341 by default. You can use the `--port` command-line
flag to change the port number, if necessary.

## Testing

Use your browser to access [localhost:9341](http://localhost:9341). The exporter should display a simple
help page. If that's OK, you can attempt to fetch a stock using something like:

[http://localhost:9341/price?symbols=GOOGL](http://localhost:9341/price?symbols=GOOGL)

The result should be similar to:

```
# HELP twelvedata_stock_price Asset Price.
# TYPE twelvedata_stock_price gauge
twelvedata_stock_price{name="Alphabet Inc.",symbol="GOOGL"} 1333.54
# HELP twelvedata_exporter_failed_queries_total Count of failed queries
# TYPE twelvedata_exporter_failed_queries_total counter
twelvedata_exporter_failed_queries_total 1
# HELP twelvedata_exporter_queries_total Count of completed queries
# TYPE twelvedata_exporter_queries_total counter
twelvedata_exporter_queries_total 5
# HELP twelvedata_exporter_query_duration_seconds Duration of queries to the upstream API
# TYPE twelvedata_exporter_query_duration_seconds summary
twelvedata_exporter_query_duration_seconds_sum 0.000144555
twelvedata_exporter_query_duration_seconds_count 4
```

## Acknowledgements

I started looking around for a prometheus compatible quotes exporter but
couldn't find anything that satisfied my needs. The closest I found was
[Tristan Colgate-McFarlane](https://github.com/tcolgate)'s [yquotes
exporter](https://github.com/tcolgate/ytwelvedata_exporter), which has stopped
working as Yahoo appears to have deprecated the endpoints required to download
stock data. My thanks to Tristan for his code, which served as the initial
template for this project.
