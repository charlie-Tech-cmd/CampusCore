package api

import (
	"encoding/json"
	"net/http"

	"campuscore/internal/models"
	"campuscore/internal/services"
)

type AdmissionHandler struct {
	service *services.AdmissionService
}

func NewAdmissionHandler(
	service *services.AdmissionService,
) *AdmissionHandler {

	return &AdmissionHandler{
		service: service,
	}
}

// SubmitApplication handles POST /admission/apply
func (h *AdmissionHandler) SubmitApplication(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var application models.AdmissionApplication

	if err := json.NewDecoder(r.Body).Decode(&application); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.SubmitApplication(&application); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// ListApplications handles GET /admission
func (h *AdmissionHandler) ListApplications(
	w http.ResponseWriter,
	r *http.Request,
) {

	applications, err := h.service.ListApplications()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(applications)
}

// GetApplication handles GET /admission/application
func (h *AdmissionHandler) GetApplication(
	w http.ResponseWriter,
	r *http.Request,
) {

	number := r.URL.Query().Get("application_no")

	application, err := h.service.GetApplication(number)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(application)
}

// ApproveApplication handles POST /admission/approve
func (h *AdmissionHandler) ApproveApplication(
	w http.ResponseWriter,
	r *http.Request,
) {

	number := r.FormValue("application_no")

	if err := h.service.ApproveApplication(number); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// RejectApplication handles POST /admission/reject
func (h *AdmissionHandler) RejectApplication(
	w http.ResponseWriter,
	r *http.Request,
) {

	number := r.FormValue("application_no")

	if err := h.service.RejectApplication(number); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
