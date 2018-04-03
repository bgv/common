package middleware

import (
	"net/http"
)

// ServerName exports servername header to all requests
func ServerName(name string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Server", name)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
