package repository

import (
	"campuscore/internal/models"
	"database/sql"
	"errors"
	"fmt"
)

// PostgresEnrollmentRepository implements models.EnrollmentRepository.
type PostgresEnrollmentRepository struct {
	db *sql.DB
}

// NewPostgresEnrollmentRepository creates a new enrollment repository.
func NewPostgresEnrollmentRepository(db *sql.DB) *PostgresEnrollmentRepository {
	return &PostgresEnrollmentRepository{db: db}
}

// Register registers a student for a course.
func (r *PostgresEnrollmentRepository) Register(enrollment *models.Enrollment) error {
	query := `
		INSERT INTO student_courses
		(student_id, course_code, session, semester, status)
		VALUES ($1, $2, $3, $4, $5);
	`

	_, err := r.db.Exec(
		query,
		enrollment.StudentID,
		enrollment.CourseCode,
		enrollment.Session,
		enrollment.Semester,
		enrollment.Status,
	)

	if err != nil {
		return fmt.Errorf("failed to register course: %w", err)
	}

	return nil
}

// FindByStudent returns every enrollment belonging to a student.
func (r *PostgresEnrollmentRepository) FindByStudent(studentID string) ([]models.Enrollment, error) {
	query := `
		SELECT
			id,
			student_id,
			course_code,
			session,
			semester,
			status,
			created_at,
			updated_at
		FROM student_courses
		WHERE student_id = $1
		ORDER BY created_at DESC;
	`

	rows, err := r.db.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enrollments []models.Enrollment

	for rows.Next() {
		var e models.Enrollment

		if err := rows.Scan(
			&e.ID,
			&e.StudentID,
			&e.CourseCode,
			&e.Session,
			&e.Semester,
			&e.Status,
			&e.CreatedAt,
			&e.UpdatedAt,
		); err != nil {
			return nil, err
		}

		enrollments = append(enrollments, e)
	}

	return enrollments, rows.Err()
}

// FindByCourse returns all students enrolled in a course.
func (r *PostgresEnrollmentRepository) FindByCourse(courseCode string) ([]models.Enrollment, error) {
	query := `
		SELECT
			id,
			student_id,
			course_code,
			session,
			semester,
			status,
			created_at,
			updated_at
		FROM student_courses
		WHERE course_code = $1;
	`

	rows, err := r.db.Query(query, courseCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enrollments []models.Enrollment

	for rows.Next() {
		var e models.Enrollment

		if err := rows.Scan(
			&e.ID,
			&e.StudentID,
			&e.CourseCode,
			&e.Session,
			&e.Semester,
			&e.Status,
			&e.CreatedAt,
			&e.UpdatedAt,
		); err != nil {
			return nil, err
		}

		enrollments = append(enrollments, e)
	}

	return enrollments, rows.Err()
}

// UpdateStatus updates an enrollment status.
func (r *PostgresEnrollmentRepository) UpdateStatus(studentID, courseCode, status string) error {
	query := `
		UPDATE student_courses
		SET status = $3
		WHERE student_id = $1
		  AND course_code = $2;
	`

	result, err := r.db.Exec(query, studentID, courseCode, status)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("enrollment not found")
	}

	return nil
}

// Delete removes an enrollment.
func (r *PostgresEnrollmentRepository) Delete(studentID, courseCode string) error {
	query := `
		DELETE FROM student_courses
		WHERE student_id = $1
		  AND course_code = $2;
	`

	result, err := r.db.Exec(query, studentID, courseCode)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("enrollment not found")
	}

	return nil
}
