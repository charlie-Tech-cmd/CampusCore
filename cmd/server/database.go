package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func mustConnectDB() *sql.DB {
	connStr := "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=campuscore sslmode=disable"

	log.Println(connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(15 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Database connected.")

	return db
}