package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"campuscore/internal/auth"
)

const (
	UserIDContextKey   contextKey = "user_id"
	UserRoleContextKey contextKey = "user_role"
)

// JWTAuth validates JWT Bearer tokens.
func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := auth.ValidateAccessToken(token)
		if err != nil {
			log.Printf("JWT validation failed: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		session := &auth.Session{
			UserID: claims.UserID,
			Role:   claims.Role,
		}

		ctx := context.WithValue(r.Context(), UserContextKey, session)
		ctx = context.WithValue(ctx, UserIDContextKey, claims.UserID)
		ctx = context.WithValue(ctx, UserRoleContextKey, claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CurrentUserID(r *http.Request) string {
	session, ok := r.Context().Value(UserContextKey).(*auth.Session)
	if !ok || session == nil {
		return ""
	}

	return session.UserID
}

func CurrentUserRole(r *http.Request) string {
	session, ok := r.Context().Value(UserContextKey).(*auth.Session)
	if !ok || session == nil {
		return ""
	}

	return session.Role
}
