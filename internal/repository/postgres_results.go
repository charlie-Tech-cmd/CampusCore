package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"campuscore/internal/models"
)

// PostgresResultRepository manages result persistence.
type PostgresResultRepository struct {
	db *sql.DB
}

// NewPostgresResultRepository creates a new result repository.
func NewPostgresResultRepository(db *sql.DB) *PostgresResultRepository {
	return &PostgresResultRepository{
		db: db,
	}
}

// Submit inserts a new academic result.
func (r *PostgresResultRepository) Submit(result *models.Result) error {
	query := `
		INSERT INTO results (
			id,
			student_id,
			course_code,
			session,
			semester,
			score,
			grade,
			grade_point,
			credit_units,
			approved,
			approved_by
		)
		VALUES (
			$1,$2,$3,$4,$5,
			$6,$7,$8,$9,$10,$11
		);
	`

	_, err := r.db.Exec(
		query,
		result.ID,
		result.StudentID,
		result.CourseCode,
		result.Session,
		result.Semester,
		result.Score,
		result.Grade,
		result.GradePoint,
		result.CreditUnits,
		result.Approved,
		result.ApprovedBy,
	)

	if err != nil {
		return fmt.Errorf("failed to submit result: %w", err)
	}

	return nil
}

// FindByStudent returns all results belonging to a student.
func (r *PostgresResultRepository) FindByStudent(studentID string) ([]models.Result, error) {
	query := `
		SELECT
			id,
			student_id,
			course_code,
			session,
			semester,
			score,
			grade,
			grade_point,
			credit_units,
			approved,
			approved_by,
			created_at,
			updated_at
		FROM results
		WHERE student_id = $1
		ORDER BY session DESC, semester;
	`

	rows, err := r.db.Query(query, studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve student results: %w", err)
	}
	defer rows.Close()

	var results []models.Result

	for rows.Next() {
		var result models.Result

		if err := rows.Scan(
			&result.ID,
			&result.StudentID,
			&result.CourseCode,
			&result.Session,
			&result.Semester,
			&result.Score,
			&result.Grade,
			&result.GradePoint,
			&result.CreditUnits,
			&result.Approved,
			&result.ApprovedBy,
			&result.CreatedAt,
			&result.UpdatedAt,
		); err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// FindByCourse returns all results for a course.
func (r *PostgresResultRepository) FindByCourse(courseCode string) ([]models.Result, error) {
	query := `
		SELECT
			id,
			student_id,
			course_code,
			session,
			semester,
			score,
			grade,
			grade_point,
			credit_units,
			approved,
			approved_by,
			created_at,
			updated_at
		FROM results
		WHERE course_code = $1
		ORDER BY student_id;
	`

	rows, err := r.db.Query(query, courseCode)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve course results: %w", err)
	}
	defer rows.Close()

	var results []models.Result

	for rows.Next() {
		var result models.Result

		if err := rows.Scan(
			&result.ID,
			&result.StudentID,
			&result.CourseCode,
			&result.Session,
			&result.Semester,
			&result.Score,
			&result.Grade,
			&result.GradePoint,
			&result.CreditUnits,
			&result.Approved,
			&result.ApprovedBy,
			&result.CreatedAt,
			&result.UpdatedAt,
		); err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// Update modifies an existing result.
func (r *PostgresResultRepository) Update(result *models.Result) error {
	query := `
		UPDATE results
		SET
			score = $2,
			grade = $3,
			grade_point = $4,
			credit_units = $5,
			approved = $6,
			approved_by = $7
		WHERE id = $1;
	`

	res, err := r.db.Exec(
		query,
		result.ID,
		result.Score,
		result.Grade,
		result.GradePoint,
		result.CreditUnits,
		result.Approved,
		result.ApprovedBy,
	)

	if err != nil {
		return fmt.Errorf("failed to update result: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("result not found")
	}

	return nil
}

// Delete removes a result.
func (r *PostgresResultRepository) Delete(id int) error {
	query := `
		DELETE FROM results
		WHERE id = $1;
	`

	res, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete result: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("result not found")
	}

	return nil
}
