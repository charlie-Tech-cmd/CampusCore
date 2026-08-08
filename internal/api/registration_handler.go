package api

import (
	"encoding/json"
	"net/http"

	"campuscore/internal/middleware"
	"campuscore/internal/models"
)

type RegistrationService interface {
	RegisterCourse(
		studentID string,
		courseCode string,
		session string,
		semester string,
	) error

	GetStudentCourses(
		studentID string,
	) ([]models.Enrollment, error)

	DropCourse(
		studentID string,
		courseCode string,
	) error
}

type RegistrationHandler struct {
	service RegistrationService
}

func NewRegistrationHandler(service RegistrationService) *RegistrationHandler {
	return &RegistrationHandler{
		service: service,
	}
}

func (h *RegistrationHandler) DropCourse(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	studentID := middleware.CurrentUserID(r)

	if studentID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	courseCode := r.URL.Query().Get("course")

	if courseCode == "" {
		http.Error(w, "missing course code", http.StatusBadRequest)
		return
	}

	err := h.service.DropCourse(
		studentID,
		courseCode,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"message": "course dropped successfully",
	})
}

type RegisterCourseRequest struct {
	CourseCode string `json:"course_code"`
	Session    string `json:"session"`
	Semester   string `json:"semester"`
}

type RegisterCourseResponse struct {
	Message string `json:"message"`
}

func (h *RegistrationHandler) RegisterCourse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	studentID := middleware.CurrentUserID(r)

	if studentID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req RegisterCourseRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := h.service.RegisterCourse(
		studentID,
		req.CourseCode,
		req.Session,
		req.Semester,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(RegisterCourseResponse{
		Message: "course registered successfully",
	})
}

func (h *RegistrationHandler) GetStudentCourses(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	studentID := middleware.CurrentUserID(r)

	if studentID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	courses, err := h.service.GetStudentCourses(studentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(courses)
}
