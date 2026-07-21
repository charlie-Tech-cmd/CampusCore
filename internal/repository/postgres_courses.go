package repository

import (
	"campuscore/internal/models"
	"database/sql"
	"errors"
	"fmt"
)

// PostgresCourseRepository implements models.CourseRepository.
type PostgresCourseRepository struct {
	db *sql.DB
}

// NewPostgresCourseRepository creates a new course repository.
func NewPostgresCourseRepository(db *sql.DB) *PostgresCourseRepository {
	return &PostgresCourseRepository{db: db}
}

// Create inserts a new course.
func (r *PostgresCourseRepository) Create(course *models.Course) error {
	query := `
		INSERT INTO courses (
			code,
			title,
			description,
			credit_units,
			department_id,
			level,
			semester,
			max_capacity,
			current_enrolled,
			is_active
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10);
	`

	_, err := r.db.Exec(
		query,
		course.Code,
		course.Title,
		course.Description,
		course.CreditUnits,
		course.DepartmentID,
		course.Level,
		course.Semester,
		course.MaxCapacity,
		course.CurrentEnrolled,
		course.IsActive,
	)

	if err != nil {
		return fmt.Errorf("failed to create course: %w", err)
	}

	return nil
}

// FindByCode returns a course by its code.
func (r *PostgresCourseRepository) FindByCode(code string) (*models.Course, error) {
	query := `
		SELECT
			code,
			title,
			description,
			credit_units,
			department_id,
			level,
			semester,
			max_capacity,
			current_enrolled,
			is_active
		FROM courses
		WHERE code = $1
		LIMIT 1;
	`

	var course models.Course

	err := r.db.QueryRow(query, code).Scan(
		&course.Code,
		&course.Title,
		&course.Description,
		&course.CreditUnits,
		&course.DepartmentID,
		&course.Level,
		&course.Semester,
		&course.MaxCapacity,
		&course.CurrentEnrolled,
		&course.IsActive,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("course %s not found", code)
		}
		return nil, err
	}

	return &course, nil
}

// Update modifies an existing course.
func (r *PostgresCourseRepository) Update(course *models.Course) error {
	query := `
		UPDATE courses
		SET
			title = $2,
			description = $3,
			credit_units = $4,
			department_id = $5,
			level = $6,
			semester = $7,
			max_capacity = $8,
			current_enrolled = $9,
			is_active = $10
		WHERE code = $1;
	`

	result, err := r.db.Exec(
		query,
		course.Code,
		course.Title,
		course.Description,
		course.CreditUnits,
		course.DepartmentID,
		course.Level,
		course.Semester,
		course.MaxCapacity,
		course.CurrentEnrolled,
		course.IsActive,
	)
	if err != nil {
		return fmt.Errorf("failed to update course: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("course not found")
	}

	return nil
}
