package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	collector := NewGpuCollector()
	prometheus.MustRegister(collector)

	err := http.ListenAndServe(":9101", promhttp.Handler())
	if err != nil {
		log.Fatal(err)
	}
}
