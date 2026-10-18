package config

import (
	"strings"
	"testing"
	"time"
)

func validConfigForValidationTests() *Config {
	return &Config{
		Server: ServerConfig{
			Host:            "localhost",
			Port:            "8080",
			ReadTimeout:     15 * time.Second,
			WriteTimeout:    15 * time.Second,
			IdleTimeout:     60 * time.Second,
			ShutdownTimeout: 30 * time.Second,
		},
		Database: DatabaseConfig{
			Host:         "localhost",
			Port:         "5432",
			User:         "campuscore",
			Password:     "test-password",
			Name:         "campuscore_test",
			SSLMode:      "disable",
			MaxOpenConns: 10,
			MaxIdleConns: 5,
		},
		Auth: AuthConfig{
			JWTSecret:          "test-secret-with-at-least-32-bytes",
			AccessTokenExpiry:  15 * time.Minute,
			RefreshTokenExpiry: 168 * time.Hour,
			Issuer:             "CampusCore",
			Audience:           "CampusCore",
		},
		Notification: NotificationConfig{
			WorkerPoolSize:  2,
			QueueSize:       10,
			MaxRetries:      3,
			RetryDelay:      5 * time.Second,
			ShutdownTimeout: 30 * time.Second,
		},
	}
}

func TestValidateJWTSecret(t *testing.T) {
	tests := []struct {
		name        string
		secret      string
		wantErrText string
	}{
		{
			name:        "empty secret",
			secret:      "",
			wantErrText: "JWT secret cannot be empty",
		},
		{
			name:        "whitespace-only secret",
			secret:      "   \t ",
			wantErrText: "JWT secret cannot be empty",
		},
		{
			name:        "default placeholder",
			secret:      "change-me-before-production",
			wantErrText: "JWT secret must not use the default placeholder",
		},
		{
			name:        "short secret",
			secret:      "short-secret",
			wantErrText: "JWT secret must be at least 32 bytes",
		},
		{
			name:   "valid secret",
			secret: "a-random-looking-test-secret-of-32-bytes-or-more",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := validConfigForValidationTests()
			cfg.Auth.JWTSecret = tt.secret

			err := Validate(cfg)

			if tt.wantErrText == "" {
				if err != nil {
					t.Fatalf("Validate() returned unexpected error: %v", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("Validate() expected error containing %q, got nil", tt.wantErrText)
			}

			if !strings.Contains(err.Error(), tt.wantErrText) {
				t.Fatalf("Validate() error = %q, want it to contain %q", err.Error(), tt.wantErrText)
			}
		})
	}
}
