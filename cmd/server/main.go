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
	server, worker := newServer(dbContainer.Pool, cfg.Server)

	// Stop the worker with a bounded deadline before closing the database.
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		worker.Stop(ctx)
	}()

	// Start serving requests and retain any unexpected listener error.
	serverErrors := startServer(server)

	// Wait until a shutdown signal or an unexpected server error occurs.
	waitForShutdown(server, serverErrors)
}
