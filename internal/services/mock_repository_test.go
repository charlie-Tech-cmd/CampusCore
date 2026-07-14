package services

import (
	"context"

	"campuscore/internal/models"
)

// mockFinancialRepository implements models.FinancialRepository for unit tests.
type mockFinancialRepository struct {
	// Function fields allow each test to customize behavior.
	getFeeStructureFn          func(ctx context.Context, departmentID, level int, feeType, session string) (*models.FeeStructure, error)
	recordPaymentFn            func(ctx context.Context, payment *models.FeePayment) error
	checkPaymentExistsFn       func(ctx context.Context, gatewayRef string) (bool, error)
	getStudentClearanceStatusFn func(ctx context.Context, studentID string) ([]models.StudentClearance, error)
	updateClearanceStatusFn    func(ctx context.Context, studentID string, officeID int, status models.ClearanceStatus, staffID string) error
	createTicketFn             func(ctx context.Context, ticket *models.SupportTicket) error
}

func (m *mockFinancialRepository) GetFeeStructure(
	ctx context.Context,
	departmentID, level int,
	feeType, session string,
) (*models.FeeStructure, error) {

	if m.getFeeStructureFn != nil {
		return m.getFeeStructureFn(ctx, departmentID, level, feeType, session)
	}

	return &models.FeeStructure{}, nil
}

func (m *mockFinancialRepository) RecordPayment(
	ctx context.Context,
	payment *models.FeePayment,
) error {

	if m.recordPaymentFn != nil {
		return m.recordPaymentFn(ctx, payment)
	}

	return nil
}

func (m *mockFinancialRepository) CheckPaymentExists(
	ctx context.Context,
	gatewayRef string,
) (bool, error) {

	if m.checkPaymentExistsFn != nil {
		return m.checkPaymentExistsFn(ctx, gatewayRef)
	}

	return false, nil
}

func (m *mockFinancialRepository) GetStudentClearanceStatus(
	ctx context.Context,
	studentID string,
) ([]models.StudentClearance, error) {

	if m.getStudentClearanceStatusFn != nil {
		return m.getStudentClearanceStatusFn(ctx, studentID)
	}

	return []models.StudentClearance{}, nil
}

func (m *mockFinancialRepository) UpdateClearanceStatus(
	ctx context.Context,
	studentID string,
	officeID int,
	status models.ClearanceStatus,
	staffID string,
) error {

	if m.updateClearanceStatusFn != nil {
		return m.updateClearanceStatusFn(ctx, studentID, officeID, status, staffID)
	}

	return nil
}

func (m *mockFinancialRepository) CreateTicket(
	ctx context.Context,
	ticket *models.SupportTicket,
) error {

	if m.createTicketFn != nil {
		return m.createTicketFn(ctx, ticket)
	}

	return nil
}