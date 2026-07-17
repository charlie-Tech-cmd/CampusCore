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