package api

import (
	"context"

	"campuscore/internal/models"
)

type mockFinancialRepository struct {
	getFeeStructureFunc           func(context.Context, int, int, string, string) (*models.FeeStructure, error)
	recordPaymentFunc             func(context.Context, *models.FeePayment) error
	checkPaymentExistsFunc        func(context.Context, string) (bool, error)
	getStudentClearanceStatusFunc func(context.Context, string) ([]models.StudentClearance, error)
	updateClearanceStatusFunc     func(context.Context, string, int, models.ClearanceStatus, string) error
	createTicketFunc              func(context.Context, *models.SupportTicket) error
}

func (m *mockFinancialRepository) GetFeeStructure(
	ctx context.Context,
	departmentID, level int,
	feeType, session string,
) (*models.FeeStructure, error) {
	if m.getFeeStructureFunc != nil {
		return m.getFeeStructureFunc(ctx, departmentID, level, feeType, session)
	}
	return nil, nil
}

func (m *mockFinancialRepository) RecordPayment(
	ctx context.Context,
	payment *models.FeePayment,
) error {
	if m.recordPaymentFunc != nil {
		return m.recordPaymentFunc(ctx, payment)
	}
	return nil
}

func (m *mockFinancialRepository) CheckPaymentExists(
	ctx context.Context,
	gatewayRef string,
) (bool, error) {
	if m.checkPaymentExistsFunc != nil {
		return m.checkPaymentExistsFunc(ctx, gatewayRef)
	}
	return false, nil
}

func (m *mockFinancialRepository) GetStudentClearanceStatus(
	ctx context.Context,
	studentID string,
) ([]models.StudentClearance, error) {
	if m.getStudentClearanceStatusFunc != nil {
		return m.getStudentClearanceStatusFunc(ctx, studentID)
	}
	return nil, nil
}

func (m *mockFinancialRepository) UpdateClearanceStatus(
	ctx context.Context,
	studentID string,
	officeID int,
	status models.ClearanceStatus,
	staffID string,
) error {
	if m.updateClearanceStatusFunc != nil {
		return m.updateClearanceStatusFunc(
			ctx,
			studentID,
			officeID,
			status,
			staffID,
		)
	}
	return nil
}

func (m *mockFinancialRepository) CreateTicket(
	ctx context.Context,
	ticket *models.SupportTicket,
) error {
	if m.createTicketFunc != nil {
		return m.createTicketFunc(ctx, ticket)
	}
	return nil
}
