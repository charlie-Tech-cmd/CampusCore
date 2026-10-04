package auth

import "testing"

func TestAccessAndRefreshTokens(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-for-jwt")

	accessToken, err := GenerateAccessToken("STU001", "student")
	if err != nil {
		t.Fatalf("failed to generate access token: %v", err)
	}

	refreshToken, err := GenerateRefreshToken("STU001", "student")
	if err != nil {
		t.Fatalf("failed to generate refresh token: %v", err)
	}

	accessClaims, err := ValidateAccessToken(accessToken)
	if err != nil {
		t.Fatalf("failed to validate access token: %v", err)
	}

	if accessClaims.TokenType != AccessTokenType {
		t.Fatalf("expected access token type %q, got %q", AccessTokenType, accessClaims.TokenType)
	}

	if accessClaims.UserID != "STU001" {
		t.Fatalf("expected user ID STU001, got %q", accessClaims.UserID)
	}

	if accessClaims.Role != "student" {
		t.Fatalf("expected role student, got %q", accessClaims.Role)
	}

	refreshClaims, err := ValidateRefreshToken(refreshToken)
	if err != nil {
		t.Fatalf("failed to validate refresh token: %v", err)
	}

	if refreshClaims.TokenType != RefreshTokenType {
		t.Fatalf("expected refresh token type %q, got %q", RefreshTokenType, refreshClaims.TokenType)
	}

	if refreshClaims.UserID != "STU001" {
		t.Fatalf("expected user ID STU001, got %q", refreshClaims.UserID)
	}

	if refreshClaims.Role != "student" {
		t.Fatalf("expected role student, got %q", refreshClaims.Role)
	}
}

func TestTokenTypesCannotBeUsedInterchangeably(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-for-jwt")

	accessToken, err := GenerateAccessToken("STU001", "student")
	if err != nil {
		t.Fatalf("failed to generate access token: %v", err)
	}

	refreshToken, err := GenerateRefreshToken("STU001", "student")
	if err != nil {
		t.Fatalf("failed to generate refresh token: %v", err)
	}

	if _, err := ValidateRefreshToken(accessToken); err == nil {
		t.Fatal("expected access token to be rejected by refresh-token validation")
	}

	if _, err := ValidateAccessToken(refreshToken); err == nil {
		t.Fatal("expected refresh token to be rejected by access-token validation")
	}
}
