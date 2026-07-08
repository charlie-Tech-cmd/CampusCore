package middleware

import (
	"context"
	"net/http"
	"campuscore/internal/auth"
	"campuscore/internal/models"
)

// contextKey defines a custom, unexported type to prevent namespace collisions in request context
type contextKey string

const (
	UserContextKey contextKey = "user_session"
)

// AuthGatekeeper wraps our route processing to enforce global identity authorization
type AuthGatekeeper struct {
	sessionMgr *auth.SessionManager
}

// NewAuthGatekeeper instantiates our gateway structure with the session engine handle
func NewAuthGatekeeper(sm *auth.SessionManager) *AuthGatekeeper {
	return &AuthGatekeeper{sessionMgr: sm}
}

// Authenticate checks for a session cookie and verifies its active status
func (ag *AuthGatekeeper) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Extract the secure session cookie
		cookie, err := r.Cookie("session_token")
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error": "Authentication required. Missing session cookie."}`))
			return
		}

		// 2. Validate the token against our running session storage tracking matrix
		session, err := ag.sessionMgr.ValidateSession(cookie.Value)
		if err != nil {
			// Session could be missing or expired due to our 15-minute inactivity rule
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error": "` + err.Error() + `"}`))
			return
		}

		// 3. Inject the active session variables down into the request context container
		ctx := context.WithValue(r.Context(), UserContextKey, session)
		
		// 4. Pass the authenticated request forward along the execution chain
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole guards access by evaluating if the active session matches allowed access levels
func (ag *AuthGatekeeper) RequireRole(allowedRoles ...models.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Retrieve the session context injected during the Authenticate step
			sessionVal := r.Context().Value(UserContextKey)
			if sessionVal == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"error": "Unauthorized processing state."}`))
				return
			}

			activeSession := sessionVal.(*auth.Session)
			roleMatched := false

			// 2. Loop through the checking arguments to confirm eligibility
			for _, role := range allowedRoles {
				if string(role) == activeSession.Role {
					roleMatched = true
					break
				}
			}

			// 3. Reject execution if the role requirement isn't satisfied
			if !roleMatched {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				_, _ = w.Write([]byte(`{"error": "Access Denied: Your account role does not have permission to execute this operation."}`))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}