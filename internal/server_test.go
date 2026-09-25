package internal

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/umatare5/twelvedata-exporter/config"
)

// newFixtureServer returns a server whose client answers from the quote fixtures.
func newFixtureServer(t *testing.T) *Server {
	t.Helper()

	cfg := config.Config{
		WebListenAddress: "0.0.0.0",
		WebListenPort:    10016,
		WebScrapePath:    "/price",
		TwelvedataAPIKey: "fixture-key",
	}

	srv, err := NewServer(&cfg)
	if err != nil {
		t.Fatalf("NewServer() error = %v, want nil", err)
	}

	srv.Client = newFixtureClient(t, map[string]string{
		spySymbol:  "quote_spy.json",
		amznSymbol: "quote_amzn.json",
	})

	return &srv
}

// scrape runs one request against the price handler and returns the exposition.
func scrape(t *testing.T, target string) string {
	t.Helper()

	rec := httptest.NewRecorder()
	newFixtureServer(t).priceHandler(rec, httptest.NewRequest(http.MethodGet, target, http.NoBody))

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	return rec.Body.String()
}

// TestNewServer checks the server holds the configuration it was given.
func TestNewServer(t *testing.T) {
	t.Parallel()

	cfg := config.Config{WebListenAddress: "127.0.0.1", WebListenPort: 9000, TwelvedataAPIKey: "secret"}

	srv, err := NewServer(&cfg)
	if err != nil {
		t.Fatalf("NewServer() error = %v, want nil", err)
	}

	if srv.Config != &cfg {
		t.Error("Config is not the one passed in")
	}

	if srv.Client.apiKey != "secret" {
		t.Errorf("apiKey = %q, want %q", srv.Client.apiKey, "secret")
	}
}

// TestPriceHandlerSplitsSymbols checks both query forms reach the same two symbols.
func TestPriceHandlerSplitsSymbols(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		target string
	}{
		{"comma separated", "/price?symbols=SPY,AMZN"},
		{"repeated parameter", "/price?symbols=SPY&symbols=AMZN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			body := scrape(t, tt.target)

			for _, symbol := range []string{spySymbol, amznSymbol} {
				if !strings.Contains(body, `symbol="`+symbol+`"`) {
					t.Errorf("exposition omits %s", symbol)
				}
			}
		})
	}
}

// TestPriceHandlerPublishesPrice checks a scrape carries the price and the collector's own counters.
func TestPriceHandlerPublishesPrice(t *testing.T) {
	t.Parallel()

	body := scrape(t, "/price?symbols=SPY")

	for _, want := range []string{
		`twelvedata_price{currency="USD",exchange="NYSE",`,
		"769.465",
		"twelvedata_queries_total",
		"twelvedata_query_duration_seconds_count",
		"twelvedata_failed_queries_total",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("exposition omits %q", want)
		}
	}
}

// TestPriceHandlerIgnoresEmptyQuery checks a scrape naming no symbol spends no credit.
func TestPriceHandlerIgnoresEmptyQuery(t *testing.T) {
	t.Parallel()

	if body := scrape(t, "/price"); body != "" {
		t.Errorf("body = %q, want empty", body)
	}
}

// TestHelpDescribesScrapeURL checks the root page states the address the exporter answers on.
func TestHelpDescribesScrapeURL(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	newFixtureServer(t).help(rec, httptest.NewRequest(http.MethodGet, "/", http.NoBody))

	for _, want := range []string{
		"<h1>Prometheus Twelvedata Exporter</h1>",
		"http://0.0.0.0:10016/price?symbols=GOOGL",
		"http://0.0.0.0:10016/price?symbols=AMZN,AAPL,MSFT",
	} {
		if !strings.Contains(rec.Body.String(), want) {
			t.Errorf("help page omits %q", want)
		}
	}
}

// TestServerStartReportsBindFailure checks a port the listener rejects returns rather than halting the process.
func TestServerStartReportsBindFailure(t *testing.T) {
	t.Parallel()

	cfg := config.Config{
		WebListenAddress: "127.0.0.1",
		WebListenPort:    -1,
		WebScrapePath:    "/price",
		TwelvedataAPIKey: "fixture-key",
	}

	srv, err := NewServer(&cfg)
	if err != nil {
		t.Fatalf("NewServer() error = %v, want nil", err)
	}

	if err := srv.Start(); err == nil {
		t.Error("Start() error = nil, want a bind error")
	}
}
