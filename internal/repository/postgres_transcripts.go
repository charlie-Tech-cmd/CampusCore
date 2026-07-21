package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"campuscore/internal/models"
)

// PostgresTranscriptRepository manages transcript persistence.
type PostgresTranscriptRepository struct {
	db *sql.DB
}

// NewPostgresTranscriptRepository creates a transcript repository.
func NewPostgresTranscriptRepository(db *sql.DB) *PostgresTranscriptRepository {
	return &PostgresTranscriptRepository{
		db: db,
	}
}

// Create inserts a transcript.
func (r *PostgresTranscriptRepository) Create(transcript *models.Transcript) error {
	query := `
		INSERT INTO transcripts
		(
			id,
			student_id,
			session,
			semester,
			cgpa,
			total_credits,
			remarks,
			generated_at,
			generated_by
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9);
	`

	_, err := r.db.Exec(
		query,
		transcript.ID,
		transcript.StudentID,
		transcript.Session,
		transcript.Semester,
		transcript.CGPA,
		transcript.TotalCredits,
		transcript.Remarks,
		transcript.GeneratedAt,
		transcript.GeneratedBy,
	)

	if err != nil {
		return fmt.Errorf("failed to create transcript: %w", err)
	}

	return nil
}

// FindByID retrieves a transcript by ID.
func (r *PostgresTranscriptRepository) FindByID(id int) (*models.Transcript, error) {
	query := `
		SELECT
			id,
			student_id,
			session,
			semester,
			cgpa,
			total_credits,
			remarks,
			generated_at,
			generated_by
		FROM transcripts
		WHERE id=$1
		LIMIT 1;
	`

	var transcript models.Transcript

	err := r.db.QueryRow(query, id).Scan(
		&transcript.ID,
		&transcript.StudentID,
		&transcript.Session,
		&transcript.Semester,
		&transcript.CGPA,
		&transcript.TotalCredits,
		&transcript.Remarks,
		&transcript.GeneratedAt,
		&transcript.GeneratedBy,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("transcript not found")
		}
		return nil, fmt.Errorf("failed to retrieve transcript: %w", err)
	}

	return &transcript, nil
}

// FindByStudent retrieves transcripts belonging to a student.
func (r *PostgresTranscriptRepository) FindByStudent(studentID string) ([]models.Transcript, error) {
	query := `
		SELECT
			id,
			student_id,
			session,
			semester,
			cgpa,
			total_credits,
			remarks,
			generated_at,
			generated_by
		FROM transcripts
		WHERE student_id=$1
		ORDER BY generated_at DESC;
	`

	rows, err := r.db.Query(query, studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve transcripts: %w", err)
	}
	defer rows.Close()

	var transcripts []models.Transcript

	for rows.Next() {
		var transcript models.Transcript

		if err := rows.Scan(
			&transcript.ID,
			&transcript.StudentID,
			&transcript.Session,
			&transcript.Semester,
			&transcript.CGPA,
			&transcript.TotalCredits,
			&transcript.Remarks,
			&transcript.GeneratedAt,
			&transcript.GeneratedBy,
		); err != nil {
			return nil, err
		}

		transcripts = append(transcripts, transcript)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transcripts, nil
}

// List retrieves all transcripts.
func (r *PostgresTranscriptRepository) List() ([]models.Transcript, error) {
	query := `
		SELECT
			id,
			student_id,
			session,
			semester,
			cgpa,
			total_credits,
			remarks,
			generated_at,
			generated_by
		FROM transcripts
		ORDER BY generated_at DESC;
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve transcripts: %w", err)
	}
	defer rows.Close()

	var transcripts []models.Transcript

	for rows.Next() {
		var transcript models.Transcript

		if err := rows.Scan(
			&transcript.ID,
			&transcript.StudentID,
			&transcript.Session,
			&transcript.Semester,
			&transcript.CGPA,
			&transcript.TotalCredits,
			&transcript.Remarks,
			&transcript.GeneratedAt,
			&transcript.GeneratedBy,
		); err != nil {
			return nil, err
		}

		transcripts = append(transcripts, transcript)
	}

	return transcripts, rows.Err()
}

// Update updates a transcript.
func (r *PostgresTranscriptRepository) Update(transcript *models.Transcript) error {
	query := `
		UPDATE transcripts
		SET
			cgpa=$2,
			total_credits=$3,
			remarks=$4,
			generated_at=$5,
			generated_by=$6
		WHERE id=$1;
	`

	result, err := r.db.Exec(
		query,
		transcript.ID,
		transcript.CGPA,
		transcript.TotalCredits,
		transcript.Remarks,
		transcript.GeneratedAt,
		transcript.GeneratedBy,
	)

	if err != nil {
		return fmt.Errorf("failed to update transcript: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("transcript not found")
	}

	return nil
}

// Delete removes a transcript.
func (r *PostgresTranscriptRepository) Delete(id int) error {
	result, err := r.db.Exec(
		`DELETE FROM transcripts WHERE id=$1`,
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to delete transcript: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("transcript not found")
	}

	return nil
}
