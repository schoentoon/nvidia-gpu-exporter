package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:9101", "Listening address")
	flag.Parse()

	collector := NewGpuCollector()
	prometheus.MustRegister(collector)

	err := http.ListenAndServe(*addr, promhttp.Handler())
	if err != nil {
		log.Fatal(err)
	}
}
