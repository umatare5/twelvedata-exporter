// Package internal is a server that uses the twelvedata API as its backend.
package internal

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Server struct
type Server struct {
	Addr      string
	Port      int
	APIKey    string
	RateLimit int
}

// New returns Twelvedata struct
func New(addr string, port int, apikey string, limit int) Server {
	return Server{
		Addr:      addr,
		Port:      port,
		APIKey:    apikey,
		RateLimit: limit,
	}
}

// Boot the server
func (s *Server) Boot() {
	reg := prometheus.NewRegistry()

	// Add standard process and Go metrics.
	reg.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(),
	)

	// Add handlers.
	http.HandleFunc("/", s.help)
	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/price", func(w http.ResponseWriter, r *http.Request) {
		s.priceHandler(w, r)
	})

	log.Print("Listening on port ", s.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.Port), nil))
}

// help returns a help message for those using the root URL.
func (s *Server) help(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "<h1>Prometheus Quotes Exporter</h1>")
	fmt.Fprintf(w, "<p>To fetch quotes, your URL must be formatted as:</p>")
	fmt.Fprintf(w, "http://localhost:%d/price?symbols=AAAA,BBBB,CCCC", s.Port)
	fmt.Fprintf(w, "<p><b>Examples:</b></p>")
	fmt.Fprintf(w, "<ul>")

	symbols := []string{
		"AMD",
		"AMZN,GOOG",
	}

	for _, symbol := range symbols {
		fmt.Fprintf(w, "<li><a href=\"http://localhost:%d/price?symbols=%s\">", s.Port, symbol)
		fmt.Fprintf(w, "http://localhost:%d/price?symbols=%s</a></li>", s.Port, symbol)
	}
}

// priceHandler handles the "/price" endpoint. It creates a new collector with
// the URL and a new prometheus registry to use that collector.
func (s *Server) priceHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("URL: %s\n", r.RequestURI)

	collector, err := newCollector(r.URL, s.APIKey, s.RateLimit)
	if err != nil {
		log.Print(err)
		return
	}

	registry := prometheus.NewRegistry()

	// These will be collected every time the /stock or /fund endpoint is reached.
	registry.MustRegister(
		&collector,
		queryCount,
		queryDuration,
		errorCount)

	// Delegate http serving to Promethues client library, which will call collector.Collect.
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}
