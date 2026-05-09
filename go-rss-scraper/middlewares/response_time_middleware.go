package middlewares

import (
	"log"
	"net/http"
	"time"
)

// ResponseTimeMiddleware logs the time taken to process each requests
func ResponseTimeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Call the next handler
		next.ServeHTTP(w, r)

		duration := time.Since(start)
		log.Printf("[%s] %s - %v", r.Method, r.URL.Path, duration)
	})
}
