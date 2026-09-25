package internal

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// TestNewTwelvedataClient checks the client points at the public API.
func TestNewTwelvedataClient(t *testing.T) {
	t.Parallel()

	client := NewTwelvedataClient("secret")

	if client.baseURL != "https://api.twelvedata.com" {
		t.Errorf("baseURL = %q, want %q", client.baseURL, "https://api.twelvedata.com")
	}

	if client.apiKey != "secret" {
		t.Errorf("apiKey = %q, want %q", client.apiKey, "secret")
	}
}

// TestGetQuoteDecodesFixture checks every field the collector reads survives the decode.
func TestGetQuoteDecodesFixture(t *testing.T) {
	t.Parallel()

	tests := []struct {
		symbol   string
		name     string
		exchange string
		currency string
		previous string
		change   string
	}{
		{spySymbol, "State Street SPDR S&P 500 ETF Trust", "NYSE", "USD", "767.17999", "2.28501"},
		{amznSymbol, "Amazon.com Inc.", "NASDAQ", "USD", "249.38000", "-0.99000488"},
	}

	client := newFixtureClient(t, map[string]string{
		spySymbol:  "quote_spy.json",
		amznSymbol: "quote_amzn.json",
	})

	for _, tt := range tests {
		t.Run(tt.symbol, func(t *testing.T) {
			t.Parallel()

			quote, err := client.GetQuote(tt.symbol)
			if err != nil {
				t.Fatalf("GetQuote(%q) error = %v, want nil", tt.symbol, err)
			}

			got := []string{quote.Symbol, quote.Name, quote.Exchange, quote.Currency, quote.PreviousClose, quote.Change}
			want := []string{tt.symbol, tt.name, tt.exchange, tt.currency, tt.previous, tt.change}

			for i, field := range got {
				if field != want[i] {
					t.Errorf("field %d = %q, want %q", i, field, want[i])
				}
			}
		})
	}
}

// TestGetQuoteSendsKeyInHeaderOnly checks the credential never reaches the query string.
func TestGetQuoteSendsKeyInHeaderOnly(t *testing.T) {
	t.Parallel()

	var got *http.Request

	body := readFixture(t, "quote_spy.json")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.Clone(r.Context())
		_, _ = w.Write(body)
	}))
	defer srv.Close()

	client := NewTwelvedataClient("secret")
	client.baseURL = srv.URL

	if _, err := client.GetQuote(spySymbol); err != nil {
		t.Fatalf("GetQuote() error = %v, want nil", err)
	}

	if auth := got.Header.Get("Authorization"); auth != "apikey secret" {
		t.Errorf("Authorization = %q, want %q", auth, "apikey secret")
	}

	if _, ok := got.URL.Query()["apikey"]; ok {
		t.Error("query carries an apikey parameter, want none")
	}
}

// TestGetQuoteEncodesSymbol checks an injected parameter stays inside the symbol value.
// An apikey parameter beats the header upstream, so an unencoded symbol would replace the key.
func TestGetQuoteEncodesSymbol(t *testing.T) {
	t.Parallel()

	var got url.Values

	body := readFixture(t, "quote_spy.json")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.URL.Query()
		_, _ = w.Write(body)
	}))
	defer srv.Close()

	client := NewTwelvedataClient("secret")
	client.baseURL = srv.URL

	const injected = "SPY&apikey=bogus"

	if _, err := client.GetQuote(injected); err != nil {
		t.Fatalf("GetQuote() error = %v, want nil", err)
	}

	if symbols := got["symbol"]; len(symbols) != 1 || symbols[0] != injected {
		t.Errorf("symbol = %v, want [%q]", symbols, injected)
	}

	if _, ok := got["apikey"]; ok {
		t.Error("injected apikey reached the query, want it kept inside the symbol")
	}
}

// TestGetQuoteDoesNotFollowRedirect checks a redirect never carries the credential onward.
func TestGetQuoteDoesNotFollowRedirect(t *testing.T) {
	t.Parallel()

	followed := 0
	body := readFixture(t, "quote_spy.json")

	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		followed++
		_, _ = w.Write(body)
	}))
	defer target.Close()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL, http.StatusFound)
	}))
	defer srv.Close()

	client := NewTwelvedataClient("secret")
	client.baseURL = srv.URL

	quote, _ := client.GetQuote(spySymbol)

	if followed != 0 {
		t.Errorf("redirect target served %d requests, want 0", followed)
	}

	if quote != nil {
		t.Errorf("GetQuote() = %v, want nil", quote)
	}
}

// TestGetQuoteReportsAbsence checks an error payload is reported as an absent quote.
// The upstream returns it with the quote fields empty, so a rejected key and an unknown symbol read alike.
func TestGetQuoteReportsAbsence(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		fixture string
		status  int
	}{
		{"unauthorized", "error_unauthorized.json", http.StatusUnauthorized},
		{"not found", "error_not_found.json", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			body := readFixture(t, tt.fixture)

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.status)
				_, _ = w.Write(body)
			}))
			defer srv.Close()

			client := NewTwelvedataClient("secret")
			client.baseURL = srv.URL

			quote, err := client.GetQuote(spySymbol)
			if quote != nil {
				t.Errorf("GetQuote() = %v, want nil", quote)
			}

			if !errors.Is(err, errQuoteAbsent) {
				t.Errorf("GetQuote() error = %v, want %v", err, errQuoteAbsent)
			}
		})
	}
}

// TestGetQuoteRejectsBrokenResponse checks a body or transport the decoder cannot use returns an error.
func TestGetQuoteRejectsBrokenResponse(t *testing.T) {
	t.Parallel()

	t.Run("malformed json", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte("{"))
		}))
		defer srv.Close()

		client := NewTwelvedataClient("secret")
		client.baseURL = srv.URL

		if _, err := client.GetQuote(spySymbol); err == nil {
			t.Error("GetQuote() error = nil, want a decode error")
		}
	})

	t.Run("unreachable host", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
		client := NewTwelvedataClient("secret")
		client.baseURL = srv.URL
		srv.Close()

		if _, err := client.GetQuote(spySymbol); err == nil {
			t.Error("GetQuote() error = nil, want a transport error")
		}
	})
}

// TestGetQuoteObservesDuration checks a completed query reaches the shared summary.
func TestGetQuoteObservesDuration(t *testing.T) {
	t.Parallel()

	client := newFixtureClient(t, map[string]string{spySymbol: "quote_spy.json"})

	before := querySampleCount(t)

	if _, err := client.GetQuote(spySymbol); err != nil {
		t.Fatalf("GetQuote() error = %v, want nil", err)
	}

	if after := querySampleCount(t); after <= before {
		t.Errorf("sample count = %d, want greater than %d", after, before)
	}
}
