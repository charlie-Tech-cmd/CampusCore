package config

import "time"

// ServerConfig contains all HTTP server settings.
type ServerConfig struct {
	// Host is the network interface the server binds to.
	// Example: "0.0.0.0" or "127.0.0.1"
	Host string

	// Port is the HTTP server port.
	// Example: "8080"
	Port string

	// ReadTimeout is the maximum duration for reading
	// the entire request, including the body.
	ReadTimeout time.Duration

	// WriteTimeout is the maximum duration before
	// timing out writes of the response.
	WriteTimeout time.Duration

	// IdleTimeout is the maximum amount of time to wait
	// for the next request when keep-alives are enabled.
	IdleTimeout time.Duration

	// ShutdownTimeout is the maximum amount of time the
	// server is allowed to gracefully shut down.
	ShutdownTimeout time.Duration
}