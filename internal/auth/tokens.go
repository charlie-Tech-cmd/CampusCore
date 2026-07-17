package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

// Session represents an active authenticated user state inside application RAM
type Session struct {
	UserID    string
	Role      string
	ExpiresAt time.Time
}

// SessionManager coordinates internal state controls with a thread-safe mutex lock
type SessionManager struct {
	mu       sync.RWMutex
	sessions map[string]Session
	duration time.Duration
}

// NewSessionManager instantiates our internal session tracking matrix with the 15-minute rule
func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]Session),
		duration: 15 * time.Minute, // Our strict 15-minute inactivity security window
	}
}

// GenerateToken creates a cryptographically strong random hex token string
func (sm *SessionManager) GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", errors.New("failed to generate secure security token")
	}
	return hex.EncodeToString(bytes), nil
}

// CreateSession generates a session attached to a token, establishing an active state
func (sm *SessionManager) CreateSession(userID string, role string) (string, error) {
	token, err := sm.GenerateToken()
	if err != nil {
		return "", err
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Assign the user mapping along with a rolling 15-minute expiration timestamp
	sm.sessions[token] = Session{
		UserID:    userID,
		Role:      role,
		ExpiresAt: time.Now().Add(sm.duration),
	}

	return token, nil
}

// ValidateSession verifies a token and dynamically extends its lifecycle if active
func (sm *SessionManager) ValidateSession(token string) (*Session, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session, exists := sm.sessions[token]
	if !exists {
		return nil, errors.New("session does not exist or has been revoked")
	}

	// Inactivity Timeout Verification Gate
	if time.Now().After(session.ExpiresAt) {
		delete(sm.sessions, token) // Clean up expired session footprint from memory
		return nil, errors.New("session has expired due to 15-minute inactivity timeout")
	}

	// Rolling Inactivity Window: Extend expiration time by another 15 minutes on active use
	session.ExpiresAt = time.Now().Add(sm.duration)
	sm.sessions[token] = session

	return &session, nil
}

// RevokeSession manually destroys an active session state (used on Explicit Logout)
func (sm *SessionManager) RevokeSession(token string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.sessions, token)
}
