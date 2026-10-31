package main

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestStartServerReportsUnexpectedListenError(t *testing.T) {
	server := &http.Server{Addr: "127.0.0.1:-1"}
	serverErrors := startServer(server)

	select {
	case err := <-serverErrors:
		if err == nil {
			t.Fatal("expected a listener error, got nil")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for the listener error")
	}
}

func TestShutdownServerStopsServing(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("create listener: %v", err)
	}

	server := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}),
	}
	serveErrors := make(chan error, 1)

	go func() {
		serveErrors <- server.Serve(listener)
	}()

	response, err := http.Get("http://" + listener.Addr().String())
	if err != nil {
		t.Fatalf("make request to running server: %v", err)
	}
	_, _ = io.Copy(io.Discard, response.Body)
	_ = response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		t.Fatalf("unexpected response status: got %d, want %d", response.StatusCode, http.StatusNoContent)
	}

	if err := shutdownServer(server, time.Second); err != nil {
		t.Fatalf("shutdown server: %v", err)
	}

	select {
	case err := <-serveErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			t.Fatalf("unexpected Serve result: got %v, want %v", err, http.ErrServerClosed)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("server did not stop after shutdown")
	}
}

// TestShutdownServerForcesCloseAfterDeadline verifies that an active
// request cannot prevent the server listener from closing indefinitely.
func TestShutdownServerForcesCloseAfterDeadline(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("create listener: %v", err)
	}

	handlerStarted := make(chan struct{})
	releaseHandler := make(chan struct{})
	defer close(releaseHandler)

	server := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			close(handlerStarted)
			<-releaseHandler
			w.WriteHeader(http.StatusNoContent)
		}),
	}

	serveErrors := make(chan error, 1)
	go func() {
		serveErrors <- server.Serve(listener)
	}()

	requestDone := make(chan struct{})
	go func() {
		response, _ := http.Get("http://" + listener.Addr().String())
		if response != nil {
			_ = response.Body.Close()
		}
		close(requestDone)
	}()

	select {
	case <-handlerStarted:
	case <-time.After(2 * time.Second):
		t.Fatal("request handler did not start")
	}

	err = shutdownServer(server, 100*time.Millisecond)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected shutdown deadline error, got %v", err)
	}

	select {
	case serveErr := <-serveErrors:
		if !errors.Is(serveErr, http.ErrServerClosed) {
			t.Fatalf("unexpected Serve result: got %v, want %v", serveErr, http.ErrServerClosed)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("server listener did not close after forced shutdown")
	}

}
