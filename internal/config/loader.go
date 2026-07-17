package config

import (
	"os"
	"strconv"
	"time"
)

// Load creates and returns the application's configuration.
func Load() (*Config, error) {
	cfg := &Config{

		Server: ServerConfig{
			Host:            getEnv("SERVER_HOST", "0.0.0.0"),
			Port:            getEnv("SERVER_PORT", "8080"),
			ReadTimeout:     getDuration("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout:    getDuration("SERVER_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:     getDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
			ShutdownTimeout: getDuration("SERVER_SHUTDOWN_TIMEOUT", 10*time.Second),
		},

		Database: DatabaseConfig{
			Host:              getEnv("DB_HOST", "127.0.0.1"),
			Port:              getEnv("DB_PORT", "5432"),
			User:              getEnv("DB_USER", "postgres"),
			Password:          getEnv("DB_PASSWORD", "postgres"),
			Name:              getEnv("DB_NAME", "campuscore"),
			SSLMode:           getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:      getInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:      getInt("DB_MAX_IDLE_CONNS", 25),
			ConnMaxLifetime:   getDuration("DB_CONN_MAX_LIFETIME", 30*time.Minute),
			ConnMaxIdleTime:   getDuration("DB_CONN_MAX_IDLE_TIME", 5*time.Minute),
		},

		Auth: AuthConfig{
			JWTSecret:          getEnv("JWT_SECRET", "change-me"),
			AccessTokenExpiry:  getDuration("JWT_ACCESS_EXPIRY", 24*time.Hour),
			RefreshTokenExpiry: getDuration("JWT_REFRESH_EXPIRY", 7*24*time.Hour),
			Issuer:             getEnv("JWT_ISSUER", "CampusCore"),
			Audience:           getEnv("JWT_AUDIENCE", "CampusCore"),
			CookieSecure:       getBool("COOKIE_SECURE", false),
			CookieHTTPOnly:     getBool("COOKIE_HTTP_ONLY", true),
			CookieSameSite:     getEnv("COOKIE_SAME_SITE", "Lax"),
		},

		Notification: NotificationConfig{
			WorkerPoolSize: getInt("NOTIFICATION_WORKERS", 5),
			QueueSize:      getInt("NOTIFICATION_QUEUE_SIZE", 100),
			MaxRetries:     getInt("NOTIFICATION_MAX_RETRIES", 3),
			RetryDelay:     getDuration("NOTIFICATION_RETRY_DELAY", 5*time.Second),
			ShutdownTimeout:getDuration("NOTIFICATION_SHUTDOWN_TIMEOUT", 10*time.Second),
		},
	}

	if err := Validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	v, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return v
}

func getBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	v, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return v
}

func getDuration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	v, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}

	return v
}