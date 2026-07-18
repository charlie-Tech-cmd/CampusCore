package api

import (
	"encoding/json"
	"net/http"

	"campuscore/internal/auth"
)

// RefreshHandler issues a new access token from a valid refresh token.
type RefreshHandler struct{}

// NewRefreshHandler creates a refresh handler.
func NewRefreshHandler() *RefreshHandler {
	return &RefreshHandler{}
}

// RefreshRequest represents the incoming refresh token payload.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RefreshResponse represents the outgoing access token payload.
type RefreshResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// RefreshToken validates a refresh token and issues a new access token.
func (h *RefreshHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)

		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Method not allowed",
		})
		return
	}

	var req RefreshRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request payload",
		})
		return
	}

	claims, err := auth.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)

		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid or expired refresh token",
		})
		return
	}

	accessToken, err := auth.GenerateAccessToken(
		claims.UserID,
		claims.Role,
	)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to generate access token",
		})
		return
	}

	response := RefreshResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int(auth.AccessTokenDuration.Seconds()),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(response)
}
