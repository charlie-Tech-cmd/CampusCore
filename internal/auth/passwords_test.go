package auth

import "testing"

// Test that a valid password is successfully hashed.
func TestHashPassword(t *testing.T) {
	password := "Student@123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if hash == "" {
		t.Fatal("expected a hash, got an empty string")
	}

	if hash == password {
		t.Fatal("hashed password should not equal the original password")
	}
}

// Test that an empty password returns an error.
func TestHashPassword_EmptyPassword(t *testing.T) {
	hash, err := HashPassword("")

	if err == nil {
		t.Fatal("expected an error for an empty password")
	}

	if hash != "" {
		t.Fatal("expected empty hash when password is empty")
	}
}

// Test that the correct password validates successfully.
func TestCheckPasswordHash_CorrectPassword(t *testing.T) {
	password := "Student@123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	if !CheckPasswordHash(password, hash) {
		t.Fatal("expected password validation to succeed")
	}
}

// Test that an incorrect password fails validation.
func TestCheckPasswordHash_WrongPassword(t *testing.T) {
	password := "Student@123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	if CheckPasswordHash("WrongPassword", hash) {
		t.Fatal("expected password validation to fail")
	}
}

// Test empty password input.
func TestCheckPasswordHash_EmptyPassword(t *testing.T) {
	password := "Student@123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	if CheckPasswordHash("", hash) {
		t.Fatal("expected empty password to fail validation")
	}
}

// Test empty hash input.
func TestCheckPasswordHash_EmptyHash(t *testing.T) {
	if CheckPasswordHash("Student@123", "") {
		t.Fatal("expected empty hash to fail validation")
	}
}

// Test invalid hash input.
func TestCheckPasswordHash_InvalidHash(t *testing.T) {
	if CheckPasswordHash("Student@123", "this-is-not-a-bcrypt-hash") {
		t.Fatal("expected invalid hash to fail validation")
	}
}

// Test that bcrypt generates unique hashes because of random salts.
func TestHashPassword_GeneratesUniqueHashes(t *testing.T) {
	password := "Student@123"

	hash1, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	if hash1 == hash2 {
		t.Fatal("expected hashes to be different because bcrypt uses random salts")
	}

	if !CheckPasswordHash(password, hash1) {
		t.Fatal("hash1 should validate")
	}

	if !CheckPasswordHash(password, hash2) {
		t.Fatal("hash2 should validate")
	}
}
