package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"campuscore/internal/models"
)

type ResultService interface {
	SubmitResult(*models.Result) error
	GetStudentResults(studentID string) ([]models.Result, error)
	GetCourseResults(courseCode string) ([]models.Result, error)
	UpdateResult(*models.Result) error
	DeleteResult(id int) error
}

type ResultHandler struct {
	service ResultService
}

func NewResultHandler(service ResultService) *ResultHandler {
	return &ResultHandler{
		service: service,
	}
}

func (h *ResultHandler) Submit(w http.ResponseWriter, r *http.Request) {
	var result models.Result

	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.SubmitResult(&result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *ResultHandler) StudentResults(w http.ResponseWriter, r *http.Request) {
	studentID := r.URL.Query().Get("student_id")

	results, err := h.service.GetStudentResults(studentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(results)
}

func (h *ResultHandler) CourseResults(w http.ResponseWriter, r *http.Request) {
	courseCode := r.URL.Query().Get("course_code")

	results, err := h.service.GetCourseResults(courseCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(results)
}

func (h *ResultHandler) Update(w http.ResponseWriter, r *http.Request) {
	var result models.Result

	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateResult(&result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ResultHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteResult(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
