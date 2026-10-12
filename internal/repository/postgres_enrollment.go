package repository

import (
	"campuscore/internal/models"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

var ErrDuplicateEnrollment = errors.New(
	"student is already registered for this course",
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
func (r *PostgresEnrollmentRepository) Register(
	enrollment *models.Enrollment,
) error {
	query := `
        INSERT INTO student_courses (
            student_id,
            course_code,
            session_id,
            status
        )
        SELECT
            $1,
            $2,
            id,
            $5
        FROM academic_sessions
        WHERE session_name = $3
          AND semester = $4
          AND registration_open = true
          AND is_active = true;
    `

	result, err := r.db.Exec(
		query,
		enrollment.StudentID,
		enrollment.CourseCode,
		enrollment.Session,
		enrollment.Semester,
		enrollment.Status,
	)
	if err != nil {
		var pqErr *pq.Error

		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return ErrDuplicateEnrollment
		}

		return fmt.Errorf("failed to register course: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(
			"failed to verify course registration: %w",
			err,
		)
	}

	if rowsAffected == 0 {
		return fmt.Errorf(
			"academic session %s / %s is not active or registration is closed",
			enrollment.Session,
			enrollment.Semester,
		)
	}

	return nil
}

// FindByStudent returns every enrollment belonging to a student.
func (r *PostgresEnrollmentRepository) FindByStudent(studentID string) ([]models.Enrollment, error) {
	query := `
SELECT
    sc.id,
    sc.student_id,
    sc.course_code,
    c.credit_units,
    sc.session_id,
    asess.session_name,
    asess.semester::text,
    sc.status,
    sc.registered_at
			FROM student_courses sc
		JOIN courses c
			ON c.code = sc.course_code
		JOIN academic_sessions asess
			ON asess.id = sc.session_id
		WHERE sc.student_id = $1
		ORDER BY sc.registered_at DESC;
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
			&e.CreditUnits,
			&e.SessionID,
			&e.Session,
			&e.Semester,
			&e.Status,
			&e.CreatedAt,
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
    sc.id,
    sc.student_id,
    sc.course_code,
    c.credit_units,
    sc.session_id,
    asess.session_name,
    asess.semester::text,
    sc.status,
    sc.registered_at
	FROM student_courses sc
		JOIN courses c
		ON c.code = sc.course_code
		JOIN academic_sessions asess
		ON asess.id = sc.session_id
		WHERE sc.course_code = $1
		ORDER BY sc.registered_at DESC;
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
			&e.CreditUnits,
			&e.SessionID,
			&e.Session,
			&e.Semester,
			&e.Status,
			&e.CreatedAt,
		); err != nil {
			return nil, err
		}
		enrollments = append(enrollments, e)
	}

	return enrollments, rows.Err()
}

// UpdateStatus updates an enrollment status.

func (r *PostgresEnrollmentRepository) UpdateStatus(
	studentID, courseCode, session, semester, status string,
) error {
	query := `
		UPDATE student_courses
		SET status = $5
		WHERE student_id = $1
		  AND course_code = $2
		  AND session_id = (
		      SELECT id
		      FROM academic_sessions
		      WHERE session_name = $3
		        AND semester = $4
		  );
	`

	result, err := r.db.Exec(
		query,
		studentID,
		courseCode,
		session,
		semester,
		status,
	)
	if err != nil {
		return fmt.Errorf("update enrollment status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check updated enrollment rows: %w", err)
	}
	if rows == 0 {
		return errors.New("enrollment not found")
	}

	return nil
}

// Delete removes an enrollment.
func (r *PostgresEnrollmentRepository) Delete(
	studentID,
	courseCode,
	session,
	semester string,
) error {
	query := `
        DELETE FROM student_courses
        WHERE student_id = $1
          AND course_code = $2
          AND session_id = (
              SELECT id
              FROM academic_sessions
              WHERE session_name = $3
                AND semester = $4
          );
    `

	result, err := r.db.Exec(
		query,
		studentID,
		courseCode,
		session,
		semester,
	)
	if err != nil {
		return fmt.Errorf("failed to delete enrollment: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(
			"failed to verify enrollment deletion: %w",
			err,
		)
	}

	if rows == 0 {
		return errors.New("enrollment not found")
	}

	return nil
}
