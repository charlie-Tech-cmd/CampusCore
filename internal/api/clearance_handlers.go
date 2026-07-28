package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"campuscore/internal/models"
)

type ClearanceService interface {
	GetStudentClearance(
		ctx context.Context,
		studentID string,
	) ([]models.StudentClearance, error)

	UpdateClearance(
		ctx context.Context,
		studentID string,
		officeID int,
		status models.ClearanceStatus,
		staffID string,
	) error

	IsStudentCleared(
		ctx context.Context,
		studentID string,
	) (bool, error)
}

type ClearanceHandler struct {
	service ClearanceService
}

func NewClearanceHandler(service ClearanceService) *ClearanceHandler {
	return &ClearanceHandler{
		service: service,
	}
}

func (h *ClearanceHandler) GetStudentClearance(
	w http.ResponseWriter,
	r *http.Request,
) {

	studentID := r.URL.Query().Get("student_id")
	if studentID == "" {
		http.Error(w, "student_id is required", http.StatusBadRequest)
		return
	}

	clearances, err := h.service.GetStudentClearance(
		context.Background(),
		studentID,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clearances)
}

func (h *ClearanceHandler) IsStudentCleared(
	w http.ResponseWriter,
	r *http.Request,
) {

	studentID := r.URL.Query().Get("student_id")
	if studentID == "" {
		http.Error(w, "student_id is required", http.StatusBadRequest)
		return
	}

	cleared, err := h.service.IsStudentCleared(
		context.Background(),
		studentID,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{
		"cleared": cleared,
	})
}

func (h *ClearanceHandler) UpdateClearance(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	officeID, err := strconv.Atoi(r.FormValue("office_id"))
	if err != nil {
		http.Error(w, "invalid office_id", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateClearance(
		context.Background(),
		r.FormValue("student_id"),
		officeID,
		models.ClearanceStatus(r.FormValue("status")),
		r.FormValue("staff_id"),
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
