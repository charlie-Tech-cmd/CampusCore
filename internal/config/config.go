package config

// Config represents the complete application configuration.
// It aggregates all configuration sections into a single object
// that can be passed throughout the application.
type Config struct {
	Server       ServerConfig
	Database     DatabaseConfig
	Auth         AuthConfig
	Notification NotificationConfig
}
