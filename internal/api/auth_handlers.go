package api

import (
	"encoding/json"
	"net/http"
	"time"

	"campuscore/internal/auth"
	"campuscore/internal/models"
)

// AuthHandler coordinates network data transitions for system access controls
type AuthHandler struct {
	userRepo   models.UserRepository
	sessionMgr *auth.SessionManager
}

// NewAuthHandler instantiates our authentication endpoint controller
func NewAuthHandler(ur models.UserRepository, sm *auth.SessionManager) *AuthHandler {
	return &AuthHandler{
		userRepo:   ur,
		sessionMgr: sm,
	}
}

// LoginRequest defines the expected JSON payload incoming from the login form submission
type LoginRequest struct {
	ID       string `json:"id"` // Matric Number or Staff ID
	Password string `json:"password"`
}

// Login processes incoming login forms and establishes stateful session cookies
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// 1. Enforce strict POST method handling
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte(`{"error": "Method not allowed. Use POST."}`))
		return
	}

	// 2. Decode incoming JSON payload safely
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error": "Invalid request payload syntax."}`))
		return
	}

	// 3. Locate the profile record row inside our repository using the input ID
	user, err := h.userRepo.FindByEmail(req.ID) // Dynamic check fallback matching lookup hooks
	if err != nil {
		// Try lookup by direct primary ID string if email parsing is skipped
		user, err = h.userRepo.FindByID(req.ID)
		if err != nil {
			// Security Strategy: Return a vague error to prevent malicious username enumeration
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error": "Invalid identification numbers or password signature."}`))
			return
		}
	}

	// 4. Validate the plain-text password input against the stored cryptographic hash
	if !auth.CheckPasswordHash(req.Password, user.PasswordHash) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error": "Invalid identification numbers or password signature."}`))
		return
	}

	// 5. Generate a cryptographically strong session token inside our memory matrix
	token, err := h.sessionMgr.CreateSession(user.ID, string(user.Role))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error": "Failed to safely issue tracking session parameters."}`))
		return
	}

	// 6. Update the user's last login tracking timestamp in the background
	_ = h.userRepo.UpdateLastLogin(user.ID)

	// 7. Embed the session token inside a highly protected, secure client cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(5 * time.Minute), // Initial matching boundary tracking tag
		HttpOnly: true,                            // Blocks browser XSS scripts from stealing tokens
		Secure:   false,                           // Set to true in production context over HTTPS
		SameSite: http.SameSiteStrictMode,         // Eliminates Cross-Site Request Forgery (CSRF) attacks
	})

	// 8. Output a clean serialization response profile
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message": "Authorization successful.", "role": "` + string(user.Role) + `"}`))
}

// Logout breaks active token paths to explicitly shut down client access states
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Extract the session cookie to identify what token to destroy
	cookie, err := r.Cookie("session_token")
	if err == nil {
		// Remotely destroy the active session context state tracking row inside memory
		h.sessionMgr.RevokeSession(cookie.Value)
	}

	// Wipe out the client's tracking cookie footprint instantly by forcing expiration parameters
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0), // Immediate expiration command
		HttpOnly: true,
		MaxAge:   -1,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message": "Session tracking token revoked successfully. Disconnected."}`))
}
