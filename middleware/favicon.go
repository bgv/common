package middleware

import (
	"net/http"
)

// Favicon mock the presence of favicon.ico to browsers
// usage: router.Use(middleware.Favicon)
func Favicon(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path[1:] == "favicon.ico" {
			w.Header().Set("Content-Type", "image/x-icon")
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

