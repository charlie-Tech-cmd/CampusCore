package repository

import (
	"database/sql"
)

type PostgresTimetableRepository struct {
	db *sql.DB
}

func NewPostgresTimetableRepository(db *sql.DB) *PostgresTimetableRepository {
	return &PostgresTimetableRepository{
		db: db,
	}
}
