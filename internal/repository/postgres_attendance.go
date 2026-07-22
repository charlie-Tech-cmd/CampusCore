package repository

import (
	"campuscore/internal/models"
	"database/sql"
	"errors"
	"fmt"
)

type PostgresAttendanceRepository struct {
	db *sql.DB
}

func NewPostgresAttendanceRepository(db *sql.DB) *PostgresAttendanceRepository {
	return &PostgresAttendanceRepository{
		db: db,
	}
}

// Create stores a new attendance record.
func (r *PostgresAttendanceRepository) Create(
	record *models.Attendance,
) error {

	query := `
		INSERT INTO attendance (
			student_id,
			course_code,
			lecturer_id,
			session,
			semester,
			class_date,
			status,
			marked_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8);
	`

	_, err := r.db.Exec(
		query,
		record.StudentID,
		record.CourseCode,
		record.LecturerID,
		record.Session,
		record.Semester,
		record.ClassDate,
		record.Status,
		record.MarkedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create attendance: %w", err)
	}

	return nil
}

// FindByID returns an attendance record by ID.
func (r *PostgresAttendanceRepository) FindByID(
	id int,
) (*models.Attendance, error) {

	query := `
		SELECT
			id,
			student_id,
			course_code,
			lecturer_id,
			session,
			semester,
			class_date,
			status,
			marked_at
		FROM attendance
		WHERE id = $1
		LIMIT 1;
	`

	var record models.Attendance

	err := r.db.QueryRow(query, id).Scan(
		&record.ID,
		&record.StudentID,
		&record.CourseCode,
		&record.LecturerID,
		&record.Session,
		&record.Semester,
		&record.ClassDate,
		&record.Status,
		&record.MarkedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("attendance record not found")
		}
		return nil, err
	}

	return &record, nil
}

// GetAll returns all attendance records.
func (r *PostgresAttendanceRepository) GetAll() ([]models.Attendance, error) {
	query := `
		SELECT
			id,
			student_id,
			course_code,
			lecturer_id,
			session,
			semester,
			class_date,
			status,
			marked_at
		FROM attendance
		ORDER BY class_date DESC;
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve attendance records: %w", err)
	}
	defer rows.Close()

	var records []models.Attendance

	for rows.Next() {
		var record models.Attendance

		if err := rows.Scan(
			&record.ID,
			&record.StudentID,
			&record.CourseCode,
			&record.LecturerID,
			&record.Session,
			&record.Semester,
			&record.ClassDate,
			&record.Status,
			&record.MarkedAt,
		); err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

// GetByStudent returns attendance records for a student.
func (r *PostgresAttendanceRepository) GetByStudent(
	studentID string,
) ([]models.Attendance, error) {

	query := `
		SELECT
			id,
			student_id,
			course_code,
			lecturer_id,
			session,
			semester,
			class_date,
			status,
			marked_at
		FROM attendance
		WHERE student_id = $1
		ORDER BY class_date DESC;
	`

	rows, err := r.db.Query(query, studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve student attendance: %w", err)
	}
	defer rows.Close()

	var records []models.Attendance

	for rows.Next() {
		var record models.Attendance

		if err := rows.Scan(
			&record.ID,
			&record.StudentID,
			&record.CourseCode,
			&record.LecturerID,
			&record.Session,
			&record.Semester,
			&record.ClassDate,
			&record.Status,
			&record.MarkedAt,
		); err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	return records, rows.Err()
}

// GetByCourse returns attendance records for a course.
func (r *PostgresAttendanceRepository) GetByCourse(
	courseCode string,
) ([]models.Attendance, error) {

	query := `
		SELECT
			id,
			student_id,
			course_code,
			lecturer_id,
			session,
			semester,
			class_date,
			status,
			marked_at
		FROM attendance
		WHERE course_code = $1
		ORDER BY class_date DESC;
	`

	rows, err := r.db.Query(query, courseCode)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve course attendance: %w", err)
	}
	defer rows.Close()

	var records []models.Attendance

	for rows.Next() {
		var record models.Attendance

		if err := rows.Scan(
			&record.ID,
			&record.StudentID,
			&record.CourseCode,
			&record.LecturerID,
			&record.Session,
			&record.Semester,
			&record.ClassDate,
			&record.Status,
			&record.MarkedAt,
		); err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	return records, rows.Err()
}

// GetByLecturer returns attendance records marked by a lecturer.
func (r *PostgresAttendanceRepository) GetByLecturer(
	lecturerID string,
) ([]models.Attendance, error) {

	query := `
		SELECT
			id,
			student_id,
			course_code,
			lecturer_id,
			session,
			semester,
			class_date,
			status,
			marked_at
		FROM attendance
		WHERE lecturer_id = $1
		ORDER BY class_date DESC;
	`

	rows, err := r.db.Query(query, lecturerID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve lecturer attendance: %w", err)
	}
	defer rows.Close()

	var records []models.Attendance

	for rows.Next() {
		var record models.Attendance

		if err := rows.Scan(
			&record.ID,
			&record.StudentID,
			&record.CourseCode,
			&record.LecturerID,
			&record.Session,
			&record.Semester,
			&record.ClassDate,
			&record.Status,
			&record.MarkedAt,
		); err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	return records, rows.Err()
}

// Delete removes an attendance record.
func (r *PostgresAttendanceRepository) Delete(id int) error {

	query := `
		DELETE FROM attendance
		WHERE id = $1;
	`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete attendance record: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("attendance record not found")
	}

	return nil
}

// Update modifies an attendance record.
func (r *PostgresAttendanceRepository) Update(
	record *models.Attendance,
) error {

	query := `
		UPDATE attendance
		SET
			student_id = $2,
			course_code = $3,
			lecturer_id = $4,
			session = $5,
			semester = $6,
			class_date = $7,
			status = $8,
			marked_at = $9
		WHERE id = $1;
	`

	result, err := r.db.Exec(
		query,
		record.ID,
		record.StudentID,
		record.CourseCode,
		record.LecturerID,
		record.Session,
		record.Semester,
		record.ClassDate,
		record.Status,
		record.MarkedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update attendance: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("attendance record not found")
	}

	return nil
}
