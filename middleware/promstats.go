package middleware

import (
	"net/http"
	"os"
	"time"

	mw "github.com/pressly/chi/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	defaultBuckets = []float64{300, 1200, 5000}
)

const (
	reqsName    = "requests_total"
	latencyName = "request_duration_milliseconds"
)

// PrometheusHandler bootstraps Prometheus for metrics collection
func PrometheusHandler() http.Handler {
	return promhttp.Handler()
}

// PrometheusStats returns a new Prometheus middleware handler.
func PrometheusStats(name string, buckets ...float64) func(http.Handler) http.Handler {

	requests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem:   "http",
			Name:        reqsName,
			Help:        "How many HTTP requests processed, partitioned by status code, method and HTTP path.",
			ConstLabels: prometheus.Labels{"service": name},
		},
		[]string{"code", "method"},
	)
	prometheus.MustRegister(requests)

	if len(buckets) == 0 {
		buckets = defaultBuckets
	}

	latency := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem:   "http",
		Name:        latencyName,
		Help:        "How long it took to process the request, partitioned by status code, method and HTTP path.",
		ConstLabels: prometheus.Labels{"service": name},
		Buckets:     buckets,
	},
		[]string{"code", "method"},
	)
	prometheus.MustRegister(latency)

	f := func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := mw.NewWrapResponseWriter(w, r.ProtoMajor)
			h.ServeHTTP(ww, r)
			requests.WithLabelValues(http.StatusText(ww.Status()), r.Method).Inc()
			latency.WithLabelValues(http.StatusText(ww.Status()), r.Method).Observe(float64(time.Since(start).Nanoseconds()) / 1e6) // milliseconds
		}
		return http.HandlerFunc(fn)
	}
	return f
}

// PrometheusStats returns a new Prometheus middleware handler.
func PrometheusDetailedStats(name string, buckets ...float64) func(http.Handler) http.Handler {

	requests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem:   "http",
			Name:        reqsName,
			Help:        "How many HTTP requests processed, partitioned by status code, method and HTTP path.",
			ConstLabels: prometheus.Labels{"service": name},
		},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(requests)

	if len(buckets) == 0 {
		buckets = defaultBuckets
	}

	latency := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem:   "http",
		Name:        latencyName,
		Help:        "How long it took to process the request, partitioned by status code, method and HTTP path.",
		ConstLabels: prometheus.Labels{"service": name},
		Buckets:     buckets,
	},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(latency)

	f := func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := mw.NewWrapResponseWriter(w, r.ProtoMajor)
			h.ServeHTTP(ww, r)
			requests.WithLabelValues(http.StatusText(ww.Status()), r.Method, r.URL.Path).Inc()
			latency.WithLabelValues(http.StatusText(ww.Status()), r.Method, r.URL.Path).Observe(float64(time.Since(start).Nanoseconds()) / 1e6) // milliseconds
		}
		return http.HandlerFunc(fn)
	}
	return f
}
