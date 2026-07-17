package api

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"campuscore/internal/auth"
	"campuscore/internal/middleware"
	"campuscore/internal/models"
)

func TestNewLecturerHandler(t *testing.T) {
	engine := &mockGovernanceEngine{}

	handler := NewLecturerHandler(engine)

	if handler == nil {
		t.Fatal("expected handler")
	}

	if handler.govEngine != engine {
		t.Fatal("expected governance engine to be assigned")
	}
}

func TestAdvanceApproval_MethodNotAllowed(t *testing.T) {
	handler := NewLecturerHandler(&mockGovernanceEngine{})

	req := httptest.NewRequest(
		http.MethodGet,
		"/lecturer/advance",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.AdvanceApproval(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf(
			"expected %d, got %d",
			http.StatusMethodNotAllowed,
			rec.Code,
		)
	}
}

func TestAdvanceApproval_Unauthorized(t *testing.T) {
	handler := NewLecturerHandler(&mockGovernanceEngine{})

	body := bytes.NewBufferString(`{
		"course_code":"CSC401"
	}`)

	req := httptest.NewRequest(
		http.MethodPost,
		"/lecturer/advance",
		body,
	)

	rec := httptest.NewRecorder()

	handler.AdvanceApproval(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected %d, got %d",
			http.StatusUnauthorized,
			rec.Code,
		)
	}
}

func TestAdvanceApproval_InvalidJSON(t *testing.T) {
	handler := NewLecturerHandler(&mockGovernanceEngine{})

	body := bytes.NewBufferString("{")

	req := httptest.NewRequest(
		http.MethodPost,
		"/lecturer/advance",
		body,
	)

	session := &auth.Session{
		UserID: "L001",
		Role:   "lecturer",
	}

	ctx := context.WithValue(
		req.Context(),
		middleware.UserContextKey,
		session,
	)

	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	handler.AdvanceApproval(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}
}

func TestAdvanceApproval_GovernanceError(t *testing.T) {
	engine := &mockGovernanceEngine{
		processApprovalAdvanceFunc: func(
			course string,
			role models.UserRole,
			user string,
		) error {
			return errors.New("workflow failed")
		},
	}

	handler := NewLecturerHandler(engine)

	body := bytes.NewBufferString(`{
		"course_code":"CSC401"
	}`)

	req := httptest.NewRequest(
		http.MethodPost,
		"/lecturer/advance",
		body,
	)

	session := &auth.Session{
		UserID: "L001",
		Role:   "lecturer",
	}

	ctx := context.WithValue(
		req.Context(),
		middleware.UserContextKey,
		session,
	)

	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	handler.AdvanceApproval(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected %d, got %d",
			http.StatusUnprocessableEntity,
			rec.Code,
		)
	}
}

func TestAdvanceApproval_Success(t *testing.T) {
	engine := &mockGovernanceEngine{
		processApprovalAdvanceFunc: func(
			course string,
			role models.UserRole,
			user string,
		) error {

			if course != "CSC401" {
				t.Fatal("wrong course")
			}

			if role != models.UserRole("lecturer") {
				t.Fatal("wrong role")
			}

			if user != "L001" {
				t.Fatal("wrong user")
			}

			return nil
		},
	}

	handler := NewLecturerHandler(engine)

	body := bytes.NewBufferString(`{
		"course_code":"CSC401"
	}`)

	req := httptest.NewRequest(
		http.MethodPost,
		"/lecturer/advance",
		body,
	)

	session := &auth.Session{
		UserID: "L001",
		Role:   "lecturer",
	}

	ctx := context.WithValue(
		req.Context(),
		middleware.UserContextKey,
		session,
	)

	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	handler.AdvanceApproval(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}
}

func TestRejectApproval_MethodNotAllowed(t *testing.T) {
	handler := NewLecturerHandler(&mockGovernanceEngine{})

	req := httptest.NewRequest(
		http.MethodGet,
		"/lecturer/reject",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.RejectApproval(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf(
			"expected %d, got %d",
			http.StatusMethodNotAllowed,
			rec.Code,
		)
	}
}

func TestRejectApproval_Unauthorized(t *testing.T) {
	handler := NewLecturerHandler(&mockGovernanceEngine{})

	body := bytes.NewBufferString(`{
		"course_code":"CSC401",
		"remarks":"Missing scores"
	}`)

	req := httptest.NewRequest(
		http.MethodPost,
		"/lecturer/reject",
		body,
	)

	rec := httptest.NewRecorder()

	handler.RejectApproval(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected %d, got %d",
			http.StatusUnauthorized,
			rec.Code,
		)
	}
}

func TestRejectApproval_InvalidJSON(t *testing.T) {
	handler := NewLecturerHandler(&mockGovernanceEngine{})

	body := bytes.NewBufferString("{")

	req := httptest.NewRequest(
		http.MethodPost,
		"/lecturer/reject",
		body,
	)

	session := &auth.Session{
		UserID: "L001",
		Role:   "lecturer",
	}

	ctx := context.WithValue(
		req.Context(),
		middleware.UserContextKey,
		session,
	)

	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	handler.RejectApproval(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}
}

func TestRejectApproval_GovernanceError(t *testing.T) {
	engine := &mockGovernanceEngine{
		processApprovalRejectionFunc: func(
			course string,
			role models.UserRole,
			user string,
			remarks string,
		) error {
			return errors.New("workflow failed")
		},
	}

	handler := NewLecturerHandler(engine)

	body := bytes.NewBufferString(`{
		"course_code":"CSC401",
		"remarks":"Missing scores"
	}`)

	req := httptest.NewRequest(
		http.MethodPost,
		"/lecturer/reject",
		body,
	)

	session := &auth.Session{
		UserID: "L001",
		Role:   "lecturer",
	}

	ctx := context.WithValue(
		req.Context(),
		middleware.UserContextKey,
		session,
	)

	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	handler.RejectApproval(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf(
			"expected %d, got %d",
			http.StatusUnprocessableEntity,
			rec.Code,
		)
	}
}

func TestRejectApproval_Success(t *testing.T) {
	engine := &mockGovernanceEngine{
		processApprovalRejectionFunc: func(
			course string,
			role models.UserRole,
			user string,
			remarks string,
		) error {

			if course != "CSC401" {
				t.Fatal("wrong course")
			}

			if role != models.UserRole("lecturer") {
				t.Fatal("wrong role")
			}

			if user != "L001" {
				t.Fatal("wrong user")
			}

			if remarks != "Missing scores" {
				t.Fatal("wrong remarks")
			}

			return nil
		},
	}

	handler := NewLecturerHandler(engine)

	body := bytes.NewBufferString(`{
		"course_code":"CSC401",
		"remarks":"Missing scores"
	}`)

	req := httptest.NewRequest(
		http.MethodPost,
		"/lecturer/reject",
		body,
	)

	session := &auth.Session{
		UserID: "L001",
		Role:   "lecturer",
	}

	ctx := context.WithValue(
		req.Context(),
		middleware.UserContextKey,
		session,
	)

	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	handler.RejectApproval(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}
}
