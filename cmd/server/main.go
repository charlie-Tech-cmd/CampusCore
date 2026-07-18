package main

import (
	"context"
	"log"

	"campuscore/internal/config"
	dbmigrate "campuscore/internal/database/migrate"
)

func main() {
	log.Println("Starting CampusCore...")

	// Load and validate application configuration.
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// Connect to PostgreSQL.
	dbContainer, err := ConnectPostgres(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbContainer.Pool.Close()

	// Run database migrations before starting the application.
	if err := dbmigrate.Run(cfg.Database.MigrationURL()); err != nil {
		log.Fatalf("failed to run database migrations: %v", err)
	}

	// Build the HTTP server and background workers.
	server, worker := newServer(dbContainer.Pool)

	// Ensure background workers shut down gracefully.
	defer worker.Stop(context.Background())

	// Start serving requests.
	startServer(server)

	// Wait until the application receives a shutdown signal.
	waitForShutdown(server)
}
