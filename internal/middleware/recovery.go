package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
)

// Recovery intercepts runtime panics globally to prevent application server crashes
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Defer executes this block at the absolute end of the request lifecycle function wrapper
		defer func() {
			if err := recover(); err != nil {
				// 1. Log the panic details along with the full cryptographic execution stack trace
				log.Printf("🚨 [CRITICAL PANIC RECOVERY] Internal Server Anomaly Encountered: %v", err)
				log.Printf("📚 [STACK TRACE]:\n%s", debug.Stack())

				// 2. Safely respond to the client with an Internal Server Error code
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error": "A critical system anomaly occurred. The operations team has been notified."}`))
			}
		}()

		// Pass the request forward down the middleware execution pipeline
		next.ServeHTTP(w, r)
	})
}
