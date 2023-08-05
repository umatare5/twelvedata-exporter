// Package internal contains the implementation of this exporter.
package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/umatare5/twelvedata-exporter/log"
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
)

// TwelvedataGatherer is an interface for Twelvedata API
type TwelvedataGatherer interface {
	GetQuote(symbol string) (float64, error)
}

// TwelvedataClient is a client for Twelvedata API
type TwelvedataClient struct {
	baseURL string
	apiKey  string
}

// QuoteResponse is a response from twelvedata Quote endpoint
type QuoteResponse struct {
	Symbol        string       `json:"symbol"`
	Name          string       `json:"name"`
	Exchange      string       `json:"exchange"`
	MicCode       string       `json:"mic_code"`
	Currency      string       `json:"currency"`
	Datetime      string       `json:"datetime"`
	Timestamp     int          `json:"timestamp"`
	Open          string       `json:"open"`
	High          string       `json:"high"`
	Low           string       `json:"low"`
	Close         string       `json:"close"`
	Volume        string       `json:"volume"`
	PreviousClose string       `json:"previous_close"`
	Change        string       `json:"change"`
	PercentChange string       `json:"percent_change"`
	AverageVolume string       `json:"average_volume"`
	IsMarketOpen  bool         `json:"is_market_open"`
	FiftyTwoWeek  fiftyTwoWeek `json:"fifty_two_week"`
}

type fiftyTwoWeek struct {
	Low               string `json:"low"`
	High              string `json:"high"`
	LowChange         string `json:"low_change"`
	HighChange        string `json:"high_change"`
	LowChangePercent  string `json:"low_change_percent"`
	HighChangePercent string `json:"high_change_percent"`
	Range             string `json:"range"`
}

// NewTwelvedataClient returns Twelvedata Client.
func NewTwelvedataClient(apiKey string) *TwelvedataClient {
	return &TwelvedataClient{
		baseURL: "https://api.twelvedata.com",
		apiKey:  apiKey,
	}
}

// GetQuote sends GET request to Twelvedata API.
func (t *TwelvedataClient) GetQuote(symbol string) (*QuoteResponse, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequestWithContext(
		context.Background(),
		"GET",
		fmt.Sprintf(t.baseURL+"/quote?symbol=%s&apikey=%s", symbol, t.apiKey),
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Error sending request to server: %s", err)
		return nil, err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Errorf("Error closing response body: %s", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error closing response body: %s", err)
		return nil, err
	}

	var data QuoteResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Errorf("Error parsing JSON: %s", err)
		return nil, err
	}

	if data.Name == "" {
		log.Errorf("Name is not included in JSON: %s", err)
		return nil, err
	}

	queryDuration.Observe(time.Since(time.Now()).Seconds())

	return &data, nil
}
