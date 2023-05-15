package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/jessicalins/instrumentation-practices-examples/middleware/httpmiddleware"
)

func main() {
	// Create a new random number generator
	rand.Seed(time.Now().UnixNano())

	// Create non-global registry.
	registry := prometheus.NewRegistry()

	// Add go runtime metrics and process collectors.
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	// Expose /metrics HTTP endpoint using the created custom registry.
	http.Handle(
		"/metrics",
		httpmiddleware.New(
			registry, nil).
			WrapHandler("/metrics", promhttp.HandlerFor(
				registry,
				promhttp.HandlerOpts{}),
			))

    // Expose /healthz endpoint used by k8s probes
	http.Handle(
		"/healthz",
		httpmiddleware.New(
			registry, nil).
			WrapHandler("/healthz",
				func() http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						// Simulate DB ping for example up to 300ms
						time.Sleep(time.Millisecond * time.Duration(rand.Intn(300)))
						w.WriteHeader(http.StatusOK)
						w.Write([]byte("OK"))
					})
				}()))

    // Expose main catch-all endpoint that fails randomly with 500
	http.Handle(
		"/",
		httpmiddleware.New(
			registry, nil).
			WrapHandler("/",
				func() http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						// Simulate some heavier work up to 1s
						time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))

						// Fail randomly to make it interesting
						if rand.Intn(30) == 3 {
							w.WriteHeader(http.StatusInternalServerError)
							return
						}

						w.WriteHeader(http.StatusOK)
						w.Write([]byte("Hello, world!"))
					})
				}()))

	log.Fatalln(http.ListenAndServe(":8080", nil))
}
