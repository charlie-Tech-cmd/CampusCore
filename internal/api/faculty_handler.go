package api

import (
	"encoding/json"
	"net/http"

	"campuscore/internal/models"
)

// FacultyService defines the business operations required by the handler.
type FacultyService interface {
	CreateFaculty(*models.Faculty) error
	GetFaculty(int) (*models.Faculty, error)
	GetFacultyByCode(string) (*models.Faculty, error)
	ListFaculties() ([]models.Faculty, error)
	UpdateFaculty(*models.Faculty) error
	DeleteFaculty(int) error
}

// FacultyHandler handles faculty endpoints.
type FacultyHandler struct {
	service FacultyService
}

// NewFacultyHandler creates a FacultyHandler.
func NewFacultyHandler(service FacultyService) *FacultyHandler {
	return &FacultyHandler{
		service: service,
	}
}

// List handles GET /faculties.
func (h *FacultyHandler) List(w http.ResponseWriter, r *http.Request) {
	faculties, err := h.service.ListFaculties()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(faculties)
}
