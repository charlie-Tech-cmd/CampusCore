package api

import (
	"encoding/json"
	"net/http"

	"campuscore/internal/models"
	"campuscore/internal/services"
)

// UserHandler handles user profile requests.
type UserHandler struct {
	service *services.UserService
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(
	service *services.UserService,
) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// GetProfile returns a user's profile.
func (h *UserHandler) GetProfile(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID := r.URL.Query().Get("id")

	profile, err := h.service.GetProfile(userID)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(profile)
}

// UpdateProfile updates a user's profile.
func (h *UserHandler) UpdateProfile(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPut {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	if err := h.service.UpdateProfile(&user); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "profile updated successfully",
	})
}
