package config

import (
	"errors"
	"fmt"
)

// Validate ensures the application configuration is valid before startup.
func Validate(cfg *Config) error {
	if cfg == nil {
		return errors.New("configuration is nil")
	}

	// ------------------------------------------------------------------
	// Server
	// ------------------------------------------------------------------

	if cfg.Server.Host == "" {
		return errors.New("server host cannot be empty")
	}

	if cfg.Server.Port == "" {
		return errors.New("server port cannot be empty")
	}

	if cfg.Server.ReadTimeout <= 0 {
		return errors.New("server read timeout must be greater than zero")
	}

	if cfg.Server.WriteTimeout <= 0 {
		return errors.New("server write timeout must be greater than zero")
	}

	if cfg.Server.IdleTimeout <= 0 {
		return errors.New("server idle timeout must be greater than zero")
	}

	if cfg.Server.ShutdownTimeout <= 0 {
		return errors.New("server shutdown timeout must be greater than zero")
	}

	// ------------------------------------------------------------------
	// Database
	// ------------------------------------------------------------------

	if cfg.Database.Host == "" {
		return errors.New("database host cannot be empty")
	}

	if cfg.Database.Port == "" {
		return errors.New("database port cannot be empty")
	}

	if cfg.Database.User == "" {
		return errors.New("database user cannot be empty")
	}

	if cfg.Database.Name == "" {
		return errors.New("database name cannot be empty")
	}

	if cfg.Database.SSLMode == "" {
		return errors.New("database sslmode cannot be empty")
	}

	if cfg.Database.MaxOpenConns <= 0 {
		return errors.New("database max open connections must be greater than zero")
	}

	if cfg.Database.MaxIdleConns <= 0 {
		return errors.New("database max idle connections must be greater than zero")
	}

	// ------------------------------------------------------------------
	// Authentication
	// ------------------------------------------------------------------

	if cfg.Auth.JWTSecret == "" {
		return errors.New("JWT secret cannot be empty")
	}

	if cfg.Auth.AccessTokenExpiry <= 0 {
		return errors.New("access token expiry must be greater than zero")
	}

	if cfg.Auth.RefreshTokenExpiry <= 0 {
		return errors.New("refresh token expiry must be greater than zero")
	}

	if cfg.Auth.Issuer == "" {
		return errors.New("JWT issuer cannot be empty")
	}

	if cfg.Auth.Audience == "" {
		return errors.New("JWT audience cannot be empty")
	}

	// ------------------------------------------------------------------
	// Notification
	// ------------------------------------------------------------------

	if cfg.Notification.WorkerPoolSize <= 0 {
		return errors.New("notification worker pool size must be greater than zero")
	}

	if cfg.Notification.QueueSize <= 0 {
		return errors.New("notification queue size must be greater than zero")
	}

	if cfg.Notification.MaxRetries < 0 {
		return errors.New("notification max retries cannot be negative")
	}

	if cfg.Notification.RetryDelay <= 0 {
		return errors.New("notification retry delay must be greater than zero")
	}

	if cfg.Notification.ShutdownTimeout <= 0 {
		return errors.New("notification shutdown timeout must be greater than zero")
	}

	return nil
}

// MustValidate panics if the configuration is invalid.
// This is intended for application bootstrap.
func MustValidate(cfg *Config) {
	if err := Validate(cfg); err != nil {
		panic(fmt.Errorf("configuration validation failed: %w", err))
	}
}
