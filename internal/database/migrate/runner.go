package migrate

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Run applies all pending database migrations.
func Run(databaseURL string) error {
	m, err := migrate.New(
		"file://database/migrations",
		databaseURL,
	)
	if err != nil {
		return fmt.Errorf("create migration instance: %w", err)
	}

	defer func() {
		srcErr, dbErr := m.Close()

		if srcErr != nil {
			fmt.Printf("migration source close error: %v\n", srcErr)
		}

		if dbErr != nil {
			fmt.Printf("migration database close error: %v\n", dbErr)
		}
	}()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}

		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}
