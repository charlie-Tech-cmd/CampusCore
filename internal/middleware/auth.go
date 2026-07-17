package middleware

import (
	"context"
	"net/http"
	"strings"

	"campuscore/internal/auth"
	"campuscore/internal/models"
)

// contextKey defines a custom, unexported type to prevent namespace collisions in request context.
type contextKey string

const (
	UserContextKey contextKey = "user_session"
)

// AuthGatekeeper wraps our route processing to enforce global identity authorization.
type AuthGatekeeper struct {
	sessionMgr *auth.SessionManager
}

// NewAuthGatekeeper instantiates our gateway structure with the session engine handle.
func NewAuthGatekeeper(sm *auth.SessionManager) *AuthGatekeeper {
	return &AuthGatekeeper{
		sessionMgr: sm,
	}
}

// Authenticate supports BOTH Session Cookies and JWT Bearer authentication.
func (ag *AuthGatekeeper) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// ----------------------------------------------------
		// 1. Try Session Cookie Authentication
		// ----------------------------------------------------
		if cookie, err := r.Cookie("session_token"); err == nil {

			session, err := ag.sessionMgr.ValidateSession(cookie.Value)
			if err == nil {
				ctx := context.WithValue(r.Context(), UserContextKey, session)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		// ----------------------------------------------------
		// 2. Try JWT Authentication
		// ----------------------------------------------------
		authHeader := r.Header.Get("Authorization")

		if strings.HasPrefix(authHeader, "Bearer ") {

			token := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := auth.ValidateAccessToken(token)
			if err == nil {

				session := &auth.Session{
					UserID: claims.UserID,
					Role:   claims.Role,
				}

				ctx := context.WithValue(r.Context(), UserContextKey, session)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		// ----------------------------------------------------
		// 3. Authentication Failed
		// ----------------------------------------------------
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)

		_, _ = w.Write([]byte(`{"error":"authentication required"}`))
	})
}

// RequireRole guards access by evaluating if the active session matches allowed access levels.
func (ag *AuthGatekeeper) RequireRole(allowedRoles ...models.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			sessionVal := r.Context().Value(UserContextKey)
			if sessionVal == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			session := sessionVal.(*auth.Session)

			for _, role := range allowedRoles {
				if string(role) == session.Role {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "Forbidden", http.StatusForbidden)
		})
	}
}
