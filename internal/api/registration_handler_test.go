package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"campuscore/internal/auth"
	"campuscore/internal/middleware"
	"campuscore/internal/models"
)

type mockRegistrationService struct {
	registerCourseFunc    func(studentID, courseCode, session, semester string) error
	getStudentCoursesFunc func(studentID string) ([]models.Enrollment, error)
	dropCourseFunc        func(studentID, courseCode, session, semester string) error
}

func (m *mockRegistrationService) RegisterCourse(
	studentID,
	courseCode,
	session,
	semester string,
) error {
	if m.registerCourseFunc != nil {
		return m.registerCourseFunc(studentID, courseCode, session, semester)
	}
	return nil
}

func (m *mockRegistrationService) GetStudentCourses(
	studentID string,
) ([]models.Enrollment, error) {
	if m.getStudentCoursesFunc != nil {
		return m.getStudentCoursesFunc(studentID)
	}
	return nil, nil
}

func (m *mockRegistrationService) DropCourse(
	studentID,
	courseCode,
	session,
	semester string,
) error {
	if m.dropCourseFunc != nil {
		return m.dropCourseFunc(studentID, courseCode, session, semester)
	}
	return nil
}

func authenticatedRequest(
	method,
	target string,
	body *bytes.Buffer,
) *http.Request {
	if body == nil {
		body = bytes.NewBuffer(nil)
	}

	req := httptest.NewRequest(method, target, body)

	session := &auth.Session{
		UserID: "STU001",
		Role:   "student",
	}

	ctx := context.WithValue(
		req.Context(),
		middleware.UserContextKey,
		session,
	)

	return req.WithContext(ctx)
}

func TestNewRegistrationHandler(t *testing.T) {
	service := &mockRegistrationService{}

	handler := NewRegistrationHandler(service)

	if handler == nil {
		t.Fatal("expected handler")
	}

	if handler.service != service {
		t.Fatal("expected registration service to be assigned")
	}
}

func TestRegistrationHandler_RegisterCourse_MethodNotAllowed(t *testing.T) {
	handler := NewRegistrationHandler(&mockRegistrationService{})

	req := httptest.NewRequest(
		http.MethodGet,
		"/student/courses",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.RegisterCourse(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf(
			"expected %d, got %d",
			http.StatusMethodNotAllowed,
			rec.Code,
		)
	}
}

func TestRegistrationHandler_RegisterCourse_Unauthorized(t *testing.T) {
	handler := NewRegistrationHandler(&mockRegistrationService{})

	req := httptest.NewRequest(
		http.MethodPost,
		"/student/courses",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.RegisterCourse(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected %d, got %d",
			http.StatusUnauthorized,
			rec.Code,
		)
	}
}

func TestRegistrationHandler_RegisterCourse_InvalidJSON(t *testing.T) {
	handler := NewRegistrationHandler(&mockRegistrationService{})

	req := authenticatedRequest(
		http.MethodPost,
		"/student/courses",
		bytes.NewBufferString("{"),
	)

	rec := httptest.NewRecorder()

	handler.RegisterCourse(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}
}

func TestRegistrationHandler_RegisterCourse_ServiceError(t *testing.T) {
	service := &mockRegistrationService{
		registerCourseFunc: func(
			studentID,
			courseCode,
			session,
			semester string,
		) error {
			if studentID != "STU001" {
				t.Fatalf("expected student STU001, got %q", studentID)
			}

			return errors.New("registration failed")
		},
	}

	handler := NewRegistrationHandler(service)

	body := bytes.NewBufferString(`{
		"course_code":"CSC401",
		"session":"2025/2026",
		"semester":"First"
	}`)

	req := authenticatedRequest(
		http.MethodPost,
		"/student/courses",
		body,
	)

	rec := httptest.NewRecorder()

	handler.RegisterCourse(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}

	if rec.Body.String() != "registration failed\n" {
		t.Fatalf(
			"expected registration error, got %q",
			rec.Body.String(),
		)
	}
}

func TestRegistrationHandler_RegisterCourse_Success(t *testing.T) {
	service := &mockRegistrationService{
		registerCourseFunc: func(
			studentID,
			courseCode,
			session,
			semester string,
		) error {
			if studentID != "STU001" {
				t.Fatalf("expected student STU001, got %q", studentID)
			}

			if courseCode != "CSC401" {
				t.Fatalf("expected CSC401, got %q", courseCode)
			}

			if session != "2025/2026" {
				t.Fatalf("expected 2025/2026, got %q", session)
			}

			if semester != "First" {
				t.Fatalf("expected First, got %q", semester)
			}

			return nil
		},
	}

	handler := NewRegistrationHandler(service)

	body := bytes.NewBufferString(`{
		"course_code":"CSC401",
		"session":"2025/2026",
		"semester":"First"
	}`)

	req := authenticatedRequest(
		http.MethodPost,
		"/student/courses",
		body,
	)

	rec := httptest.NewRecorder()

	handler.RegisterCourse(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}

	var response RegisterCourseResponse

	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Message != "course registered successfully" {
		t.Fatalf(
			"unexpected message: %q",
			response.Message,
		)
	}
}

func TestGetStudentCourses_MethodNotAllowed(t *testing.T) {
	handler := NewRegistrationHandler(&mockRegistrationService{})

	req := httptest.NewRequest(
		http.MethodPost,
		"/student/courses",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.GetStudentCourses(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf(
			"expected %d, got %d",
			http.StatusMethodNotAllowed,
			rec.Code,
		)
	}
}

func TestGetStudentCourses_Unauthorized(t *testing.T) {
	handler := NewRegistrationHandler(&mockRegistrationService{})

	req := httptest.NewRequest(
		http.MethodGet,
		"/student/courses",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.GetStudentCourses(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected %d, got %d",
			http.StatusUnauthorized,
			rec.Code,
		)
	}
}

func TestGetStudentCourses_ServiceError(t *testing.T) {
	service := &mockRegistrationService{
		getStudentCoursesFunc: func(studentID string) ([]models.Enrollment, error) {
			return nil, errors.New("failed to load courses")
		},
	}

	handler := NewRegistrationHandler(service)

	req := authenticatedRequest(
		http.MethodGet,
		"/student/courses",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.GetStudentCourses(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected %d, got %d",
			http.StatusInternalServerError,
			rec.Code,
		)
	}
}

func TestGetStudentCourses_Success(t *testing.T) {
	service := &mockRegistrationService{
		getStudentCoursesFunc: func(studentID string) ([]models.Enrollment, error) {
			if studentID != "STU001" {
				t.Fatalf("expected STU001, got %q", studentID)
			}

			return []models.Enrollment{}, nil
		},
	}

	handler := NewRegistrationHandler(service)

	req := authenticatedRequest(
		http.MethodGet,
		"/student/courses",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.GetStudentCourses(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}
}

func TestDropCourse_MethodNotAllowed(t *testing.T) {
	handler := NewRegistrationHandler(&mockRegistrationService{})

	req := httptest.NewRequest(
		http.MethodPost,
		"/student/courses",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.DropCourse(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf(
			"expected %d, got %d",
			http.StatusMethodNotAllowed,
			rec.Code,
		)
	}
}

func TestDropCourse_Unauthorized(t *testing.T) {
	handler := NewRegistrationHandler(&mockRegistrationService{})

	req := httptest.NewRequest(
		http.MethodDelete,
		"/student/courses",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.DropCourse(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected %d, got %d",
			http.StatusUnauthorized,
			rec.Code,
		)
	}
}

func TestDropCourse_MissingCourse(t *testing.T) {
	handler := NewRegistrationHandler(&mockRegistrationService{})

	req := authenticatedRequest(
		http.MethodDelete,
		"/student/courses?session=2025%2F2026&semester=First",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.DropCourse(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}
}

func TestDropCourse_MissingSessionOrSemester(t *testing.T) {
	handler := NewRegistrationHandler(&mockRegistrationService{})

	req := authenticatedRequest(
		http.MethodDelete,
		"/student/courses?course=CSC401",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.DropCourse(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}
}

func TestDropCourse_ServiceError(t *testing.T) {
	service := &mockRegistrationService{
		dropCourseFunc: func(
			studentID,
			courseCode,
			session,
			semester string,
		) error {
			return errors.New("drop failed")
		},
	}

	handler := NewRegistrationHandler(service)

	req := authenticatedRequest(
		http.MethodDelete,
		"/student/courses?course=CSC401&session=2025%2F2026&semester=First",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.DropCourse(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}
}

func TestDropCourse_Success(t *testing.T) {
	service := &mockRegistrationService{
		dropCourseFunc: func(
			studentID,
			courseCode,
			session,
			semester string,
		) error {
			if studentID != "STU001" {
				t.Fatalf("expected STU001, got %q", studentID)
			}

			if courseCode != "CSC401" {
				t.Fatalf("expected CSC401, got %q", courseCode)
			}

			if session != "2025/2026" {
				t.Fatalf("expected 2025/2026, got %q", session)
			}

			if semester != "First" {
				t.Fatalf("expected First, got %q", semester)
			}

			return nil
		},
	}

	handler := NewRegistrationHandler(service)

	req := authenticatedRequest(
		http.MethodDelete,
		"/student/courses?course=CSC401&session=2025%2F2026&semester=First",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.DropCourse(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}
}
