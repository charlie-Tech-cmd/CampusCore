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
