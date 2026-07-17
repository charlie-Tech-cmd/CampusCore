package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecovery(t *testing.T) {
	handler := Recovery(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected %d, got %d",
			http.StatusInternalServerError,
			rr.Code,
		)
	}

	body := rr.Body.String()

	if !strings.Contains(body, "error") {
		t.Fatalf("expected error response, got %s", body)
	}
}

func TestRecovery_NoPanic(t *testing.T) {
	handler := Recovery(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}

	if rr.Body.String() != "OK" {
		t.Fatalf("unexpected body %q", rr.Body.String())
	}
}
