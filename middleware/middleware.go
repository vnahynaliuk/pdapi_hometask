package middleware

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

// Define counters and histograms for request metrics.
var (
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "endpoint"},
	)
	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
)

func init() {
	// Register metrics at package initialization.
	prometheus.MustRegister(HttpRequestsTotal, HttpRequestDuration)
}

// LoggingMiddleware logs the HTTP method, URI, and remote address for each request.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

// MetricsMiddleware collects metrics for each HTTP request.
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timer := prometheus.NewTimer(HttpRequestDuration.WithLabelValues(r.Method, r.RequestURI))
		defer timer.ObserveDuration()
		HttpRequestsTotal.WithLabelValues(r.Method, r.RequestURI).Inc()
		next.ServeHTTP(w, r)
	})
}
