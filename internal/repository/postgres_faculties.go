package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"campuscore/internal/models"
)

// PostgresFacultyRepository manages faculty persistence.
type PostgresFacultyRepository struct {
	db *sql.DB
}

// NewPostgresFacultyRepository creates a faculty repository.
func NewPostgresFacultyRepository(db *sql.DB) *PostgresFacultyRepository {
	return &PostgresFacultyRepository{
		db: db,
	}
}

// Create inserts a faculty.
func (r *PostgresFacultyRepository) Create(faculty *models.Faculty) error {
	query := `
		INSERT INTO faculties
		(
			id,
			code,
			name,
			description,
			is_active
		)
		VALUES ($1, $2, $3, $4, $5);
	`

	_, err := r.db.Exec(
		query,
		faculty.ID,
		faculty.Code,
		faculty.Name,
		faculty.Description,
		faculty.IsActive,
	)

	if err != nil {
		return fmt.Errorf("failed to create faculty: %w", err)
	}

	return nil
}

// FindByID retrieves a faculty by ID.
func (r *PostgresFacultyRepository) FindByID(id int) (*models.Faculty, error) {
	query := `
		SELECT
			id,
			code,
			name,
			description,
			is_active,
			created_at,
			updated_at
		FROM faculties
		WHERE id = $1
		LIMIT 1;
	`

	var faculty models.Faculty

	err := r.db.QueryRow(query, id).Scan(
		&faculty.ID,
		&faculty.Code,
		&faculty.Name,
		&faculty.Description,
		&faculty.IsActive,
		&faculty.CreatedAt,
		&faculty.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("faculty not found")
		}
		return nil, fmt.Errorf("failed to retrieve faculty: %w", err)
	}

	return &faculty, nil
}

// FindByCode retrieves a faculty by code.
func (r *PostgresFacultyRepository) FindByCode(code string) (*models.Faculty, error) {
	query := `
		SELECT
			id,
			code,
			name,
			description,
			is_active,
			created_at,
			updated_at
		FROM faculties
		WHERE code = $1
		LIMIT 1;
	`

	var faculty models.Faculty

	err := r.db.QueryRow(query, code).Scan(
		&faculty.ID,
		&faculty.Code,
		&faculty.Name,
		&faculty.Description,
		&faculty.IsActive,
		&faculty.CreatedAt,
		&faculty.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("faculty not found")
		}
		return nil, fmt.Errorf("failed to retrieve faculty: %w", err)
	}

	return &faculty, nil
}

// List retrieves all faculties.
func (r *PostgresFacultyRepository) List() ([]models.Faculty, error) {
	query := `
		SELECT
			id,
			code,
			name,
			description,
			is_active,
			created_at,
			updated_at
		FROM faculties
		ORDER BY name;
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve faculties: %w", err)
	}
	defer rows.Close()

	var faculties []models.Faculty

	for rows.Next() {
		var faculty models.Faculty

		if err := rows.Scan(
			&faculty.ID,
			&faculty.Code,
			&faculty.Name,
			&faculty.Description,
			&faculty.IsActive,
			&faculty.CreatedAt,
			&faculty.UpdatedAt,
		); err != nil {
			return nil, err
		}

		faculties = append(faculties, faculty)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return faculties, nil
}

// Update modifies a faculty.
func (r *PostgresFacultyRepository) Update(faculty *models.Faculty) error {
	query := `
		UPDATE faculties
		SET
			code = $2,
			name = $3,
			description = $4,
			is_active = $5
		WHERE id = $1;
	`

	result, err := r.db.Exec(
		query,
		faculty.ID,
		faculty.Code,
		faculty.Name,
		faculty.Description,
		faculty.IsActive,
	)

	if err != nil {
		return fmt.Errorf("failed to update faculty: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("faculty not found")
	}

	return nil
}

// Delete removes a faculty.
func (r *PostgresFacultyRepository) Delete(id int) error {
	query := `
		DELETE FROM faculties
		WHERE id = $1;
	`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete faculty: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("faculty not found")
	}

	return nil
}
