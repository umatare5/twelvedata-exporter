package internal

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

// Symbols the fixtures carry, captured from the live API.
const (
	spySymbol  = "SPY"
	amznSymbol = "AMZN"
)

// readFixture returns a payload captured from the live Twelvedata API.
func readFixture(t *testing.T, name string) []byte {
	t.Helper()

	body, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("os.ReadFile(%q) error = %v, want nil", name, err)
	}

	return body
}

// newFixtureClient serves one fixture per symbol and returns a client aimed at it.
// Every fixture is read before the server starts, because the handler runs on a goroutine that must not fail the test.
func newFixtureClient(t *testing.T, fixtures map[string]string) *TwelvedataClient {
	t.Helper()

	bodies := make(map[string][]byte, len(fixtures))
	for symbol, name := range fixtures {
		bodies[symbol] = readFixture(t, name)
	}

	absent := readFixture(t, "error_not_found.json")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, ok := bodies[r.URL.Query().Get("symbol")]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write(absent)

			return
		}

		_, _ = w.Write(body)
	}))
	t.Cleanup(srv.Close)

	client := NewTwelvedataClient("fixture-key")
	client.baseURL = srv.URL

	return client
}

// gather runs one gather and fails the case rather than returning its error.
func gather(t *testing.T, gatherer prometheus.Gatherer) []*dto.MetricFamily {
	t.Helper()

	families, err := gatherer.Gather()
	if err != nil {
		t.Fatalf("Gather() error = %v, want nil", err)
	}

	return families
}

// gaugeValues returns every gauge one symbol published, keyed by metric name.
func gaugeValues(t *testing.T, gatherer prometheus.Gatherer, symbol string) map[string]float64 {
	t.Helper()

	values := map[string]float64{}

	for _, family := range gather(t, gatherer) {
		for _, metric := range family.GetMetric() {
			for _, label := range metric.GetLabel() {
				if label.GetName() == "symbol" && label.GetValue() == symbol {
					values[family.GetName()] = metric.GetGauge().GetValue()
				}
			}
		}
	}

	return values
}

// quoteLabels returns the label set carried by one symbol's price series.
func quoteLabels(t *testing.T, gatherer prometheus.Gatherer, symbol string) map[string]string {
	t.Helper()

	for _, family := range gather(t, gatherer) {
		if family.GetName() != "twelvedata_price" {
			continue
		}

		for _, metric := range family.GetMetric() {
			labels := make(map[string]string, len(metric.GetLabel()))
			for _, label := range metric.GetLabel() {
				labels[label.GetName()] = label.GetValue()
			}

			if labels["symbol"] == symbol {
				return labels
			}
		}
	}

	return map[string]string{}
}

// familyTypes returns the metric type of every family a gather produced.
func familyTypes(t *testing.T, gatherer prometheus.Gatherer) map[string]string {
	t.Helper()

	types := map[string]string{}
	for _, family := range gather(t, gatherer) {
		types[family.GetName()] = family.GetType().String()
	}

	return types
}

// seriesCount reports how many series one family carries.
func seriesCount(t *testing.T, gatherer prometheus.Gatherer, name string) int {
	t.Helper()

	for _, family := range gather(t, gatherer) {
		if family.GetName() == name {
			return len(family.GetMetric())
		}
	}

	return 0
}

// counterValue reads a package-level counter through a registry of its own.
func counterValue(t *testing.T, collector prometheus.Collector, name string) float64 {
	t.Helper()

	reg := prometheus.NewRegistry()
	reg.MustRegister(collector)

	for _, family := range gather(t, reg) {
		if family.GetName() == name {
			return family.GetMetric()[0].GetCounter().GetValue()
		}
	}

	t.Fatalf("%s is absent from the registry", name)

	return 0
}

// querySampleCount reports how many durations the shared summary has observed.
func querySampleCount(t *testing.T) uint64 {
	t.Helper()

	reg := prometheus.NewRegistry()
	reg.MustRegister(queryDuration)

	for _, family := range gather(t, reg) {
		if family.GetName() == "twelvedata_query_duration_seconds" {
			return family.GetMetric()[0].GetSummary().GetSampleCount()
		}
	}

	t.Fatal("twelvedata_query_duration_seconds is absent from the registry")

	return 0
}
