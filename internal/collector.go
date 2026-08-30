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

// quoteLabelNames are the label names attached to every quote metric.
var quoteLabelNames = []string{"symbol", "name", "exchange", "currency"}

// Metrics descriptions.
var (
	changePriceDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "change_price"),
		"Changed price since last close price.",
		quoteLabelNames, nil,
	)

	changePercentDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "change_percent"),
		"Changed percent since last close price.",
		quoteLabelNames, nil,
	)

	volumeDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "volume"),
		"Trading volume during the bar.",
		quoteLabelNames, nil,
	)

	previousClosePriceDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "previous_close_price"),
		"Closing price of the previous day.",
		quoteLabelNames, nil,
	)

	priceDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "price"),
		"Real-time or the latest available price.",
		quoteLabelNames, nil,
	)

	httpRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "http_requests_total",
		Help:      "The total number of requests labeled by response code",
	},
		quoteLabelNames,
	)
)

// Collector collects Quote Metrics.
type Collector struct {
	client  *TwelvedataClient
	symbols []string
}

// newCollector returns an initialized exporter.
func newCollector(client *TwelvedataClient, symbols []string) *Collector {
	return &Collector{
		client:  client,
		symbols: symbols,
	}
}

// Describe outputs description for prometheus timeseries.
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- changePriceDesc
	ch <- changePercentDesc
	ch <- volumeDesc
	ch <- priceDesc
	httpRequestsTotal.Describe(ch)
}

// Collect retrieves quote data and outputs Prometheus compatible time series
// on the output channel.
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	queryCount.Inc()

	for _, symbol := range c.symbols {
		quote, err := c.client.GetQuote(symbol)
		if err != nil || quote == nil {
			continue
		}

		c.processMetrics(quote, ch)
	}
}

func (c *Collector) processMetrics(quote *QuoteResponse, ch chan<- prometheus.Metric) {
	isCached := false

	labels := c.createLabelValues(quote.Symbol, quote)
	changedPrice := parseFloatOrZero(quote.Change)
	changedPercent := parseFloatOrZero(quote.PercentChange)
	currentVolume := parseFloatOrZero(quote.Volume)
	previousClosePrice := parseFloatOrZero(quote.PreviousClose)

	ch <- prometheus.MustNewConstMetric(changePriceDesc, prometheus.GaugeValue, changedPrice, labels...)
	ch <- prometheus.MustNewConstMetric(changePercentDesc, prometheus.GaugeValue, changedPercent, labels...)
	ch <- prometheus.MustNewConstMetric(volumeDesc, prometheus.GaugeValue, currentVolume, labels...)
	ch <- prometheus.MustNewConstMetric(previousClosePriceDesc, prometheus.GaugeValue, previousClosePrice, labels...)
	ch <- prometheus.MustNewConstMetric(priceDesc, prometheus.GaugeValue, previousClosePrice+changedPrice, labels...)

	httpRequestsTotal.Collect(ch)

	// TODO: Implement caching. isCached is always false.
	c.logRetrievedData(quote.Symbol, isCached, previousClosePrice+changedPrice)
}

// parseFloatOrZero converts a numeric string to a float64, returning 0 when parsing fails.
func parseFloatOrZero(value string) float64 {
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}

	return parsed
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
