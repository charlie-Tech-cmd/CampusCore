package api

import (
	"encoding/json"
	"net/http"

	"campuscore/internal/auth"
)

type RefreshHandler struct{}

func NewRefreshHandler() *RefreshHandler {
	return &RefreshHandler{}
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func (h *RefreshHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {

	// Only POST is allowed.
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)

		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Method not allowed",
		})
		return
	}

	// Decode request.
	var req RefreshRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request",
		})
		return
	}

	// Validate refresh token.
	claims, err := auth.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)

		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid or expired refresh token",
		})
		return
	}

	// Generate a new access token.
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

	// Respond with the new access token.
	response := RefreshResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int(auth.AccessTokenDuration.Seconds()),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(response)
}
