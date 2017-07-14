package middleware

import (
	"net/http"
	"time"

	mw "github.com/go-chi/chi/middleware"
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

// PromMiddleware is a handler that exposes prometheus metrics for the number of requests,
// the latency and the response size, partitioned by status code, method and HTTP path.
type PromMiddleware struct {
	reqs    *prometheus.CounterVec
	latency *prometheus.HistogramVec
	service string
}

func Prometheus(service string, buckets ...float64) *PromMiddleware {
	prom := &PromMiddleware{
		service: service,
	}

	if len(buckets) == 0 {
		buckets = defaultBuckets
	}

	prom.registerMetrics(buckets...)

	return prom
}

// HandlerFunc returns Prometheus http.HandlerFunc for metrics collection
func (pmw *PromMiddleware) HandlerFunc() http.HandlerFunc {
	return promhttp.Handler().ServeHTTP
}

// Middleware returns a new Prometheus middleware handler.
func (pmw *PromMiddleware) Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := mw.NewWrapResponseWriter(w, r.ProtoMajor)
		h.ServeHTTP(ww, r)
		pmw.reqs.WithLabelValues(http.StatusText(ww.Status()), r.Method).Inc()
		pmw.latency.WithLabelValues(http.StatusText(ww.Status()), r.Method).Observe(float64(time.Since(start).Nanoseconds()) / 1e6) // milliseconds
	})
}

func (pmw *PromMiddleware) registerMetrics(buckets ...float64) {
	pmw.reqs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem:   "http",
			Name:        reqsName,
			Help:        "How many HTTP requests processed, partitioned by status code, method and HTTP path.",
			ConstLabels: prometheus.Labels{"service": pmw.service},
		},
		[]string{"code", "method"},
	)

	pmw.latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem:   "http",
		Name:        latencyName,
		Help:        "How long it took to process the request, partitioned by status code, method and HTTP path.",
		ConstLabels: prometheus.Labels{"service": pmw.service},
		Buckets:     buckets,
	},
		[]string{"code", "method"},
	)

	prometheus.MustRegister(pmw.reqs, pmw.latency)
}
