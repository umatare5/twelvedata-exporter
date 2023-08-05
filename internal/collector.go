// Package internal contains the implementation of this exporter.
package internal

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/umatare5/twelvedata-exporter/log"
)

const (
	namespace = "twelvedata"
)

// Metrics descriptions
var (
	change_price = prometheus.NewDesc( //nolint:golint,revive
		prometheus.BuildFQName(namespace, "", "change_price"),
		"Changed price since last close price.",
		[]string{"symbol", "name", "exchange", "currency"}, nil,
	)

	change_percent = prometheus.NewDesc( //nolint:golint,revive
		prometheus.BuildFQName(namespace, "", "change_percent"),
		"Changed percent since last close price.",
		[]string{"symbol", "name", "exchange", "currency"}, nil,
	)

	volume = prometheus.NewDesc( //nolint:golint,revive
		prometheus.BuildFQName(namespace, "", "volume"),
		"Trading volume during the bar.",
		[]string{"symbol", "name", "exchange", "currency"}, nil,
	)

	previous_close_price = prometheus.NewDesc( //nolint:golint,revive
		prometheus.BuildFQName(namespace, "", "previous_close_price"),
		"Closing price of the previous day.",
		[]string{"symbol", "name", "exchange", "currency"}, nil,
	)

	price = prometheus.NewDesc( //nolint:golint,revive
		prometheus.BuildFQName(namespace, "", "price"),
		"Real-time or the latest available price.",
		[]string{"symbol", "name", "exchange", "currency"}, nil,
	)

	httpRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "http_requests_total",
		Help:      "The total number of requests labeled by response code",
	},
		[]string{"symbol", "name", "exchange", "currency"},
	)
)

// Collector collects Quote Metrics
type Collector struct {
	client  *TwelvedataClient
	symbols []string
}

// newCollector returns an initialized exporter
func newCollector(client *TwelvedataClient, symbols []string) *Collector {
	return &Collector{
		client:  client,
		symbols: symbols,
	}
}

// Describe outputs description for prometheus timeseries.
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- change_price
	ch <- change_percent
	ch <- volume
	ch <- price
	httpRequestsTotal.Describe(ch)
}

// Collect retrieves quote data and outputs Prometheus compatible time series
// on the output channel.
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	queryCount.Inc()

	for _, symbol := range c.symbols {
		quote, _ := c.client.GetQuote(symbol)
		if quote == nil {
			continue
		}

		c.processMetrics(quote, ch)
	}
}

func (c *Collector) processMetrics(quote *QuoteResponse, ch chan<- prometheus.Metric) {
	isCached := false

	labels := c.createLabelValues(quote.Symbol, quote)
	changedPrice, _ := strconv.ParseFloat(quote.Change, 64)
	changedPercent, _ := strconv.ParseFloat(quote.PercentChange, 64)
	currentVolume, _ := strconv.ParseFloat(quote.Volume, 64)
	previousClosePrice, _ := strconv.ParseFloat(quote.PreviousClose, 64)

	ch <- prometheus.MustNewConstMetric(change_price, prometheus.GaugeValue, changedPrice, labels...)
	ch <- prometheus.MustNewConstMetric(change_percent, prometheus.GaugeValue, changedPercent, labels...)
	ch <- prometheus.MustNewConstMetric(volume, prometheus.GaugeValue, currentVolume, labels...)
	ch <- prometheus.MustNewConstMetric(previous_close_price, prometheus.GaugeValue, previousClosePrice, labels...)
	ch <- prometheus.MustNewConstMetric(price, prometheus.GaugeValue, previousClosePrice+changedPrice, labels...)

	httpRequestsTotal.Collect(ch)

	// TODO: Implement caching. isCached is always false.
	c.logRetrievedData(quote.Symbol, isCached, previousClosePrice+changedPrice)
}

// createLabelValues creates label values for a given symbol and its quote data.
func (c *Collector) createLabelValues(symbol string, quote *QuoteResponse) []string {
	return []string{symbol, quote.Name, quote.Exchange, quote.Currency}
}

// logRetrievedData logs the retrieved data for a given symbol and its quote data.
func (c *Collector) logRetrievedData(symbol string, cached bool, currentPrice float64) {
	cachedMsg := ""
	if cached {
		cachedMsg = " (cached)"
	}

	log.Infof("Retrieved %s%s, price: %f\n", symbol, cachedMsg, currentPrice)
}
