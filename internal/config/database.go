package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	// Import the official PostgreSQL driver for standard database/sql execution
	_ "github.com/lib/pq"
)

// DBContainer wraps our active pool instance to pass around cleanly
type DBContainer struct {
	Pool *sql.DB
}

// ConnectPostgres establishes an optimized connection pool based on environment variables
func ConnectPostgres() (*DBContainer, error) {
	// 1. Collect connection details dynamically from environment variables
	host := getEnv("DB_HOST", "127.0.0.1")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	pass := getEnv("DB_PASSWORD", "postgres")
	name := getEnv("DB_NAME", "campuscore")

	// Construct the standard PostgreSQL connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", 
		host, port, user, pass, name)

	// 2. Open the connection handle (does not actually ping the DB yet)
	pool, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %w", err)
	}

	// 3. Configure Connection Pool Parameters for High Concurrency
	// Limits maximum open connections to prevent exhausting database resources
	pool.SetMaxOpenConns(25) 
	// Limits connections sitting idle to save server memory
	pool.SetMaxIdleConns(25) 
	// Connection lifetime rule to recycle old connections safely
	pool.SetConnMaxLifetime(30 * time.Minute)
	pool.SetConnMaxIdleTime(5 * time.Minute)

	// 4. Force a physical connection test (Ping) to ensure database is online
	if err := pool.Ping(); err != nil {
		pool.Close()
		return nil, fmt.Errorf("database unreachable via ping: %w", err)
	}

	log.Println("🔌 [Database Layer] PostgreSQL connection pool initialized successfully.")
	return &DBContainer{Pool: pool}, nil
}

// Helper utility to read environment settings with a fallback value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}