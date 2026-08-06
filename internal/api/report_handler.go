package api

import (
	"encoding/json"
	"net/http"

	"campuscore/internal/services"
)

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(
	service *services.ReportService,
) *ReportHandler {

	return &ReportHandler{
		service: service,
	}
}

func (h *ReportHandler) DashboardSummary(
	w http.ResponseWriter,
	r *http.Request,
) {

	summary, err := h.service.DashboardSummary()
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(summary)
}

func (h *ReportHandler) AdmissionReport(
	w http.ResponseWriter,
	r *http.Request,
) {

	report, err := h.service.AdmissionReport()
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(report)
}
