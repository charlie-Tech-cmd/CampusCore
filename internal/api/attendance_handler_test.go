package api

import (
	"context"
	"net/http"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"campuscore/internal/auth"
	"campuscore/internal/middleware"
	"campuscore/internal/models"
	"campuscore/internal/services"
)

type attendanceRepositorySpy struct {
	requestedStudentID string
	called             bool
	recordByID         *models.Attendance
}

func (r *attendanceRepositorySpy) FindByID(int) (*models.Attendance, error) {
	return r.recordByID, nil
}

func (r *attendanceRepositorySpy) Create(*models.Attendance) error {
	return nil
}

func (r *attendanceRepositorySpy) GetAll() ([]models.Attendance, error) {
	return nil, nil
}

func (r *attendanceRepositorySpy) GetByStudent(studentID string) ([]models.Attendance, error) {
	r.called = true
	r.requestedStudentID = studentID
	return []models.Attendance{}, nil
}

func (r *attendanceRepositorySpy) GetByCourse(string) ([]models.Attendance, error) {
	return nil, nil
}

func (r *attendanceRepositorySpy) GetByLecturer(string) ([]models.Attendance, error) {
	return nil, nil
}

func (r *attendanceRepositorySpy) Update(*models.Attendance) error {
	return nil
}

func (r *attendanceRepositorySpy) Delete(int) error {
	return nil
}

func attendanceRequest(studentID string) *http.Request {
	url := "/api/v1/attendance/student"
	if studentID != "" {
		url += "?student_id=" + studentID
	}

	req := httptest.NewRequest(http.MethodGet, url, nil)
	ctx := context.WithValue(
		req.Context(),
		middleware.UserContextKey,
		&auth.Session{UserID: "STU001", Role: "student"},
	)

	return req.WithContext(ctx)
}

func TestListStudentAttendance_RejectsAnotherStudentsRecords(t *testing.T) {
	repo := &attendanceRepositorySpy{}
	handler := NewAttendanceHandler(services.NewAttendanceService(repo))
	rec := httptest.NewRecorder()

	handler.ListStudentAttendance(rec, attendanceRequest("STU002"))

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, rec.Code)
	}

	if repo.called {
		t.Fatal("repository must not be called for another student's records")
	}
}

func TestListStudentAttendance_AllowsOwnRecords(t *testing.T) {
	repo := &attendanceRepositorySpy{}
	handler := NewAttendanceHandler(services.NewAttendanceService(repo))
	rec := httptest.NewRecorder()

	handler.ListStudentAttendance(rec, attendanceRequest("STU001"))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if !repo.called || repo.requestedStudentID != "STU001" {
		t.Fatalf("expected repository query for STU001; called=%v, studentID=%q",
			repo.called, repo.requestedStudentID)
	}
}

func TestListStudentAttendance_DefaultsToAuthenticatedStudent(t *testing.T) {
	repo := &attendanceRepositorySpy{}
	handler := NewAttendanceHandler(services.NewAttendanceService(repo))
	rec := httptest.NewRecorder()

	handler.ListStudentAttendance(rec, attendanceRequest(""))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if !repo.called || repo.requestedStudentID != "STU001" {
		t.Fatalf("expected repository query for STU001 when student_id is omitted; called=%v, studentID=%q",
			repo.called, repo.requestedStudentID)
	}
}

func TestGetAttendance_RejectsAnotherStudentsRecord(t *testing.T) {
	repo := &attendanceRepositorySpy{
		recordByID: &models.Attendance{
			ID:        42,
			StudentID: "STU002",
		},
	}

	handler := NewAttendanceHandler(
		services.NewAttendanceService(repo),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/attendance/get?id=42",
		nil,
	)

	ctx := context.WithValue(
		req.Context(),
		middleware.UserContextKey,
		&auth.Session{
			UserID: "STU001",
			Role:   "student",
		},
	)
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()
	handler.GetAttendance(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf(
			"expected status %d, got %d; body: %s",
			http.StatusForbidden,
			rec.Code,
			rec.Body.String(),
		)
	}
}

func TestGetAttendance_AllowsOwnRecord(t *testing.T) {
	repo := &attendanceRepositorySpy{
		recordByID: &models.Attendance{
			ID:        42,
			StudentID: "STU001",
		},
	}

	handler := NewAttendanceHandler(
		services.NewAttendanceService(repo),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/attendance/get?id=42",
		nil,
	)

	ctx := context.WithValue(
		req.Context(),
		middleware.UserContextKey,
		&auth.Session{
			UserID: "STU001",
			Role:   "student",
		},
	)
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()
	handler.GetAttendance(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d; body: %s",
			http.StatusOK,
			rec.Code,
			rec.Body.String(),
		)
	}

	var record models.Attendance
	if err := json.NewDecoder(rec.Body).Decode(&record); err != nil {
		t.Fatalf("failed to decode attendance response: %v", err)
	}

	if record.ID != 42 || record.StudentID != "STU001" {
		t.Fatalf(
			"expected attendance record 42 for STU001, got ID=%d, StudentID=%q",
			record.ID,
			record.StudentID,
		)
	}
}
