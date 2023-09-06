package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	prometheus.MustRegister(testCounter)
	prometheus.MustRegister(testLantency)
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/test", test)
	mux.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9999", mux))

}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Home Page")
}

func test(w http.ResponseWriter, r *http.Request) {
	testCounter.Inc()
	var Status string
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		testLantency.WithLabelValues(Status).Observe(v)
	}))
	defer func() {
		timer.ObserveDuration()
	}()
	fmt.Fprint(w, "Segunda chamada")
}

var testCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "number_requests_test",
		Help: "no",
	},
)

var testLantency = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_request_get_duration_seconds",
		Help:    "Latency of test page",
		Buckets: prometheus.LinearBuckets(0.01, 0.05, 10),
	},
	[]string{"Status"},
)
