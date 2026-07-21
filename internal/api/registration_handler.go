package api

import (
	"encoding/json"
	"net/http"
)

type RegistrationService interface {
	RegisterCourse(
		studentID string,
		courseCode string,
		session string,
		semester string,
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

type RegisterCourseRequest struct {
	StudentID  string `json:"student_id"`
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

	var req RegisterCourseRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := h.service.RegisterCourse(
		req.StudentID,
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
