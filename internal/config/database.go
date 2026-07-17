package config

import "time"

// DatabaseConfig contains PostgreSQL connection and
// connection pool settings.
type DatabaseConfig struct {
	// Connection Settings
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string

	// Connection Pool Settings
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// loadDatabaseConfig loads PostgreSQL configuration from the environment.
func loadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:     getEnv("DB_HOST", "127.0.0.1"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		Name:     getEnv("DB_NAME", "campuscore"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),

		MaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 10),

		ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", "30m"),
		ConnMaxIdleTime: getEnvAsDuration("DB_CONN_MAX_IDLE_TIME", "5m"),
	}
}
