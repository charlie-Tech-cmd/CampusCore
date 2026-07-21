package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"campuscore/internal/models"
)

// CourseService defines the service contract used by the handler.
type CourseService interface {
	CreateCourse(*models.Course) error
	GetCourse(string) (*models.Course, error)
	ListCourses() ([]models.Course, error)
	ListDepartmentCourses(int) ([]models.Course, error)
	UpdateCourse(*models.Course) error
	DeleteCourse(string) error
}

// CourseHandler handles course endpoints.
type CourseHandler struct {
	service CourseService
}

// NewCourseHandler creates a new course handler.
func NewCourseHandler(service CourseService) *CourseHandler {
	return &CourseHandler{
		service: service,
	}
}

// Create creates a new course.
func (h *CourseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var course models.Course

	if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreateCourse(&course); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(course)
}

// Get returns a single course.
func (h *CourseHandler) Get(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	course, err := h.service.GetCourse(code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(course)
}

// List returns all courses.
func (h *CourseHandler) List(w http.ResponseWriter, r *http.Request) {
	courses, err := h.service.ListCourses()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(courses)
}

// ListByDepartment returns all courses for a department.
func (h *CourseHandler) ListByDepartment(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("department_id")

	departmentID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid department id", http.StatusBadRequest)
		return
	}

	courses, err := h.service.ListDepartmentCourses(departmentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(courses)
}

// Update updates a course.
func (h *CourseHandler) Update(w http.ResponseWriter, r *http.Request) {
	var course models.Course

	if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateCourse(&course); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(course)
}

// Delete deletes a course.
func (h *CourseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	if err := h.service.DeleteCourse(code); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
