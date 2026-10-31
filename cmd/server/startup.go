package main

import (
	"log"
	"net/http"
)

// startServer starts the HTTP server and reports unexpected listener errors.
func startServer(server *http.Server) <-chan error {
	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("Server listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	return serverErrors
}
