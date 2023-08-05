// Package internal contains the implementation of this exporter.
package internal

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/umatare5/twelvedata-exporter/config"
	"github.com/umatare5/twelvedata-exporter/log"
)

// Server struct
type Server struct {
	ListenAddrAndPort string
	ScrapePath        string
	Client            *TwelvedataClient
}

// NewServer returns Twelvedata struct
func NewServer(config *config.Config) (Server, error) {
	return Server{
		ListenAddrAndPort: config.WebListenAddress + ":" + strconv.Itoa(config.WebListenPort),
		ScrapePath:        config.WebScrapePath,
		Client:            NewTwelvedataClient(config.TwelvedataAPIKey),
	}, nil
}

// Start starts the server
func (s *Server) Start() {
	reg := prometheus.NewRegistry()

	// Add standard process and Go metrics.
	reg.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(),
	)

	// Register handlers.
	http.HandleFunc("/", s.help)
	http.HandleFunc(s.ScrapePath, func(w http.ResponseWriter, r *http.Request) {
		s.priceHandler(w, r)
	})

	log.Infof("Listening on port %s", s.ListenAddrAndPort)
	srv := &http.Server{
		Addr:         s.ListenAddrAndPort,
		Handler:      nil,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	log.Fatal(srv.ListenAndServe())
}

// priceHandler handles the "/price" endpoint. It creates a new collector with
// the URL and a new prometheus registry to use that collector.
func (s *Server) priceHandler(w http.ResponseWriter, r *http.Request) {
	// The typical query is formatted as: ?symbols=AAA,BBB...&symbols=CCC,DDD...
	// We fetch all symbols into a single slice.
	syms := r.URL.Query()["symbols"]
	if len(syms) == 0 {
		log.Infof("missing symbols in the query: %s", r.RequestURI)
		return
	}
	log.Infof("URL: %s\n", r.RequestURI)

	var symbols []string
	for _, sym := range syms {
		symbols = append(symbols, strings.Split(sym, ",")...)
	}

	registry := prometheus.NewRegistry()

	// These will be collected every time the /stock or /fund endpoint is reached.
	registry.MustRegister(
		newCollector(s.Client, symbols),
		queryCount,
		queryDuration,
		errorCount,
	)

	// Delegate http serving to Promethues client library, which will call collector.Collect.
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

// help returns a help message for those using the root URL.
func (s *Server) help(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "<h1>Prometheus Twelvedta Exporter</h1>")
	fmt.Fprintf(w, "<p>To fetch the price of quotes, your URL must be formatted as:</p>")
	fmt.Fprintf(w, "http://%s/price?symbols=AAAA,BBBB,CCCC", s.ListenAddrAndPort)
	fmt.Fprintf(w, "<p><b>Examples:</b></p>")
	fmt.Fprintf(w, "<ul>")

	symbols := []string{
		"GOOGL",
		"AMZN,AAPL,MSFT",
	}

	for _, symbol := range symbols {
		fmt.Fprintf(w, "<li><a href=\"http://%s/price?symbols=%s\">", s.ListenAddrAndPort, symbol)
		fmt.Fprintf(w, "http://%s/price?symbols=%s</a></li>", s.ListenAddrAndPort, symbol)
	}
}
