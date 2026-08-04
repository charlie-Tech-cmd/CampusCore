package repository

import (
	"database/sql"
	"errors"
	"time"

	"campuscore/internal/models"
)

type PostgresAdmissionRepository struct {
	db *sql.DB
}

func NewPostgresAdmissionRepository(
	db *sql.DB,
) *PostgresAdmissionRepository {

	return &PostgresAdmissionRepository{
		db: db,
	}
}

func (r *PostgresAdmissionRepository) SubmitApplication(
	application *models.AdmissionApplication,
) error {

	// TODO:
	// INSERT INTO admission_applications ...

	return nil
}

func (r *PostgresAdmissionRepository) FindByApplicationNo(
	applicationNo string,
) (*models.AdmissionApplication, error) {

	// TODO:
	// SELECT ...

	return nil, sql.ErrNoRows
}

func (r *PostgresAdmissionRepository) ListApplications() ([]models.AdmissionApplication, error) {

	return []models.AdmissionApplication{}, nil
}

func (r *PostgresAdmissionRepository) ApproveApplication(
	applicationNo string,
) error {

	now := time.Now()

	_ = now

	// TODO:
	// UPDATE status='approved'

	return nil
}

func (r *PostgresAdmissionRepository) RejectApplication(
	applicationNo string,
) error {

	if applicationNo == "" {
		return errors.New("application number required")
	}

	// TODO:
	// UPDATE status='rejected'

	return nil
}
