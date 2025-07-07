package middleware

import (
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"runtime/debug"
)

func NewPanicMiddleware(logger *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rcv := recover(); rcv != nil {
					// Log the panic and stack trace
					logger.Printf("Panic recovered: %v\nStack trace:\n%s", rcv, debug.Stack())

					// Respond to the client with a 500 Internal Server Error
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// PanicRecoveryMiddleware recovers from panics and logs the stack trace.
func PanicRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rcv := recover(); rcv != nil {
				// Log the panic and stack trace
				log.Printf("Panic recovered: %v\nStack trace:\n%s", rcv, debug.Stack())

				// Respond to the client with a 500 Internal Server Error
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
