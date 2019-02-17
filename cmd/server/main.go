package main

import (
	"fmt"
	"log"

	"github.com/gerlacdt/prom_example/pkg/http"
	"github.com/gerlacdt/prom_example/pkg/inmemory"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	inFlightGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "in_flight_requests",
		Help: "A gauge of requests currently being served by the wrapped handler.",
	})

	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_requests_total",
			Help: "A counter for requests to the wrapped handler.",
		},
		[]string{"code", "method"},
	)

	// duration is partitioned by the HTTP method and handler. It uses custom
	// buckets based on the expected request duration.
	duration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_duration_seconds",
			Help:    "A histogram of latencies for requests.",
			Buckets: []float64{.1, .3, .8, 2},
		},
		[]string{"handler", "method", "code"},
	)

	// responseSize has no labels, making it a zero-dimensional
	// ObserverVec.
	responseSize := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "response_size_bytes",
			Help:    "A histogram of response sizes for requests.",
			Buckets: []float64{200, 500, 900, 1500},
		},
		[]string{"code", "method"},
	)

	// Register all of the metrics in the standard registry.
	prometheus.MustRegister(inFlightGauge, counter, duration, responseSize)

	// var (
	// 	requestCounter = promauto.NewCounter(prometheus.CounterOpts{
	// 		Name: "http_processed_requests_total",
	// 		Help: "The total number of processed http requests",
	// 		// []string{"code", "method"},
	// 	})
	// )

	postService := inmemory.New()
	postHandler := http.New(postService)
	promPostHandler := wrapHandler(postHandler, "/v1/posts",
		&promMiddleware{inFlightGauge: inFlightGauge, counter: counter, duration: duration, responseSize: responseSize})

	healthHandler := wrapHandler(&http.HealthHandler{}, "/heath",
		&promMiddleware{inFlightGauge: inFlightGauge, counter: counter, duration: duration, responseSize: responseSize})

	// middleware := http.NewMiddleware(requestCounter, h)
	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/health", healthHandler)
	http.Handle("/v1/posts", promPostHandler)
	http.Handle("/v1/posts/", promPostHandler)

	// start server
	port := ":8080"
	fmt.Printf("Started server on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

type promMiddleware struct {
	inFlightGauge prometheus.Gauge
	counter       *prometheus.CounterVec
	duration      *prometheus.HistogramVec
	responseSize  *prometheus.HistogramVec
}

func wrapHandler(h http.Handler, path string, proms *promMiddleware) http.Handler {
	return promhttp.InstrumentHandlerInFlight(proms.inFlightGauge,
		promhttp.InstrumentHandlerDuration(proms.duration.MustCurryWith(prometheus.Labels{"handler": path}),
			promhttp.InstrumentHandlerCounter(proms.counter,
				promhttp.InstrumentHandlerResponseSize(proms.responseSize, h),
			),
		),
	)
}
