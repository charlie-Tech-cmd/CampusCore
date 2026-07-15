package api

import (
	"encoding/json"
	"net/http"
)

// PaymentHandler handles payment requests.
type PaymentHandler struct {
    paymentService PaymentService
}

// NewPaymentHandler creates a PaymentHandler.
func NewPaymentHandler(paymentService PaymentService) *PaymentHandler {
    return &PaymentHandler{
        paymentService: paymentService,
    }
}

// PaymentRequest represents a payment request.
type PaymentRequest struct {
	StudentID  string  `json:"student_id"`
	Reference  string  `json:"reference"`
	Amount     float64 `json:"amount"`
	FeeType    string  `json:"fee_type"`
	Session    string  `json:"session"`
}

// VerifyPayment validates and records a payment.
func (h *PaymentHandler) VerifyPayment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req PaymentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	err := h.paymentService.ProcessPayment(
		r.Context(),
		req.StudentID,
		req.Reference,
		req.Amount,
		req.FeeType,
		req.Session,
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "payment recorded successfully",
	})
}