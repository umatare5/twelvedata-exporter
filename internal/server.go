// Package internal provides the HTTP server implementation for the exporter.
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

// Server represents the HTTP server for the exporter.
type Server struct {
	Client *TwelvedataClient // Twelvedata API client
	Config *config.Config    // Configuration for the server
}

// NewServer initializes and returns a new Server instance.
func NewServer(config *config.Config) (Server, error) {
	return Server{
		Client: NewTwelvedataClient(config.TwelvedataAPIKey),
		Config: config,
	}, nil
}

// Start configures and launches the HTTP server to serve metrics and help pages.
func (s *Server) Start() {
	reg := prometheus.NewRegistry()

	// Register standard process and Go metrics.
	reg.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(),
	)

	// Register HTTP handlers.
	http.HandleFunc("/", s.help)
	http.HandleFunc(s.Config.WebScrapePath, func(w http.ResponseWriter, r *http.Request) {
		s.priceHandler(w, r)
	})

	listenAddr := s.Config.WebListenAddress + ":" + strconv.Itoa(s.Config.WebListenPort)
	log.Infof("Starting the Twelvedata exporter on %s", listenAddr)
	srv := &http.Server{
		Addr:         listenAddr,
		Handler:      nil,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}
	
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}

// priceHandler registers Prometheus metrics and serves them via HTTP.
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

	// Register the Twelvedata collector and metrics.
	registry.MustRegister(
		newCollector(s.Client, symbols),
		queryCount,
		queryDuration,
		errorCount,
	)

	// Serve metrics using Prometheus client library.
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		ErrorHandling: promhttp.ContinueOnError,
	})
	h.ServeHTTP(w, r)
}

// help generates and serves an HTML help page for the root URL.
func (s *Server) help(w http.ResponseWriter, _ *http.Request) {
	listenAddrAndPort := s.Config.WebListenAddress + ":" + strconv.Itoa(s.Config.WebListenPort)

	var builder strings.Builder
	builder.WriteString("<h1>Prometheus Twelvedata Exporter</h1>")
	builder.WriteString("<p>To fetch the price of quotes, your URL must be formatted as:</p>")
	builder.WriteString(fmt.Sprintf("http://%s%s?symbols=AAAA,BBBB,CCCC", listenAddrAndPort, s.Config.WebScrapePath))
	builder.WriteString("<p><b>Examples:</b></p>")
	builder.WriteString("<ul>")

	symbols := []string{
		"GOOGL",
		"AMZN,AAPL,MSFT",
	}

	for _, symbol := range symbols {
		builder.WriteString(fmt.Sprintf("<li><a href=\"http://%s%s?symbols=%s\">", listenAddrAndPort, s.Config.WebScrapePath, symbol))
		builder.WriteString(fmt.Sprintf("http://%s%s?symbols=%s</a></li>", listenAddrAndPort, s.Config.WebScrapePath, symbol))
	}
	builder.WriteString("</ul>")

	if _, err := w.Write([]byte(builder.String())); err != nil {
		log.Errorf("Error writing response: %v", err)
	}
}
