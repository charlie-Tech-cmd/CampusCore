package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims defines the payload stored inside every JWT.
type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`

	jwt.RegisteredClaims
}

// Token lifetimes.
const (
	AccessTokenDuration  = 15 * time.Minute
	RefreshTokenDuration = 7 * 24 * time.Hour
)

// jwtSecret loads the signing secret from the environment.
func jwtSecret() ([]byte, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET environment variable is not set")
	}

	return []byte(secret), nil
}

// GenerateAccessToken creates a signed access token.
func GenerateAccessToken(userID, role string) (string, error) {
	return generateToken(userID, role, AccessTokenDuration)
}

// GenerateRefreshToken creates a signed refresh token.
func GenerateRefreshToken(userID, role string) (string, error) {
	return generateToken(userID, role, RefreshTokenDuration)
}

// generateToken creates a signed JWT.
func generateToken(userID, role string, duration time.Duration) (string, error) {
	secret, err := jwtSecret()
	if err != nil {
		return "", err
	}

	now := time.Now()

	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "CampusCore",
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secret)
}

// ValidateAccessToken validates an access token.
func ValidateAccessToken(tokenString string) (*Claims, error) {
	return validateToken(tokenString)
}

// ValidateRefreshToken validates a refresh token.
func ValidateRefreshToken(tokenString string) (*Claims, error) {
	return validateToken(tokenString)
}

// validateToken verifies signature and expiration.
func validateToken(tokenString string) (*Claims, error) {
	secret, err := jwtSecret()
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}

			return secret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
