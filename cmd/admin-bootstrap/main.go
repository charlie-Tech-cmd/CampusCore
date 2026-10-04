package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"syscall"

	"campuscore/internal/auth"
	"campuscore/internal/config"

	_ "github.com/lib/pq"
	"golang.org/x/term"
)

const (
	adminID        = "CHARLES-001"
	adminFirstName = "Charles"
	adminSurname   = "Owoicho"
	adminEmail     = "charles@gmail.com"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fatal(err)
	}

	dbContainer, err := connectPostgres(cfg.Database)
	if err != nil {
		fatal(err)
	}
	defer dbContainer.Close()

	if err := createAdmin(dbContainer, adminID, adminFirstName, adminSurname, adminEmail); err != nil {
		fatal(err)
	}
}

func connectPostgres(cfg config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("opening database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("connecting to database: %w", err)
	}

	return db, nil
}

func createAdmin(db *sql.DB, id, firstName, surname, email string) error {
	var exists bool

	err := db.QueryRow(
		`SELECT EXISTS(
			SELECT 1 FROM users WHERE id = $1 OR email = $2
		)`,
		id,
		email,
	).Scan(&exists)
	if err != nil {
		return fmt.Errorf("checking existing user: %w", err)
	}

	if exists {
		return fmt.Errorf("a user with ID %q or email %q already exists", id, email)
	}

	password, err := readPassword("Admin password: ")
	if err != nil {
		return err
	}

	confirmPassword, err := readPassword("Confirm password: ")
	if err != nil {
		return err
	}

	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	if password != confirmPassword {
		return fmt.Errorf("passwords do not match")
	}

	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return fmt.Errorf("hashing password: %w", err)
	}

	_, err = db.ExecContext(
		context.Background(),
		`INSERT INTO users (
			id,
			surname,
			first_name,
			email,
			password_hash,
			role
		) VALUES ($1, $2, $3, $4, $5, $6)`,
		id,
		surname,
		firstName,
		email,
		passwordHash,
		"admin",
	)
	if err != nil {
		return fmt.Errorf("creating admin user: %w", err)
	}

	fmt.Printf("Admin user %s created successfully.\n", id)
	return nil
}

func readPassword(prompt string) (string, error) {
	fmt.Print(prompt)

	password, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()

	if err != nil {
		return "", fmt.Errorf("reading password: %w", err)
	}

	return strings.TrimSpace(string(password)), nil
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

var _ = bufio.NewReader
