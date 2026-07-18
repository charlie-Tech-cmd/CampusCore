package api

import (
	"encoding/json"
	"net/http"
	"time"

	"campuscore/internal/auth"
	"campuscore/internal/middleware"
	"campuscore/internal/models"
)

// AuthHandler coordinates network data transitions for system access controls.
type AuthHandler struct {
	userRepo   models.UserRepository
	sessionMgr *auth.SessionManager
}

// NewAuthHandler instantiates our authentication endpoint controller.
func NewAuthHandler(ur models.UserRepository, sm *auth.SessionManager) *AuthHandler {
	return &AuthHandler{
		userRepo:   ur,
		sessionMgr: sm,
	}
}

// LoginRequest defines the expected JSON payload.
type LoginRequest struct {
	ID       string `json:"id"` // Matric Number or Staff ID
	Password string `json:"password"`
}

// RegisterResponse defines the response returned after a successful registration.
type RegisterResponse struct {
	Message string `json:"message"`
	UserID  string `json:"user_id"`
	Role    string `json:"role"`
}

// Login authenticates a user and issues both a session cookie and JWT.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Only POST is allowed.
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte(`{"error": "Method not allowed. Use POST."}`))
		return
	}
	// Decode request.
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error": "Invalid request payload syntax."}`))
		return
	}
	// Lookup user.
	user, err := h.userRepo.FindByEmail(req.ID)
	if err != nil {
		user, err = h.userRepo.FindByID(req.ID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error": "Invalid identification numbers or password signature."}`))
			return
		}
	}

	// Verify password.
	if !auth.CheckPasswordHash(req.Password, user.PasswordHash) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error": "Invalid identification numbers or password signature."}`))
		return
	}
	// Create session.
	sessionToken, err := h.sessionMgr.CreateSession(user.ID, string(user.Role))
	if err != nil {
		http.Error(w, `{"error":"Failed to safely issue tracking session parameters."}`, http.StatusInternalServerError)
		return
	}

	// Create JWT.
	accessToken, err := auth.GenerateAccessToken(user.ID, string(user.Role))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error": "Failed to generate access token."}`))
		return
	}

	// Update login timestamp.
	_ = h.userRepo.UpdateLastLogin(user.ID)

	// Set session cookie.
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		Expires:  time.Now().Add(5 * time.Minute),
		HttpOnly: true,
		Secure:   false, // Change to true in production
		SameSite: http.SameSiteStrictMode,
	})

	// Return JWT.
	response := map[string]any{
		"message":      "Authorization successful.",
		"role":         string(user.Role),
		"access_token": accessToken,
		"token_type":   "Bearer",
		"expires_in":   int(auth.AccessTokenDuration.Seconds()),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error": "Failed to encode response."}`))
		return
	}
}

// Register creates a new CampusCore user account.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Only POST requests are allowed.
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)

		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Method not allowed. Use POST.",
		})
		return
	}

	// Decode incoming request.
	var req models.UserOnboarding

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request payload.",
		})
		return
	}

	// Validate password confirmation.
	if req.Password != req.ConfirmPassword {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Passwords do not match.",
		})
		return
	}

	// Ensure the ID does not already exist.
	if _, err := h.userRepo.FindByID(req.ID); err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)

		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "User already exists.",
		})
		return
	}

	// Ensure the email does not already exist.
	if _, err := h.userRepo.FindByEmail(req.Email); err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)

		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Email already registered.",
		})
		return
	}

	// Hash the password.
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to secure password.",
		})
		return
	}

	// Build the user entity.
	user := &models.User{
		ID:           req.ID,
		Surname:      req.Surname,
		FirstName:    req.FirstName,
		MiddleName:   req.MiddleName,
		Email:        req.Email,
		Phone:        req.Phone,
		PasswordHash: hashedPassword,
		Role:         req.Role,
		DepartmentID: req.DepartmentID,
		Level:        req.Level,
	}

	// Persist the new user.
	if err := h.userRepo.Create(user); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to create user account.",
		})
		return
	}

	// Respond.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(RegisterResponse{
		Message: "Account created successfully.",
		UserID:  user.ID,
		Role:    string(user.Role),
	})
}

// Logout revokes the current session.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err == nil {
		h.sessionMgr.RevokeSession(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "Session tracking token revoked successfully. Disconnected.",
	})
}

// Me returns the currently authenticated user's profile.
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	sessionVal := r.Context().Value(middleware.UserContextKey)
	if sessionVal == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error": "Unauthorized"}`))
		return
	}

	session := sessionVal.(*auth.Session)

	user, err := h.userRepo.FindByID(session.UserID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error": "User not found"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(user)
}
