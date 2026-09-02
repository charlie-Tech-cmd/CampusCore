package repository

import (
	"campuscore/internal/models"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// PostgresGovernanceRepository implements the models.GovernanceRepository contract.
type PostgresGovernanceRepository struct {
	db *sql.DB
}

// NewPostgresGovernanceRepository instantiates the governance repository.
func NewPostgresGovernanceRepository(db *sql.DB) *PostgresGovernanceRepository {
	return &PostgresGovernanceRepository{db: db}
}

// GetApprovalStatus retrieves the latest workflow state for a course.
// The approvals table stores session_id, while session name and semester
// are resolved through academic_sessions.
func (r *PostgresGovernanceRepository) GetApprovalStatus(courseCode string) (*models.Approval, error) {
	query := `
		SELECT
			a.id,
			a.course_code,
			s.session_name,
			s.semester::text,
			a.current_state,
			COALESCE(a.action_by, ''),
			COALESCE(a.remarks, ''),
			a.updated_at
		FROM approvals a
		JOIN academic_sessions s
			ON s.id = a.session_id
		WHERE a.course_code = $1
		ORDER BY a.updated_at DESC
		LIMIT 1;
	`

	row := r.db.QueryRow(query, courseCode)

	var approval models.Approval

	err := row.Scan(
		&approval.ID,
		&approval.CourseCode,
		&approval.Session,
		&approval.Semester,
		&approval.CurrentState,
		&approval.ActionBy,
		&approval.Remarks,
		&approval.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No workflow record exists yet. Return a virtual initial
			// submitted state using the currently active academic session.
			var session string
			var semester string

			sessionQuery := `
				SELECT session_name, semester::text
				FROM academic_sessions
				WHERE is_active = true
				ORDER BY id DESC
				LIMIT 1;
			`

			sessionErr := r.db.QueryRow(sessionQuery).Scan(
				&session,
				&semester,
			)

			if sessionErr != nil {
				if errors.Is(sessionErr, sql.ErrNoRows) {
					return nil, errors.New("no active academic session configured")
				}

				return nil, fmt.Errorf(
					"failed to determine active academic session: %w",
					sessionErr,
				)
			}

			return &models.Approval{
				CourseCode:   courseCode,
				Session:      session,
				Semester:     semester,
				CurrentState: models.StatusSubmitted,
				Remarks:      "Initial system submission baseline.",
				UpdatedAt:    time.Now(),
			}, nil
		}

		return nil, fmt.Errorf(
			"database query failure matching approval status: %w",
			err,
		)
	}

	return &approval, nil
}

// UpdateApprovalState updates the workflow state using the session associated
// with the existing approval. If no approval exists yet, the active academic
// session is used.
func (r *PostgresGovernanceRepository) UpdateApprovalState(
	courseCode string,
	newState models.ResultStatus,
	staffID string,
	remarks string,
) error {
	query := `
		WITH target_session AS (
			SELECT session_id
			FROM approvals
			WHERE course_code = $1
			ORDER BY updated_at DESC
			LIMIT 1
		),
		active_session AS (
			SELECT id AS session_id
			FROM academic_sessions
			WHERE is_active = true
			ORDER BY id DESC
			LIMIT 1
		),
		resolved_session AS (
			SELECT session_id FROM target_session
			UNION ALL
			SELECT session_id
			FROM active_session
			WHERE NOT EXISTS (
				SELECT 1 FROM target_session
			)
			LIMIT 1
		)
		INSERT INTO approvals (
			course_code,
			session_id,
			current_state,
			action_by,
			remarks,
			updated_at
		)
		SELECT
			$1,
			session_id,
			$2,
			$3,
			$4,
			CURRENT_TIMESTAMP
		FROM resolved_session
		ON CONFLICT (course_code, session_id)
		DO UPDATE SET
			current_state = EXCLUDED.current_state,
			action_by = EXCLUDED.action_by,
			remarks = EXCLUDED.remarks,
			updated_at = CURRENT_TIMESTAMP;
	`

	result, err := r.db.Exec(
		query,
		courseCode,
		newState,
		staffID,
		remarks,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to commit approval state update to postgres layer: %w",
			err,
		)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(
			"failed to verify approval state update: %w",
			err,
		)
	}

	if rowsAffected == 0 {
		return errors.New("unable to resolve an active academic session for approval state update")
	}

	return nil
}
