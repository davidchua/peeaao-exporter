package main

import (
	"log"
	"net/http"
	"os"

	"github.com/davidchua/peeaao-exporter/internal/pinger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	addr := os.Getenv("ADDR")
	targets := os.Getenv("TARGETS")
	locations := os.Getenv("LOCATIONS")
	authToken := os.Getenv("AUTH_TOKEN")

	if targets == "" {
		log.Fatalf("environment variable TARGETS is not set")
	}

	if authToken == "" {
		log.Fatalf("environment variable AUTH_TOKEN is not set")
	}

	if addr == "" {
		addr = ":9100"
	}

	// prometheus exporter
	pingCollector := pinger.NewPingCollector(targets, locations, authToken)
	prometheus.MustRegister(pingCollector)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(addr, nil))
}
