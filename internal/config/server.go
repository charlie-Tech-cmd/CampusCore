package config

import "time"

type ServerConfig struct {
	Host string
	Port string

	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

// loadServerConfig loads HTTP server configuration from the environment.
func loadServerConfig() ServerConfig {
	return ServerConfig{
		Host:            getEnv("SERVER_HOST", "0.0.0.0"),
		Port:            getEnv("SERVER_PORT", "8080"),
		ReadTimeout:     getEnvAsDuration("SERVER_READ_TIMEOUT", "15s"),
		WriteTimeout:    getEnvAsDuration("SERVER_WRITE_TIMEOUT", "15s"),
		IdleTimeout:     getEnvAsDuration("SERVER_IDLE_TIMEOUT", "60s"),
		ShutdownTimeout: getEnvAsDuration("SERVER_SHUTDOWN_TIMEOUT", "30s"),
	}
}
