package api

import (
	"context"
	"errors"
	"fmt"

	"campuscore/internal/models"
)

// ClearanceService handles student clearance operations.
type ClearanceService struct {
	repo models.FinancialRepository
}

// NewClearanceService creates a new ClearanceService.
func NewClearanceService(r models.FinancialRepository) *ClearanceService {
	return &ClearanceService{
		repo: r,
	}
}

// GetChecklistStatus returns a student's clearance checklist.
func (s *ClearanceService) GetChecklistStatus(ctx context.Context, studentID string) ([]models.StudentClearance, error) {
	if studentID == "" {
		return nil, errors.New("student ID is required")
	}

	return s.repo.GetStudentClearanceStatus(ctx, studentID)
}

// ProcessOfficeSignOff updates the clearance status for an office.
func (s *ClearanceService) ProcessOfficeSignOff(
	ctx context.Context,
	studentID string,
	officeID int,
	status string,
	staffID string,
) error {

	if studentID == "" {
		return errors.New("student ID is required")
	}

	if staffID == "" {
		return errors.New("staff ID is required")
	}

	var clearanceStatus models.ClearanceStatus

	switch models.ClearanceStatus(status) {
	case models.ClearancePending,
		models.ClearanceSubmitted,
		models.ClearanceCleared:
		clearanceStatus = models.ClearanceStatus(status)

	default:
		return fmt.Errorf("invalid clearance status: %s", status)
	}

	return s.repo.UpdateClearanceStatus(
		ctx,
		studentID,
		officeID,
		clearanceStatus,
		staffID,
	)
}