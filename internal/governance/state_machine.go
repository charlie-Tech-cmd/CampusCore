package governance

import (
	"campuscore/internal/models"
	"errors"
	"fmt"
	"strings"
)

// Engine handles the logical validations for our institutional workflow states.
type Engine struct {
	repo models.GovernanceRepository
}

// NewEngine instantiates our governance workflow controller.
func NewEngine(r models.GovernanceRepository) *Engine {
	return &Engine{repo: r}
}

// ProcessApprovalAdvance moves a result batch forward through the institutional hierarchy.
func (e *Engine) ProcessApprovalAdvance(
	courseCode string,
	currentActorRole models.UserRole,
	staffID string,
) error {
	approval, err := e.repo.GetApprovalStatus(courseCode)
	if err != nil {
		return fmt.Errorf("failed to check current workflow state: %w", err)
	}

	var nextState models.ResultStatus

	switch approval.CurrentState {
	case models.StatusSubmitted:
		// The current role model has no explicit HOD role.
		// Admin is therefore the only role authorized to perform
		// the HOD approval tier until a dedicated HOD role exists.
		if currentActorRole != models.RoleAdmin {
			return errors.New(
				"governance conflict: only an authorized administrator can approve a primary submission",
			)
		}

		nextState = models.StatusHODApproved

	case models.StatusHODApproved:
		// Lecturer currently represents the next academic approval tier.
		// Admin retains an institutional override capability.
		if currentActorRole != models.RoleAdmin &&
			currentActorRole != models.RoleLecturer {
			return errors.New(
				"governance conflict: only an authorized lecturer or administrator can advance an HOD-approved result",
			)
		}

		nextState = models.StatusDeanApproved

	case models.StatusDeanApproved:
		// There is currently no dedicated Senate role in the user model.
		if currentActorRole != models.RoleAdmin {
			return errors.New(
				"governance conflict: only an authorized administrator can grant institutional final approval",
			)
		}

		nextState = models.StatusSenateApproved

	case models.StatusSenateApproved:
		if currentActorRole != models.RoleAdmin {
			return errors.New(
				"governance conflict: only an authorized administrator can finalize an academic record",
			)
		}

		nextState = models.StatusFinalized

	case models.StatusFinalized:
		return errors.New(
			"invalid operation: this academic record is finalized and locked against changes",
		)

	default:
		return fmt.Errorf(
			"invalid governance state: %q",
			approval.CurrentState,
		)
	}

	return e.repo.UpdateApprovalState(
		courseCode,
		nextState,
		staffID,
		"Forwarded to next governance tier.",
	)
}

// ProcessApprovalRejection processes rollback steps while enforcing
// authorization for the current workflow tier.
func (e *Engine) ProcessApprovalRejection(
	courseCode string,
	currentActorRole models.UserRole,
	staffID string,
	remarks string,
) error {
	cleanRemarks := strings.TrimSpace(remarks)

	if len(cleanRemarks) < 10 {
		return errors.New(
			"validation error: you must provide an explicit reason string (minimum 10 characters) to reject a result batch",
		)
	}

	approval, err := e.repo.GetApprovalStatus(courseCode)
	if err != nil {
		return fmt.Errorf("failed to check current workflow state: %w", err)
	}

	var targetBackwardState models.ResultStatus

	switch approval.CurrentState {
	case models.StatusSubmitted:
		return errors.New(
			"invalid operation: cannot reject a batch that is currently at initial submission level",
		)

	case models.StatusHODApproved:
		// The next academic tier is responsible for rejecting an
		// HOD-approved submission back to submitted.
		if currentActorRole != models.RoleAdmin &&
			currentActorRole != models.RoleLecturer {
			return errors.New(
				"governance conflict: only an authorized lecturer or administrator can reject an HOD-approved result",
			)
		}

		targetBackwardState = models.StatusSubmitted

	case models.StatusDeanApproved:
		// Admin is the only currently defined institutional authority
		// capable of handling this higher-level rollback.
		if currentActorRole != models.RoleAdmin {
			return errors.New(
				"governance conflict: only an authorized administrator can reject a dean-approved result",
			)
		}

		targetBackwardState = models.StatusHODApproved

	case models.StatusSenateApproved:
		if currentActorRole != models.RoleAdmin {
			return errors.New(
				"governance conflict: only an authorized administrator can reject a senate-approved result",
			)
		}

		targetBackwardState = models.StatusDeanApproved

	case models.StatusFinalized:
		return errors.New(
			"critical security failure: finalized transcripts cannot be rejected via standard endpoints",
		)

	default:
		return fmt.Errorf(
			"invalid governance state: %q",
			approval.CurrentState,
		)
	}

	return e.repo.UpdateApprovalState(
		courseCode,
		targetBackwardState,
		staffID,
		cleanRemarks,
	)
}
