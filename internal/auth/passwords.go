package auth

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

const (
	// DefaultCost sets our work factor workload balance. 
	// 10 offers high security on our ASUS VivoBook workspace while keeping performance crisp.
	DefaultCost = 10
)

// HashPassword takes a plain-text string and turns it into a secure 60-character Bcrypt hash.
func HashPassword(password string) (string, error) {
	if len(password) == 0 {
		return "", errors.New("password cannot be empty")
	}

	// GenerateFromPassword automatically handles random salt generation and injection
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to cryptographically hash password: %w", err)
	}

	return string(hashedBytes), nil
}

// CheckPasswordHash compares an incoming plain-text login attempt against your stored database hash.
// Returns true if it matches perfectly, otherwise returns false.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}