package auth

import (
	"testing"
	"time"
)

// Test that a new session manager is initialized correctly.
func TestNewSessionManager(t *testing.T) {
	sm := NewSessionManager()

	if sm == nil {
		t.Fatal("expected SessionManager, got nil")
	}

	if sm.sessions == nil {
		t.Fatal("expected sessions map to be initialized")
	}

	if len(sm.sessions) != 0 {
		t.Fatal("expected empty sessions map")
	}

	if sm.duration != 15*time.Minute {
		t.Fatalf("expected duration %v, got %v", 15*time.Minute, sm.duration)
	}
}

// Test secure token generation.
func TestGenerateToken(t *testing.T) {
	sm := NewSessionManager()

	token, err := sm.GenerateToken()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if token == "" {
		t.Fatal("expected non-empty token")
	}

	// 32 random bytes = 64 hex characters
	if len(token) != 64 {
		t.Fatalf("expected token length 64, got %d", len(token))
	}
}

// Ensure every generated token is unique.
func TestGenerateTokenUniqueness(t *testing.T) {
	sm := NewSessionManager()

	token1, err := sm.GenerateToken()
	if err != nil {
		t.Fatal(err)
	}

	token2, err := sm.GenerateToken()
	if err != nil {
		t.Fatal(err)
	}

	if token1 == token2 {
		t.Fatal("expected unique tokens")
	}
}

// Test successful session creation.
func TestCreateSession(t *testing.T) {
	sm := NewSessionManager()

	token, err := sm.CreateSession("STU001", "student")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if token == "" {
		t.Fatal("expected session token")
	}

	if len(sm.sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sm.sessions))
	}

	session, ok := sm.sessions[token]
	if !ok {
		t.Fatal("expected session stored in map")
	}

	if session.UserID != "STU001" {
		t.Fatalf("expected user STU001, got %s", session.UserID)
	}

	if session.Role != "student" {
		t.Fatalf("expected role student, got %s", session.Role)
	}
}

// Test successful session validation.
func TestValidateSession(t *testing.T) {
	sm := NewSessionManager()

	token, err := sm.CreateSession("STU001", "student")
	if err != nil {
		t.Fatal(err)
	}

	session, err := sm.ValidateSession(token)
	if err != nil {
		t.Fatalf("expected valid session, got %v", err)
	}

	if session.UserID != "STU001" {
		t.Fatalf("expected STU001, got %s", session.UserID)
	}

	if session.Role != "student" {
		t.Fatalf("expected student, got %s", session.Role)
	}
}

// Test validation of an unknown token.
func TestValidateSession_InvalidToken(t *testing.T) {
	sm := NewSessionManager()

	_, err := sm.ValidateSession("invalid-token")

	if err == nil {
		t.Fatal("expected error for invalid session")
	}
}

// Test session revocation.
func TestRevokeSession(t *testing.T) {
	sm := NewSessionManager()

	token, err := sm.CreateSession("STU001", "student")
	if err != nil {
		t.Fatal(err)
	}

	sm.RevokeSession(token)

	if len(sm.sessions) != 0 {
		t.Fatal("expected session to be removed")
	}

	_, err = sm.ValidateSession(token)
	if err == nil {
		t.Fatal("expected revoked session to be invalid")
	}
}

// Test expired session removal.
func TestExpiredSession(t *testing.T) {
	sm := NewSessionManager()

	token, err := sm.CreateSession("STU001", "student")
	if err != nil {
		t.Fatal(err)
	}

	// Force expiration.
	session := sm.sessions[token]
	session.ExpiresAt = time.Now().Add(-time.Minute)
	sm.sessions[token] = session

	_, err = sm.ValidateSession(token)

	if err == nil {
		t.Fatal("expected expired session")
	}

	if _, exists := sm.sessions[token]; exists {
		t.Fatal("expected expired session to be deleted")
	}
}

// Test that successful validation extends the session expiry.
func TestValidateSessionExtendsExpiration(t *testing.T) {
	sm := NewSessionManager()

	token, err := sm.CreateSession("STU001", "student")
	if err != nil {
		t.Fatal(err)
	}

	oldExpiry := sm.sessions[token].ExpiresAt

	time.Sleep(20 * time.Millisecond)

	_, err = sm.ValidateSession(token)
	if err != nil {
		t.Fatal(err)
	}

	newExpiry := sm.sessions[token].ExpiresAt

	if !newExpiry.After(oldExpiry) {
		t.Fatal("expected expiration time to be extended")
	}
}