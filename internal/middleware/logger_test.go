package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	var buf bytes.Buffer

	// Capture log output.
	originalOutput := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(originalOutput)

	handler := Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("created"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/students", nil)
	req.RemoteAddr = "127.0.0.1:12345"

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	logOutput := buf.String()

	expected := []string{
		"GET",
		"/students",
		"Status: 201",
		"127.0.0.1:12345",
	}

	for _, item := range expected {
		if !strings.Contains(logOutput, item) {
			t.Fatalf("expected log to contain %q\nActual log:\n%s", item, logOutput)
		}
	}
}
