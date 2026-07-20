package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"campuscore/internal/models"
)

// PostgresDepartmentRepository manages department persistence.
type PostgresDepartmentRepository struct {
	db *sql.DB
}

// NewPostgresDepartmentRepository creates a department repository.
func NewPostgresDepartmentRepository(db *sql.DB) *PostgresDepartmentRepository {
	return &PostgresDepartmentRepository{
		db: db,
	}
}

// Create inserts a new department.
func (r *PostgresDepartmentRepository) Create(department *models.Department) error {
	query := `
		INSERT INTO departments (id, name, faculty_id)
		VALUES ($1, $2, $3);
	`

	_, err := r.db.Exec(
		query,
		department.ID,
		department.Name,
		department.FacultyID,
	)

	if err != nil {
		return fmt.Errorf("failed to create department: %w", err)
	}

	return nil
}

// FindByID retrieves a department by its ID.
func (r *PostgresDepartmentRepository) FindByID(id int) (*models.Department, error) {
	query := `
		SELECT id, name, faculty_id
		FROM departments
		WHERE id = $1
		LIMIT 1;
	`

	var department models.Department

	err := r.db.QueryRow(query, id).Scan(
		&department.ID,
		&department.Name,
		&department.FacultyID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("department not found")
		}
		return nil, fmt.Errorf("failed to retrieve department: %w", err)
	}

	return &department, nil
}

// FindAll retrieves all departments.
func (r *PostgresDepartmentRepository) FindAll() ([]models.Department, error) {
	query := `
		SELECT id, name, faculty_id
		FROM departments
		ORDER BY name;
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve departments: %w", err)
	}
	defer rows.Close()

	var departments []models.Department

	for rows.Next() {
		var department models.Department

		if err := rows.Scan(
			&department.ID,
			&department.Name,
			&department.FacultyID,
		); err != nil {
			return nil, err
		}

		departments = append(departments, department)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return departments, nil
}

// Update modifies an existing department.
func (r *PostgresDepartmentRepository) Update(department *models.Department) error {
	query := `
		UPDATE departments
		SET
			name = $2,
			faculty_id = $3
		WHERE id = $1;
	`

	result, err := r.db.Exec(
		query,
		department.ID,
		department.Name,
		department.FacultyID,
	)

	if err != nil {
		return fmt.Errorf("failed to update department: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("department not found")
	}

	return nil
}

// Delete removes a department.
func (r *PostgresDepartmentRepository) Delete(id int) error {
	query := `
		DELETE FROM departments
		WHERE id = $1;
	`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete department: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("department not found")
	}

	return nil
}
