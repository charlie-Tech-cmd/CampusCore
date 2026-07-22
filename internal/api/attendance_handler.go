package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"campuscore/internal/models"
	"campuscore/internal/services"
)

type AttendanceHandler struct {
	service *services.AttendanceService
}

func NewAttendanceHandler(
	service *services.AttendanceService,
) *AttendanceHandler {
	return &AttendanceHandler{
		service: service,
	}
}

// MarkAttendance handles attendance creation.
func (h *AttendanceHandler) MarkAttendance(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var record models.Attendance

	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.MarkAttendance(&record); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(record); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetAttendance returns a single attendance record.
func (h *AttendanceHandler) GetAttendance(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid attendance id", http.StatusBadRequest)
		return
	}

	record, err := h.service.GetAttendance(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(record); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ListAttendance returns all attendance records.
func (h *AttendanceHandler) ListAttendance(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	records, err := h.service.ListAttendance()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(records); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ListStudentAttendance returns attendance for a student.
func (h *AttendanceHandler) ListStudentAttendance(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	studentID := r.URL.Query().Get("student_id")

	records, err := h.service.ListStudentAttendance(studentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(records); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ListCourseAttendance returns attendance for a course.
func (h *AttendanceHandler) ListCourseAttendance(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	courseCode := r.URL.Query().Get("course_code")

	records, err := h.service.ListCourseAttendance(courseCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(records); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ListLecturerAttendance returns attendance recorded by a lecturer.
func (h *AttendanceHandler) ListLecturerAttendance(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lecturerID := r.URL.Query().Get("lecturer_id")

	records, err := h.service.ListLecturerAttendance(lecturerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(records); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// UpdateAttendance updates an attendance record.
func (h *AttendanceHandler) UpdateAttendance(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var record models.Attendance

	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateAttendance(&record); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(record); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// DeleteAttendance removes an attendance record.
func (h *AttendanceHandler) DeleteAttendance(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid attendance id", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteAttendance(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
