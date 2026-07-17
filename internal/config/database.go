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