package main

import (
	"context"
	"log"
)

func main() {
	log.Println("Starting CampusCore...")

	db := mustConnectDB()
	defer db.Close()

	server, worker := newServer(db)
	defer worker.Stop(context.Background())

	startServer(server)

	waitForShutdown(server)
}