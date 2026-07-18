package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"campuscore/internal/auth"
	"campuscore/internal/middleware"
	"campuscore/internal/models"
)

func TestNewAuthHandler(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")
	sessionMgr := auth.NewSessionManager()

	handler := NewAuthHandler(nil, sessionMgr)

	if handler == nil {
		t.Fatal("expected handler, got nil")
	}

	if handler.userRepo != nil {
		t.Fatal("expected nil userRepo")
	}

	if handler.sessionMgr != sessionMgr {
		t.Fatal("session manager not assigned")
	}
}

func TestLogin_MethodNotAllowed(t *testing.T) {
	repo := &mockUserRepository{}

	sessionMgr := auth.NewSessionManager()

	handler := NewAuthHandler(repo, sessionMgr)

	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected %d, got %d",
			http.StatusMethodNotAllowed,
			res.StatusCode,
		)
	}

	if ct := res.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("expected application/json, got %s", ct)
	}

	body, _ := io.ReadAll(res.Body)

	expected := `{"error": "Method not allowed. Use POST."}`

	if strings.TrimSpace(string(body)) != expected {
		t.Fatalf("expected %s, got %s", expected, string(body))
	}
}

func TestLogin_InvalidJSON(t *testing.T) {
	repo := &mockUserRepository{}

	sessionMgr := auth.NewSessionManager()

	handler := NewAuthHandler(repo, sessionMgr)

	req := httptest.NewRequest(
		http.MethodPost,
		"/login",
		strings.NewReader("{invalid json"),
	)

	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d",
			http.StatusBadRequest,
			res.StatusCode,
		)
	}

	if ct := res.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("expected application/json, got %s", ct)
	}

	body, _ := io.ReadAll(res.Body)

	expected := `{"error": "Invalid request payload syntax."}`

	if strings.TrimSpace(string(body)) != expected {
		t.Fatalf("expected %s, got %s",
			expected,
			string(body),
		)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := &mockUserRepository{
		findByEmailFunc: func(string) (*models.User, error) {
			return nil, errors.New("not found")
		},
		findByIDFunc: func(string) (*models.User, error) {
			return nil, errors.New("not found")
		},
	}

	sessionMgr := auth.NewSessionManager()

	handler := NewAuthHandler(repo, sessionMgr)

	body := `{
		"id":"STU001",
		"password":"password123"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/login",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d",
			http.StatusUnauthorized,
			res.StatusCode,
		)
	}

	if ct := res.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("expected application/json, got %s", ct)
	}

	respBody, _ := io.ReadAll(res.Body)

	expected := `{"error": "Invalid identification numbers or password signature."}`

	if strings.TrimSpace(string(respBody)) != expected {
		t.Fatalf("expected %s, got %s",
			expected,
			string(respBody),
		)
	}
}

func TestLogin_InvalidPassword(t *testing.T) {
	repo := &mockUserRepository{
		findByEmailFunc: func(string) (*models.User, error) {
			return &models.User{
				ID:           "STU001",
				Email:        "student@example.com",
				PasswordHash: "$2a$10$2b2b2b2b2b2b2b2b2b2b2OQ4vN1X6G5Tz7Q5gW2g5sVw1xYzYzYz", // any invalid bcrypt hash is fine
				Role:         models.RoleStudent,
			}, nil
		},
	}

	sessionMgr := auth.NewSessionManager()
	handler := NewAuthHandler(repo, sessionMgr)

	body := `{
		"id":"student@example.com",
		"password":"wrongpassword"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/login",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d",
			http.StatusUnauthorized,
			res.StatusCode,
		)
	}

	resp, _ := io.ReadAll(res.Body)

	expected := `{"error": "Invalid identification numbers or password signature."}`

	if strings.TrimSpace(string(resp)) != expected {
		t.Fatalf("expected %s, got %s",
			expected,
			string(resp),
		)
	}
}

func TestLogin_Success(t *testing.T) {
	hashedPassword, err := auth.HashPassword("password123")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	updateCalled := false

	repo := &mockUserRepository{
		findByEmailFunc: func(string) (*models.User, error) {
			return &models.User{
				ID:           "STU001",
				Email:        "student@example.com",
				PasswordHash: hashedPassword,
				Role:         models.RoleStudent,
			}, nil
		},
		updateLastLoginFunc: func(id string) error {
			updateCalled = true

			if id != "STU001" {
				t.Fatalf("expected STU001, got %s", id)
			}

			return nil
		},
	}

	sessionMgr := auth.NewSessionManager()
	handler := NewAuthHandler(repo, sessionMgr)

	body := `{
		"id":"student@example.com",
		"password":"password123"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/login",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if !updateCalled {
		t.Fatal("expected UpdateLastLogin to be called")
	}

	cookies := res.Cookies()

	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}

	cookie := cookies[0]

	if cookie.Name != "session_token" {
		t.Fatalf("expected session_token cookie, got %s", cookie.Name)
	}

	if cookie.Value == "" {
		t.Fatal("expected non-empty session token")
	}

	if !cookie.HttpOnly {
		t.Fatal("expected HttpOnly cookie")
	}

	var response map[string]any

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if response["message"] != "Authorization successful." {
		t.Fatal("wrong message")
	}

	if response["role"] != "student" {
		t.Fatal("wrong role")
	}

	if response["access_token"] == "" {
		t.Fatal("missing access token")
	}

	if response["token_type"] != "Bearer" {
		t.Fatal("wrong token type")
	}

}

func TestLogout_WithSession(t *testing.T) {
	repo := &mockUserRepository{}

	sessionMgr := auth.NewSessionManager()

	token, err := sessionMgr.CreateSession("STU001", "student")
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	handler := NewAuthHandler(repo, sessionMgr)

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: token,
	})

	rec := httptest.NewRecorder()

	handler.Logout(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected %d, got %d",
			http.StatusOK,
			res.StatusCode,
		)
	}

	// Session should have been revoked.
	if _, err := sessionMgr.ValidateSession(token); err == nil {
		t.Fatal("expected session to be revoked")
	}

	cookies := res.Cookies()

	if len(cookies) == 0 {
		t.Fatal("expected logout cookie")
	}

	if cookies[0].MaxAge != -1 {
		t.Fatalf("expected MaxAge -1, got %d", cookies[0].MaxAge)
	}
}

func TestLogout_WithoutSession(t *testing.T) {
	repo := &mockUserRepository{}

	sessionMgr := auth.NewSessionManager()

	handler := NewAuthHandler(repo, sessionMgr)

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)

	rec := httptest.NewRecorder()

	handler.Logout(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected %d, got %d",
			http.StatusOK,
			res.StatusCode,
		)
	}

	cookies := res.Cookies()

	if len(cookies) == 0 {
		t.Fatal("expected expired cookie")
	}
}
func TestNewStudentHandler(t *testing.T) {
	handler := NewStudentHandler(nil, nil)

	if handler == nil {
		t.Fatal("expected handler")
	}

	if handler.academicService != nil {
		t.Fatal("expected nil academic service")
	}

	if handler.ticketService != nil {
		t.Fatal("expected nil ticket service")
	}
}

func TestWriteJSON(t *testing.T) {
	rec := httptest.NewRecorder()

	payload := map[string]string{
		"message": "ok",
	}

	writeJSON(rec, http.StatusCreated, payload)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("expected %d got %d",
			http.StatusCreated,
			res.StatusCode,
		)
	}

	if ct := res.Header.Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected application/json got %s", ct)
	}
}

func TestGetSession(t *testing.T) {
	expected := &auth.Session{
		UserID: "STU001",
		Role:   "student",
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	ctx := context.WithValue(
		req.Context(),
		middleware.UserContextKey,
		expected,
	)

	req = req.WithContext(ctx)

	session, ok := getSession(req)

	if !ok {
		t.Fatal("expected session")
	}

	if session.UserID != expected.UserID {
		t.Fatal("unexpected user id")
	}
}

func TestRegisterCourse_MethodNotAllowed(t *testing.T) {
	handler := NewStudentHandler(nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/student/register", nil)
	rec := httptest.NewRecorder()

	handler.RegisterCourse(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected %d, got %d",
			http.StatusMethodNotAllowed,
			rec.Code,
		)
	}
}

func TestRegisterCourse_Unauthorized(t *testing.T) {
	handler := NewStudentHandler(nil, nil)

	body := `{
		"course_code":"CSC401",
		"session":"2026/2027",
		"semester":"First"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/student/register",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	handler.RegisterCourse(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d",
			http.StatusUnauthorized,
			rec.Code,
		)
	}
}

func TestRegisterCourse_InvalidJSON(t *testing.T) {
	handler := NewStudentHandler(nil, nil)

	session := &auth.Session{
		UserID: "STU001",
		Role:   "student",
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/student/register",
		strings.NewReader("{invalid"),
	)

	ctx := context.WithValue(
		req.Context(),
		middleware.UserContextKey,
		session,
	)

	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	handler.RegisterCourse(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}
}

func TestRegisterCourse_ServiceError(t *testing.T) {
	academic := &mockAcademicService{
		registerCourseFunc: func(studentID, courseCode, session, semester string) error {
			return errors.New("registration failed")
		},
	}

	handler := NewStudentHandler(academic, nil)

	reqBody := `{
		"course_code":"CSC401",
		"session":"2025/2026",
		"semester":"First"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/student/register",
		strings.NewReader(reqBody),
	)

	session := &auth.Session{
		UserID: "STU001",
	}

	ctx := context.WithValue(
		req.Context(),
		middleware.UserContextKey,
		session,
	)

	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	handler.RegisterCourse(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected %d, got %d",
			http.StatusUnprocessableEntity,
			rec.Code,
		)
	}
}

func TestRegisterCourse_Success(t *testing.T) {
	academic := &mockAcademicService{
		registerCourseFunc: func(studentID, courseCode, session, semester string) error {
			return nil
		},
	}

	handler := NewStudentHandler(academic, nil)

	reqBody := `{
		"course_code":"CSC401",
		"session":"2025/2026",
		"semester":"First"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/student/register",
		strings.NewReader(reqBody),
	)

	session := &auth.Session{
		UserID: "STU001",
	}

	ctx := context.WithValue(
		req.Context(),
		middleware.UserContextKey,
		session,
	)

	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	handler.RegisterCourse(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}
}

func TestSubmitTicket_MethodNotAllowed(t *testing.T) {
	handler := NewStudentHandler(nil, &mockTicketService{})

	req := httptest.NewRequest(http.MethodGet, "/ticket", nil)
	rec := httptest.NewRecorder()

	handler.SubmitTicket(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected %d got %d",
			http.StatusMethodNotAllowed,
			rec.Code)
	}
}

func TestSubmitTicket_Unauthorized(t *testing.T) {
	handler := NewStudentHandler(nil, &mockTicketService{})

	req := httptest.NewRequest(http.MethodPost, "/ticket", nil)
	rec := httptest.NewRecorder()

	handler.SubmitTicket(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d got %d",
			http.StatusUnauthorized,
			rec.Code)
	}
}

func TestSubmitTicket_InvalidJSON(t *testing.T) {
	handler := NewStudentHandler(nil, &mockTicketService{})

	session := &auth.Session{
		UserID: "STU001",
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/ticket",
		strings.NewReader("{invalid"),
	)

	ctx := context.WithValue(
		req.Context(),
		middleware.UserContextKey,
		session,
	)

	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	handler.SubmitTicket(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected %d got %d",
			http.StatusBadRequest,
			rec.Code)
	}
}

func TestSubmitTicket_ServiceError(t *testing.T) {
	ticketSvc := &mockTicketService{
		submitTicketFunc: func(
			ctx context.Context,
			ticket *models.SupportTicket,
		) error {
			return errors.New("service failed")
		},
	}

	handler := NewStudentHandler(nil, ticketSvc)

	body := `{
        "category":"payment",
        "subject":"Receipt",
        "message":"Missing receipt"
    }`

	req := httptest.NewRequest(
		http.MethodPost,
		"/ticket",
		strings.NewReader(body),
	)

	ctx := context.WithValue(
		req.Context(),
		middleware.UserContextKey,
		&auth.Session{UserID: "STU001"},
	)

	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	handler.SubmitTicket(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected %d got %d",
			http.StatusUnprocessableEntity,
			rec.Code)
	}
}

func TestSubmitTicket_Success(t *testing.T) {
	ticketSvc := &mockTicketService{
		submitTicketFunc: func(
			ctx context.Context,
			ticket *models.SupportTicket,
		) error {

			if ticket.StudentID != "STU001" {
				t.Fatal("wrong student")
			}

			return nil
		},
	}

	handler := NewStudentHandler(nil, ticketSvc)

	body := `{
        "category":"payment",
        "subject":"Receipt",
        "message":"Missing receipt"
    }`

	req := httptest.NewRequest(
		http.MethodPost,
		"/ticket",
		strings.NewReader(body),
	)

	ctx := context.WithValue(
		req.Context(),
		middleware.UserContextKey,
		&auth.Session{UserID: "STU001"},
	)

	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	handler.SubmitTicket(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected %d got %d",
			http.StatusCreated,
			rec.Code)
	}
}
