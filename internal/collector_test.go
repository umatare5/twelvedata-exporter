package internal

import (
	"maps"
	"slices"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

// spyGauges is what the SPY fixture publishes, with price as previous close plus change.
var spyGauges = map[string]float64{
	"twelvedata_change_percent":       0.29784501,
	"twelvedata_change_price":         2.28501,
	"twelvedata_previous_close_price": 767.17999,
	"twelvedata_price":                769.465,
	"twelvedata_volume":               86217,
}

// gatherCollector registers a collector alone and returns the gatherer.
// The registry is pedantic, so a series Collect publishes without a matching Describe fails the gather.
func gatherCollector(t *testing.T, client *TwelvedataClient, symbols ...string) prometheus.Gatherer {
	t.Helper()

	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(newCollector(client, symbols))

	return reg
}

// TestCollectorPublishesQuote checks an ETF fixture reaches every series with its real values.
func TestCollectorPublishesQuote(t *testing.T) {
	t.Parallel()

	reg := gatherCollector(t, newFixtureClient(t, map[string]string{spySymbol: "quote_spy.json"}), spySymbol)
	got := gaugeValues(t, reg, spySymbol)

	if !maps.Equal(got, spyGauges) {
		t.Errorf("gauges = %v, want %v", got, spyGauges)
	}
}

// TestCollectorLabelsQuote checks every series carries the four labels that identify the instrument.
// A second symbol is registered, because a label set must come from one series rather than from whichever sorted first.
func TestCollectorLabelsQuote(t *testing.T) {
	t.Parallel()

	client := newFixtureClient(t, map[string]string{
		spySymbol:  "quote_spy.json",
		amznSymbol: "quote_amzn.json",
	})
	reg := gatherCollector(t, client, spySymbol, amznSymbol)

	want := map[string]string{
		"symbol":   spySymbol,
		"name":     "State Street SPDR S&P 500 ETF Trust",
		"exchange": "NYSE",
		"currency": "USD",
	}

	if got := quoteLabels(t, reg, spySymbol); !maps.Equal(got, want) {
		t.Errorf("labels = %v, want %v", got, want)
	}

	if got := quoteLabels(t, reg, amznSymbol)["name"]; got != "Amazon.com Inc." {
		t.Errorf("AMZN name = %q, want %q", got, "Amazon.com Inc.")
	}
}

// TestCollectorPublishesGauges checks a bar is a running aggregate rather than a counter.
func TestCollectorPublishesGauges(t *testing.T) {
	t.Parallel()

	reg := gatherCollector(t, newFixtureClient(t, map[string]string{spySymbol: "quote_spy.json"}), spySymbol)
	types := familyTypes(t, reg)

	for _, name := range slices.Sorted(maps.Keys(spyGauges)) {
		if types[name] != "GAUGE" {
			t.Errorf("%s type = %q, want %q", name, types[name], "GAUGE")
		}
	}
}

// TestCollectorPublishesEverySymbol checks one scrape covers a stock and an ETF.
func TestCollectorPublishesEverySymbol(t *testing.T) {
	t.Parallel()

	client := newFixtureClient(t, map[string]string{
		spySymbol:  "quote_spy.json",
		amznSymbol: "quote_amzn.json",
	})
	reg := gatherCollector(t, client, spySymbol, amznSymbol)

	if got := seriesCount(t, reg, "twelvedata_price"); got != 2 {
		t.Errorf("twelvedata_price series = %d, want 2", got)
	}

	if got := gaugeValues(t, reg, amznSymbol)["twelvedata_price"]; got != 248.38999512 {
		t.Errorf("AMZN price = %v, want 248.38999512", got)
	}
}

// TestCollectorSkipsAbsentSymbol checks a rejected symbol publishes nothing rather than a zero.
func TestCollectorSkipsAbsentSymbol(t *testing.T) {
	t.Parallel()

	client := newFixtureClient(t, map[string]string{spySymbol: "quote_spy.json"})
	reg := gatherCollector(t, client, spySymbol, "NOTAREALSYMBOL")

	if got := seriesCount(t, reg, "twelvedata_price"); got != 1 {
		t.Errorf("twelvedata_price series = %d, want 1", got)
	}
}

// TestCollectorCountsScrape checks each Collect marks one query.
func TestCollectorCountsScrape(t *testing.T) {
	t.Parallel()

	client := newFixtureClient(t, map[string]string{spySymbol: "quote_spy.json"})

	before := counterValue(t, queryCount, "twelvedata_queries_total")
	gather(t, gatherCollector(t, client, spySymbol))

	if after := counterValue(t, queryCount, "twelvedata_queries_total"); after <= before {
		t.Errorf("queryCount = %v, want greater than %v", after, before)
	}
}

// TestCollectorDescribe checks Describe announces every series Collect publishes.
func TestCollectorDescribe(t *testing.T) {
	t.Parallel()

	ch := make(chan *prometheus.Desc, 16)
	newCollector(nil, nil).Describe(ch)
	close(ch)

	announced := make([]string, 0, cap(ch))
	for desc := range ch {
		announced = append(announced, desc.String())
	}

	want := append(slices.Sorted(maps.Keys(spyGauges)), "twelvedata_http_requests_total")
	if len(announced) != len(want) {
		t.Errorf("Describe() = %d descriptions, want %d", len(announced), len(want))
	}

	for _, name := range want {
		if !announces(announced, name) {
			t.Errorf("Describe() omits %s", name)
		}
	}
}

// announces reports whether any description carries the fully qualified name.
func announces(descriptions []string, name string) bool {
	return slices.ContainsFunc(descriptions, func(desc string) bool {
		return strings.Contains(desc, `fqName: "`+name+`"`)
	})
}

// TestParseFloatOrZero checks a field the upstream omits or malforms reads as zero.
func TestParseFloatOrZero(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
		want  float64
	}{
		{"price", "769.465", 769.465},
		{"negative change", "-0.99000488", -0.99000488},
		{"volume", "86217", 86217},
		{"missing field", "", 0},
		{"malformed", "n/a", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := parseFloatOrZero(tt.value); got != tt.want {
				t.Errorf("parseFloatOrZero(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}
