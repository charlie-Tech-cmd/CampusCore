package main

import (
	"log"
	"net/http"
)

func startServer(server *http.Server) {
	go func() {
		log.Println("Server listening on http://localhost:8080")

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
}