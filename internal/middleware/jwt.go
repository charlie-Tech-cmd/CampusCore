package middleware

import (
	"context"
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
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDContextKey, claims.UserID)
		ctx = context.WithValue(ctx, UserRoleContextKey, claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CurrentUserID(r *http.Request) string {
	id, _ := r.Context().Value(UserIDContextKey).(string)
	return id
}

func CurrentUserRole(r *http.Request) string {
	role, _ := r.Context().Value(UserRoleContextKey).(string)
	return role
}
