package api

import (
	"encoding/json"
	"net/http"

	"campuscore/internal/services"
)

// ReportingHandler handles reporting endpoints.
type ReportingHandler struct {
	service *services.ReportingService
}

// NewReportingHandler creates a reporting handler.
func NewReportingHandler(
	service *services.ReportingService,
) *ReportingHandler {

	return &ReportingHandler{
		service: service,
	}
}

// GetEnrollmentSummary handles GET /reports/enrollment.
func (h *ReportingHandler) GetEnrollmentSummary(
	w http.ResponseWriter,
	r *http.Request,
) {

	summary, err := h.service.GetEnrollmentSummary()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

// GetPaymentSummary handles GET /reports/payments.
func (h *ReportingHandler) GetPaymentSummary(
	w http.ResponseWriter,
	r *http.Request,
) {

	summary, err := h.service.GetPaymentSummary()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

// GetAcademicPerformanceSummary handles GET /reports/academic.
func (h *ReportingHandler) GetAcademicPerformanceSummary(
	w http.ResponseWriter,
	r *http.Request,
) {

	summary, err := h.service.GetAcademicPerformanceSummary()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

// GetClearanceSummary handles GET /reports/clearance.
func (h *ReportingHandler) GetClearanceSummary(
	w http.ResponseWriter,
	r *http.Request,
) {

	summary, err := h.service.GetClearanceSummary()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}
