package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"campuscore/internal/auth"
	"campuscore/internal/models"
)

func TestAuthenticate_MissingCookie(t *testing.T) {
	sm := auth.NewSessionManager()
	gatekeeper := NewAuthGatekeeper(sm)

	handler := gatekeeper.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestAuthenticate_InvalidSession(t *testing.T) {
	sm := auth.NewSessionManager()
	gatekeeper := NewAuthGatekeeper(sm)

	handler := gatekeeper.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: "invalid-token",
	})

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestAuthenticate_ValidSession(t *testing.T) {
	sm := auth.NewSessionManager()
	gatekeeper := NewAuthGatekeeper(sm)

	token, err := sm.CreateSession("STU001", "student")
	if err != nil {
		t.Fatal(err)
	}

	called := false

	handler := gatekeeper.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true

		session := r.Context().Value(UserContextKey)
		if session == nil {
			t.Fatal("expected session in context")
		}

		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: token,
	})

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if !called {
		t.Fatal("expected handler to be called")
	}

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
}

func TestRequireRole_NoContext(t *testing.T) {
	sm := auth.NewSessionManager()
	gatekeeper := NewAuthGatekeeper(sm)

	handler := gatekeeper.RequireRole(models.RoleStudent)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("handler should not execute")
		}),
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestRequireRole_Authorized(t *testing.T) {
	sm := auth.NewSessionManager()
	gatekeeper := NewAuthGatekeeper(sm)

	session := &auth.Session{
		UserID: "STU001",
		Role:   string(models.RoleStudent),
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(withSession(req, session))

	rr := httptest.NewRecorder()

	called := false

	handler := gatekeeper.RequireRole(models.RoleStudent)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		}),
	)

	handler.ServeHTTP(rr, req)

	if !called {
		t.Fatal("handler should execute")
	}

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
}

func TestRequireRole_Forbidden(t *testing.T) {
	sm := auth.NewSessionManager()
	gatekeeper := NewAuthGatekeeper(sm)

	session := &auth.Session{
		UserID: "STU001",
		Role:   string(models.RoleStudent),
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(withSession(req, session))

	rr := httptest.NewRecorder()

	handler := gatekeeper.RequireRole(models.RoleAdmin)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("handler should not execute")
		}),
	)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected %d got %d", http.StatusForbidden, rr.Code)
	}
}

func withSession(r *http.Request, session *auth.Session) context.Context {
	return context.WithValue(r.Context(), UserContextKey, session)
}