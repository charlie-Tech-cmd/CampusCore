package main

import (
	"context"
	"log"

	"campuscore/internal/config"
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

	// Build the HTTP server.
	server, worker := newServer(dbContainer.Pool)
	defer worker.Stop(context.Background())

	startServer(server)

	waitForShutdown(server)
}