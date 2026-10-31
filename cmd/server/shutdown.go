package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const shutdownTimeout = 5 * time.Second

// waitForShutdown waits for an operating-system signal or an unexpected
// server error, then shuts down the HTTP server.
func waitForShutdown(server *http.Server, serverErrors <-chan error) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(stop)

	select {
	case receivedSignal := <-stop:
		log.Printf("Received %s; shutting down...", receivedSignal)
	case err := <-serverErrors:
		log.Printf("HTTP server stopped unexpectedly: %v", err)
	}

	if err := shutdownServer(server, shutdownTimeout); err != nil {
		log.Printf("HTTP server shutdown: %v", err)
	}

	log.Println("Server stopped.")
}

// shutdownServer attempts graceful shutdown and force-closes the server if
// the shutdown deadline expires or graceful shutdown otherwise fails.
func shutdownServer(server *http.Server, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		if closeErr := server.Close(); closeErr != nil && !errors.Is(closeErr, http.ErrServerClosed) {
			return fmt.Errorf("graceful shutdown failed: %w (force close failed: %v)", err, closeErr)
		}
		return err
	}

	return nil
}
