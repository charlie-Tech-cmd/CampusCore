package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	// "campuscore/internal/models"

)

func TestNewPaymentHandler(t *testing.T) {
	service := &mockPaymentService{}

	handler := NewPaymentHandler(service)

	if handler == nil {
		t.Fatal("expected handler")
	}

	if handler.paymentService != service {
		t.Fatal("payment service not assigned")
	}
}

func TestVerifyPayment_MethodNotAllowed(t *testing.T) {
	handler := NewPaymentHandler(&mockPaymentService{})

	req := httptest.NewRequest(http.MethodGet, "/payment", nil)
	rec := httptest.NewRecorder()

	handler.VerifyPayment(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected %d, got %d",
			http.StatusMethodNotAllowed,
			rec.Code,
		)
	}
}

func TestVerifyPayment_InvalidJSON(t *testing.T) {
	handler := NewPaymentHandler(&mockPaymentService{})

	req := httptest.NewRequest(
		http.MethodPost,
		"/payment",
		strings.NewReader("{invalid"),
	)

	rec := httptest.NewRecorder()

	handler.VerifyPayment(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}
}

func TestVerifyPayment_ServiceError(t *testing.T) {
	service := &mockPaymentService{
		processPaymentFunc: func(
			ctx context.Context,
			studentID,
			reference string,
			amount float64,
			feeType,
			session string,
		) error {
			return errors.New("payment failed")
		},
	}

	handler := NewPaymentHandler(service)

	body := `{
		"student_id":"STU001",
		"reference":"REF123",
		"amount":5000,
		"fee_type":"school",
		"session":"2026/2027"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/payment",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	handler.VerifyPayment(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}
}

func TestVerifyPayment_Success(t *testing.T) {
	service := &mockPaymentService{
		processPaymentFunc: func(
			ctx context.Context,
			studentID,
			reference string,
			amount float64,
			feeType,
			session string,
		) error {
			return nil
		},
	}

	handler := NewPaymentHandler(service)

	body := `{
		"student_id":"STU001",
		"reference":"REF123",
		"amount":5000,
		"fee_type":"school",
		"session":"2026/2027"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/payment",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	handler.VerifyPayment(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected %d, got %d",
			http.StatusCreated,
			rec.Code,
		)
	}

	var response map[string]string

	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if response["message"] != "payment recorded successfully" {
		t.Fatal("unexpected response")
	}
}