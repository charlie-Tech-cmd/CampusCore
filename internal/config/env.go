package config

import (
	"os"
	"strconv"
	"time"
)

// getEnv returns the value of an environment variable,
// or the provided fallback if it is not set.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// getEnvAsInt returns an integer environment variable,
// or the fallback value if parsing fails.
func getEnvAsInt(key string, fallback int) int {
	value := getEnv(key, "")

	if value == "" {
		return fallback
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return intValue
}

// getEnvAsBool returns a boolean environment variable,
// or the fallback value if parsing fails.
func getEnvAsBool(key string, fallback bool) bool {
	value := getEnv(key, "")

	if value == "" {
		return fallback
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return boolValue
}

// getEnvAsDuration returns a time.Duration parsed from
// an environment variable, or the fallback duration.
func getEnvAsDuration(key, fallback string) time.Duration {
	value := getEnv(key, fallback)

	duration, err := time.ParseDuration(value)
	if err != nil {
		duration, _ = time.ParseDuration(fallback)
	}

	return duration
}
