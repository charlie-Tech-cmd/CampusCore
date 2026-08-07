package api

import (
	"encoding/json"
	"net/http"

	"campuscore/internal/services"
)

type DashboardHandler struct {
	service *services.DashboardService
}

func NewDashboardHandler(
	service *services.DashboardService,
) *DashboardHandler {
	return &DashboardHandler{
		service: service,
	}
}

// GET /api/student/dashboard
func (h *DashboardHandler) GetStudentDashboard(
	w http.ResponseWriter,
	r *http.Request,
) {
	// Temporary: until JWT middleware injects the authenticated user.
	studentID := r.URL.Query().Get("student_id")

	if studentID == "" {
		http.Error(
			w,
			"student_id is required",
			http.StatusBadRequest,
		)
		return
	}

	dashboard, err := h.service.GetStudentDashboard(studentID)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(dashboard)
}
