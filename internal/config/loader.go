package config

import "github.com/joho/godotenv"

// Load loads the application configuration from the environment.
func Load() (*Config, error) {
	// Load .env if present. Environment variables already set
	// in the OS take precedence.
	_ = godotenv.Load()

	cfg := &Config{
		Server:       loadServerConfig(),
		Database:     loadDatabaseConfig(),
		Auth:         loadAuthConfig(),
		Notification: loadNotificationConfig(),
	}

	if err := Validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
