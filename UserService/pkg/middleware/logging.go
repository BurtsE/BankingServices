package middleware

import (
	"log"
	"net/http"
)

// LoggerMiddleware logs incoming requests.
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request received: %s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
