package config

import "time"

// NotificationConfig contains settings for the background
// notification processing engine.
type NotificationConfig struct {
	// WorkerPoolSize specifies the number of worker goroutines
	// that process notification jobs concurrently.
	WorkerPoolSize int

	// QueueSize specifies the maximum number of notification
	// jobs that can be queued before producers block or reject
	// new jobs.
	QueueSize int

	// MaxRetries defines how many times a failed notification
	// should be retried before being marked as failed.
	MaxRetries int

	// RetryDelay specifies the delay between retry attempts.
	RetryDelay time.Duration

	// ShutdownTimeout specifies how long the notification
	// engine is allowed to complete pending jobs during
	// graceful shutdown.
	ShutdownTimeout time.Duration
}

// loadNotificationConfig loads notification worker configuration from the environment.
func loadNotificationConfig() NotificationConfig {
	return NotificationConfig{
		WorkerPoolSize: getEnvAsInt("WORKER_POOL_SIZE", 5),
		QueueSize:      getEnvAsInt("NOTIFICATION_QUEUE_SIZE", 100),
		MaxRetries:     getEnvAsInt("NOTIFICATION_MAX_RETRIES", 3),
		RetryDelay:     getEnvAsDuration("NOTIFICATION_RETRY_DELAY", "5s"),
		ShutdownTimeout: getEnvAsDuration(
			"NOTIFICATION_SHUTDOWN_TIMEOUT",
			"30s",
		),
	}
}
