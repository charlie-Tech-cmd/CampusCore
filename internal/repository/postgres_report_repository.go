package repository

import (
	"database/sql"

	"campuscore/internal/models"
)

type PostgresReportRepository struct {
	db *sql.DB
}

func NewPostgresReportRepository(
	db *sql.DB,
) *PostgresReportRepository {

	return &PostgresReportRepository{
		db: db,
	}
}

func (r *PostgresReportRepository) GetDashboardSummary() (
	*models.DashboardSummary,
	error,
) {

	return &models.DashboardSummary{}, nil
}
