package main

import (
	"log"
	"net/http"
)

func startServer(server *http.Server) {
	go func() {
		log.Printf("Server listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
}
