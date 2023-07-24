// Package internal is a server that uses the twelvedata API as its backend.
package internal

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/kofalt/go-memoize"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// These are metrics for the collector itself
	queryDuration = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "twelvedata_query_duration_seconds",
			Help: "Duration of queries to the upstream API",
		},
	)
	queryCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "twelvedata_queries_total",
			Help: "Count of completed queries",
		},
	)
	errorCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "twelvedata_failed_queries_total",
			Help: "Count of failed queries",
		},
	)

	// Cache external API consuming calls for 10 minutes.
	cache *memoize.Memoizer = memoize.NewMemoizer(10*time.Minute, 20*time.Minute)
)

// collector holds data for a prometheus collector.
type collector struct {
	apikey  string
	limit   int
	symbols []string
}

// newCollector returns a new collector object with parsed data from the URL object.
func newCollector(myURL *url.URL, apiKey string, limit int) (collector, error) {
	var symbols []string

	// The typical query is formatted as: ?symbols=AAA,BBB...&symbols=CCC,DDD...
	// We fetch all symbols into a single slice.
	querySymbols, exists := myURL.Query()["symbols"]
	if !exists {
		return collector{}, fmt.Errorf("missing symbols in the query")
	}

	for _, qValue := range querySymbols {
		symbols = append(symbols, strings.Split(qValue, ",")...)
	}

	return collector{apiKey, limit, symbols}, nil
}

// Describe outputs description for prometheus timeseries.
func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	// Must send one description, or the registry panics.
	ch <- prometheus.NewDesc("dummy", "dummy", nil, nil)
}

// Collect retrieves quote data and outputs Prometheus compatible time series on
// the output channel.
func (c *collector) Collect(ch chan<- prometheus.Metric) {
	queryCount.Inc()

	for _, symbol := range c.symbols {
		quote, cached := c.fetchQuoteData(symbol)
		if quote == nil {
			continue
		}

		ls := []string{"symbol", "name", "exchange", "currency"}
		lvs := c.createLabelValues(symbol, quote)

		changedPrice, _ := strconv.ParseFloat(quote.Change, 64)
		changedPercent, _ := strconv.ParseFloat(quote.PercentChange, 64)
		currentVolume, _ := strconv.ParseFloat(quote.Volume, 64)
		previousClosePrice, _ := strconv.ParseFloat(quote.PreviousClose, 64)
		currentPrice := previousClosePrice + changedPrice

		c.logRetrievedData(symbol, cached, currentPrice)

		c.sendPrometheusMetrics(ch, ls, lvs, changedPrice, changedPercent, currentVolume, currentPrice)
	}
}

// fetchQuoteData fetches quote data for a single symbol using the cachedFetcher.
func (c *collector) fetchQuoteData(symbol string) (*Quote, bool) {
	cachedFetcher := func() (interface{}, error) {
		res, err := FetchQuote(symbol, c.apikey)
		if err != nil {
			errorCount.Inc()
			log.Printf("Error looking up %s: %v\n", symbol, err)
			return nil, nil
		}
		return res, nil
	}

	start := time.Now()
	qret, err, cached := cache.Memoize(symbol, cachedFetcher)
	queryDuration.Observe(time.Since(start).Seconds())

	if err != nil {
		errorCount.Inc()
		log.Printf("Error looking up %s: %v\n", symbol, err)
		return nil, false
	}

	quote, ok := qret.(*Quote)
	if !ok {
		errorCount.Inc()
		log.Printf("Invalid quote data for %s: %v\n", symbol, qret)
		return nil, false
	}

	return quote, cached
}

// createLabelValues creates label values for a given symbol and its quote data.
func (c *collector) createLabelValues(symbol string, quote *Quote) []string {
	return []string{symbol, quote.Name, quote.Exchange, quote.Currency}
}

// logRetrievedData logs the retrieved data for a given symbol and its quote data.
func (c *collector) logRetrievedData(symbol string, cached bool, currentPrice float64) {
	cachedMsg := ""
	if cached {
		cachedMsg = " (cached)"
	}

	log.Printf("Retrieved %s%s, price: %f\n", symbol, cachedMsg, currentPrice)

	// Temporary rate-limit
	if c.limit != 0 {
		time.Sleep(60 * time.Second / time.Duration(c.limit))
	}
}

// sendPrometheusMetrics sends Prometheus metrics for a given symbol and its quote data.
func (c *collector) sendPrometheusMetrics(
	ch chan<- prometheus.Metric, ls, lvs []string, changedPrice, changedPercent, currentVolume, currentPrice float64,
) {
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc("twelvedata_stock_change_price", "Changed price since last close price.", ls, nil),
		prometheus.GaugeValue,
		changedPrice,
		lvs...,
	)

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc("twelvedata_stock_change_percent", "Changed percent since last close price.", ls, nil),
		prometheus.GaugeValue,
		changedPercent,
		lvs...,
	)

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc("twelvedata_stock_volume", "Trading volume during the bar.", ls, nil),
		prometheus.GaugeValue,
		currentVolume,
		lvs...,
	)

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc("twelvedata_stock_price", "Real-time or the latest available price.", ls, nil),
		prometheus.GaugeValue,
		currentPrice,
		lvs...,
	)
}
