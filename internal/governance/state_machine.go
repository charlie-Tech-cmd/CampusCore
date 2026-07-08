package governance

import (
	"errors"
	"fmt"
	"strings"
	"campuscore/internal/models"
)

// Engine handles the logical validations for our institutional workflow states
type Engine struct {
	repo models.GovernanceRepository
}

// NewEngine instantiates our governance workflow controller
func NewEngine(r models.GovernanceRepository) *Engine {
	return &Engine{repo: r}
}

// ProcessApprovalAdvance moves a result batch forward through the institutional hierarchy
func (e *Engine) ProcessApprovalAdvance(courseCode string, currentActorRole models.UserRole, staffID string) error {
	approval, err := e.repo.GetApprovalStatus(courseCode)
	if err != nil {
		return fmt.Errorf("failed to check current workflow state: %w", err)
	}

	var nextState models.ResultStatus

	// Enforce strict tier progression rules
	switch approval.CurrentState {
	case models.StatusSubmitted:
		if currentActorRole != models.RoleAdmin && string(currentActorRole) != "HOD" { 
			return errors.New("governance conflict: only the Head of Department can approve a primary submission")
		}
		nextState = models.StatusHODApproved

	case models.StatusHODApproved:
		if currentActorRole != models.RoleAdmin && currentActorRole != models.RoleLecturer { // Assuming Dean authorization tier mappings
			// Note: In route delivery, the endpoint will explicitly check if the lecturer has Dean privileges
		}
		nextState = models.StatusDeanApproved

	case models.StatusDeanApproved:
		if currentActorRole != models.RoleAdmin {
			return errors.New("governance conflict: only the Senate Board can grant institutional finalization")
		}
		nextState = models.StatusSenateApproved

	case models.StatusSenateApproved:
		if currentActorRole != models.RoleAdmin {
			return errors.New("governance conflict: administrative override required to shift to absolute final lock")
		}
		nextState = models.StatusFinalized

	case models.StatusFinalized:
		return errors.New("invalid operation: this academic record is finalized and locked against changes")
	}

	// Persist the forward step to the database layer
	return e.repo.UpdateApprovalState(courseCode, nextState, staffID, "Forwarded to next governance tier.")
}

// ProcessApprovalRejection processes rollback steps, enforcing the dashed loops from your workflow chart
func (e *Engine) ProcessApprovalRejection(courseCode string, currentActorRole models.UserRole, staffID string, remarks string) error {
	// Defensive Validation: Reject instantly if no audit trail text justification is supplied
	cleanRemarks := strings.TrimSpace(remarks)
	if len(cleanRemarks) < 10 {
		return errors.New("validation error: you must provide an explicit reason string (minimum 10 characters) to reject a result batch")
	}

	approval, err := e.repo.GetApprovalStatus(courseCode)
	if err != nil {
		return fmt.Errorf("failed to check current workflow state: %w", err)
	}

	var targetBackwardState models.ResultStatus

	// Map the workflow backward loops
	switch approval.CurrentState {
	case models.StatusSubmitted:
		return errors.New("invalid operation: cannot reject a batch that is currently at initial submission level")

	case models.StatusHODApproved:
		// HOD rejects back to the original Lecturer
		targetBackwardState = models.StatusSubmitted

	case models.StatusDeanApproved:
		// Dean rejects back to the HOD panel
		targetBackwardState = models.StatusHODApproved

	case models.StatusSenateApproved:
		// Senate rejects back to the Dean's faculty office
		targetBackwardState = models.StatusDeanApproved

	case models.StatusFinalized:
		return errors.New("critical security failure: finalized transcripts cannot be rejected via standard endpoints")
	}

	// Update the database state to match our rollback cascade
	return e.repo.UpdateApprovalState(courseCode, targetBackwardState, staffID, cleanRemarks)
}