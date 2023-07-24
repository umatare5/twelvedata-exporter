// Package internal is a server that uses the twelvedata API as its backend.
package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	twelvedataQuoteURL = "https://api.twelvedata.com/quote?symbol=%s&apikey=%s"
)

// Quote is a response from twelvedata Quote endpoint
type Quote struct {
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

// FetchQuote returns the current value of a symbol.
func FetchQuote(symbol, apikey string) (*Quote, error) {
	symbol = strings.ToUpper(symbol)

	resp, err := http.Get(fmt.Sprintf(twelvedataQuoteURL, symbol, apikey))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data Quote
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil, err
	}

	if data.Name == "" {
		fmt.Println("Name is not included in JSON:", err)
		return nil, err
	}

	return &data, nil
}
