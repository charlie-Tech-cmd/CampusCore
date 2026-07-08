package services

import (
	"context"
	"errors"

	"campuscore/internal/models"
)

// ClearanceService handles student clearance.
type ClearanceService struct {
	repo models.FinancialRepository
}

// NewClearanceService creates a ClearanceService.
func NewClearanceService(repo models.FinancialRepository) *ClearanceService {
	return &ClearanceService{
		repo: repo,
	}
}

// GetStudentClearance returns a student's clearance records.
func (s *ClearanceService) GetStudentClearance(
	ctx context.Context,
	studentID string,
) ([]models.StudentClearance, error) {

	if studentID == "" {
		return nil, errors.New("student ID is required")
	}

	return s.repo.GetStudentClearanceStatus(ctx, studentID)
}

// UpdateClearance updates the clearance status for an office.
func (s *ClearanceService) UpdateClearance(
	ctx context.Context,
	studentID string,
	officeID int,
	status models.ClearanceStatus,
	staffID string,
) error {

	if studentID == "" {
		return errors.New("student ID is required")
	}

	if officeID <= 0 {
		return errors.New("invalid office ID")
	}

	switch status {
	case models.ClearancePending,
		models.ClearanceSubmitted,
		models.ClearanceCleared:
		// Valid status.
	default:
		return errors.New("invalid clearance status")
	}

	return s.repo.UpdateClearanceStatus(
		ctx,
		studentID,
		officeID,
		status,
		staffID,
	)
}

// IsStudentCleared reports whether every clearance has been completed.
func (s *ClearanceService) IsStudentCleared(
	ctx context.Context,
	studentID string,
) (bool, error) {

	if studentID == "" {
		return false, errors.New("student ID is required")
	}

	clearances, err := s.repo.GetStudentClearanceStatus(ctx, studentID)
	if err != nil {
		return false, err
	}

	if len(clearances) == 0 {
		return false, nil
	}

	for _, clearance := range clearances {
		if clearance.Status != models.ClearanceCleared {
			return false, nil
		}
	}

	return true, nil
}