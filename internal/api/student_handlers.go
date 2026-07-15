package api

import (
	"encoding/json"
	"net/http"

	"campuscore/internal/auth"
	"campuscore/internal/middleware"
	"campuscore/internal/models"
)

// StudentHandler handles student requests.
type StudentHandler struct {
	academicService AcademicService
	ticketService   TicketService
}

// NewStudentHandler creates a StudentHandler.
func NewStudentHandler(
	academicService AcademicService,
	ticketService TicketService,
) *StudentHandler {
	return &StudentHandler{
		academicService: academicService,
		ticketService:   ticketService,
	}
}

// CourseRegistrationRequest represents a course registration request.
type CourseRegistrationRequest struct {
	CourseCode string `json:"course_code"`
	Session    string `json:"session"`
	Semester   string `json:"semester"`
}

// TicketSubmissionRequest represents a support ticket request.
type TicketSubmissionRequest struct {
	Category string `json:"category"`
	Subject  string `json:"subject"`
	Message  string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func getSession(r *http.Request) (*auth.Session, bool) {
	session, ok := r.Context().Value(middleware.UserContextKey).(*auth.Session)
	return session, ok
}

// RegisterCourse registers a student for a course.
func (h *StudentHandler) RegisterCourse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method not allowed",
		})
		return
	}

	session, ok := getSession(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
		return
	}

	var req CourseRegistrationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}

	if err := h.academicService.RegisterCourse(
		session.UserID,
		req.CourseCode,
		req.Session,
		req.Semester,
	); err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "course registered successfully",
	})
}

// SubmitTicket creates a support ticket.
func (h *StudentHandler) SubmitTicket(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method not allowed",
		})
		return
	}

	session, ok := getSession(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
		return
	}

	var req TicketSubmissionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}

	ticket := &models.SupportTicket{
		StudentID: session.UserID,
		Category:  req.Category,
		Subject:   req.Subject,
		Message:   req.Message,
		Status:    "open",
	}

	if err := h.ticketService.SubmitHelpdeskTicket(r.Context(), ticket); err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"message": "ticket submitted successfully",
	})
}