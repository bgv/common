package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pressly/chi"
)

func TestPrometheusMiddleware(t *testing.T) {
	recorder := httptest.NewRecorder()

    // Configure Prometheus middleware with service name test and latency buckets
	metrics := Prometheus("test", 300, 500, 1000, 2000)

    r := chi.NewRouter()

	r.Use(metrics.Middleware)

	r.Get("/metrics", metrics.HandlerFunc())
	r.Get("/ok", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

    // Make a test request to record some metrics before we test the /metrics enpoint
	req1, err := http.NewRequest("GET", "/ok", nil)
	if err != nil {
		t.Error(err)
	}
	req2, err := http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		t.Error(err)
	}

    // Record and test the request results
	r.ServeHTTP(recorder, req1)
	r.ServeHTTP(recorder, req2)
	body := recorder.Body.String()

	if !strings.Contains(body, reqsName) {
		t.Errorf("body does not contain request total entry '%s'", reqsName)
	}
	if !strings.Contains(body, latencyName) {
		t.Errorf("body does not contain request duration entry '%s'", latencyName)
	}
}
