package repository

import (
	"campuscore/internal/models"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// PostgresGovernanceRepository implements the models.GovernanceRepository contract
type PostgresGovernanceRepository struct {
	db *sql.DB
}

// NewPostgresGovernanceRepository instantiates our data access controller with a DB pool pointer
func NewPostgresGovernanceRepository(db *sql.DB) *PostgresGovernanceRepository {
	return &PostgresGovernanceRepository{db: db}
}

// GetApprovalStatus retrieves the active state workflow item for a target course module
func (r *PostgresGovernanceRepository) GetApprovalStatus(courseCode string) (*models.Approval, error) {
	query := `
		SELECT id, course_code, session, semester, current_state, action_by, remarks, updated_at 
		FROM approvals 
		WHERE course_code = $1 
		LIMIT 1;`

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
			// Defensive Strategy: If no workflow record exists yet, return an initial virtual "submitted" state container
			return &models.Approval{
				CourseCode:   courseCode,
				Session:      "2026/2027", // Derived from system state configs
				Semester:     "First",
				CurrentState: models.StatusSubmitted,
				Remarks:      "Initial system submission baseline.",
				UpdatedAt:    time.Now(),
			}, nil
		}
		return nil, fmt.Errorf("database query failure matching approval status: %w", err)
	}

	return &approval, nil
}

// UpdateApprovalState executes an atomic update transaction to safely alter the validation state
func (r *PostgresGovernanceRepository) UpdateApprovalState(courseCode string, newState models.ResultStatus, staffID string, remarks string) error {
	// We use an UPSERT (INSERT ... ON CONFLICT) statement to modify the row if it exists,
	// or initialize it safely if this is the first approval step.
	query := `
		INSERT INTO approvals (course_code, session, semester, current_state, action_by, remarks, updated_at)
		VALUES ($1, '2026/2027', 'First', $2, $3, $4, CURRENT_TIMESTAMP)
		ON CONFLICT (course_code, session, semester) 
		DO UPDATE SET 
			current_state = EXCLUDED.current_state,
			action_by = EXCLUDED.action_by,
			remarks = EXCLUDED.remarks,
			updated_at = CURRENT_TIMESTAMP;`

	_, err := r.db.Exec(query, courseCode, newState, staffID, remarks)
	if err != nil {
		return fmt.Errorf("failed to commit approval state update to postgres layer: %w", err)
	}

	return nil
}
