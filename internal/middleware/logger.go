package middleware

import (
	"log"
	"net/http"
	"time"
)

// responseWriterInterceptor wraps standard http.ResponseWriter to capture the HTTP status code
type responseWriterInterceptor struct {
	http.ResponseWriter
	statusCode int
}

func (rwi *responseWriterInterceptor) WriteHeader(code int) {
	rwi.statusCode = code
	rwi.ResponseWriter.WriteHeader(code)
}

// Logger is a structural middleware that logs request execution times and statuses
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Intercept the response headers to record the outgoing status code
		interceptor := &responseWriterInterceptor{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // Default status code if WriteHeader is not explicitly called
		}

		// Pass the request along the pipeline execution chain to the next handler
		next.ServeHTTP(interceptor, r)

		// Calculate total time taken to process the query execution pipeline
		duration := time.Since(startTime)

		log.Printf("📥 [%s] %s %s | Status: %d | Execution Time: %v",
			r.RemoteAddr,
			r.Method,
			r.URL.Path,
			interceptor.statusCode,
			duration,
		)
	})
}