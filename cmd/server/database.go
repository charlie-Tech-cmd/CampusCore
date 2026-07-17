package main

import (
	"database/sql"
	"fmt"
	"log"

	"campuscore/internal/config"

	_ "github.com/lib/pq"
)

// DBContainer wraps the PostgreSQL connection pool.
type DBContainer struct {
	Pool *sql.DB
}

// ConnectPostgres creates and validates the PostgreSQL connection pool.
func ConnectPostgres(cfg config.DatabaseConfig) (*DBContainer, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.SSLMode,
	)

	log.Printf(
		"Connecting to PostgreSQL (host=%s port=%s database=%s sslmode=%s)",
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
	)
	pool, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("opening database connection: %w", err)
	}

	pool.SetMaxOpenConns(cfg.MaxOpenConns)
	pool.SetMaxIdleConns(cfg.MaxIdleConns)
	pool.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	pool.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	if err := pool.Ping(); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	log.Println("Database connected.")

	return &DBContainer{
		Pool: pool,
	}, nil
}
