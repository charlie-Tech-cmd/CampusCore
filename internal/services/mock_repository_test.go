package services

import (
	"context"
	"errors"
	"testing"

	"campuscore/internal/models"
)

// mockFinancialRepository implements models.FinancialRepository for unit tests.
type mockFinancialRepository struct {
	// Function fields allow each test to customize behavior.
	getFeeStructureFn           func(ctx context.Context, departmentID, level int, feeType, session string) (*models.FeeStructure, error)
	recordPaymentFn             func(ctx context.Context, payment *models.FeePayment) error
	checkPaymentExistsFn        func(ctx context.Context, gatewayRef string) (bool, error)
	getStudentClearanceStatusFn func(ctx context.Context, studentID string) ([]models.StudentClearance, error)
	updateClearanceStatusFn     func(ctx context.Context, studentID string, officeID int, status models.ClearanceStatus, staffID string) error
	createTicketFn              func(ctx context.Context, ticket *models.SupportTicket) error
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

func TestGetStudentClearance_EmptyStudentID(t *testing.T) {
	repo := &mockFinancialRepository{}
	service := NewClearanceService(repo)

	_, err := service.GetStudentClearance(context.Background(), "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "student ID is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetStudentClearance_RepositoryError(t *testing.T) {
	expectedErr := errors.New("database unavailable")

	repo := &mockFinancialRepository{
		getStudentClearanceStatusFn: func(
			ctx context.Context,
			studentID string,
		) ([]models.StudentClearance, error) {
			return nil, expectedErr
		},
	}

	service := NewClearanceService(repo)

	_, err := service.GetStudentClearance(
		context.Background(),
		"STU001",
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}

func TestUpdateClearance_Success(t *testing.T) {
	repo := &mockFinancialRepository{
		updateClearanceStatusFn: func(
			ctx context.Context,
			studentID string,
			officeID int,
			status models.ClearanceStatus,
			staffID string,
		) error {

			if studentID != "STU001" {
				t.Fatalf("expected STU001, got %q", studentID)
			}

			if officeID != 1 {
				t.Fatalf("expected officeID 1, got %d", officeID)
			}

			if status != models.ClearanceCleared {
				t.Fatalf("unexpected status %q", status)
			}

			if staffID != "STAFF001" {
				t.Fatalf("expected STAFF001, got %q", staffID)
			}

			return nil
		},
	}

	service := NewClearanceService(repo)

	err := service.UpdateClearance(
		context.Background(),
		"STU001",
		1,
		models.ClearanceCleared,
		"STAFF001",
	)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestUpdateClearance_ValidationErrors(t *testing.T) {
	service := NewClearanceService(&mockFinancialRepository{})

	tests := []struct {
		name      string
		studentID string
		officeID  int
		status    models.ClearanceStatus
		wantErr   string
	}{
		{
			name:      "missing student ID",
			studentID: "",
			officeID:  1,
			status:    models.ClearancePending,
			wantErr:   "student ID is required",
		},
		{
			name:      "invalid office ID",
			studentID: "STU001",
			officeID:  0,
			status:    models.ClearancePending,
			wantErr:   "invalid office ID",
		},
		{
			name:      "invalid status",
			studentID: "STU001",
			officeID:  1,
			status:    models.ClearanceStatus("invalid"),
			wantErr:   "invalid clearance status",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := service.UpdateClearance(
				context.Background(),
				tc.studentID,
				tc.officeID,
				tc.status,
				"STAFF001",
			)

			if err == nil {
				t.Fatal("expected error")
			}

			if err.Error() != tc.wantErr {
				t.Fatalf("expected %q got %q", tc.wantErr, err.Error())
			}
		})
	}
}

func TestUpdateClearance_RepositoryError(t *testing.T) {
	expectedErr := errors.New("database error")

	repo := &mockFinancialRepository{
		updateClearanceStatusFn: func(
			ctx context.Context,
			studentID string,
			officeID int,
			status models.ClearanceStatus,
			staffID string,
		) error {
			return expectedErr
		},
	}

	service := NewClearanceService(repo)

	err := service.UpdateClearance(
		context.Background(),
		"STU001",
		1,
		models.ClearancePending,
		"STAFF001",
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}

func TestIsStudentCleared_AllCleared(t *testing.T) {
	repo := &mockFinancialRepository{
		getStudentClearanceStatusFn: func(
			ctx context.Context,
			studentID string,
		) ([]models.StudentClearance, error) {

			return []models.StudentClearance{
				{OfficeID: 1, Status: models.ClearanceCleared},
				{OfficeID: 2, Status: models.ClearanceCleared},
				{OfficeID: 3, Status: models.ClearanceCleared},
			}, nil
		},
	}

	service := NewClearanceService(repo)

	cleared, err := service.IsStudentCleared(
		context.Background(),
		"STU001",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !cleared {
		t.Fatal("expected student to be cleared")
	}
}

func TestIsStudentCleared_NotCleared(t *testing.T) {
	repo := &mockFinancialRepository{
		getStudentClearanceStatusFn: func(
			ctx context.Context,
			studentID string,
		) ([]models.StudentClearance, error) {

			return []models.StudentClearance{
				{OfficeID: 1, Status: models.ClearanceCleared},
				{OfficeID: 2, Status: models.ClearancePending},
			}, nil
		},
	}

	service := NewClearanceService(repo)

	cleared, err := service.IsStudentCleared(
		context.Background(),
		"STU001",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cleared {
		t.Fatal("expected student NOT to be cleared")
	}
}

func TestIsStudentCleared_NoRecords(t *testing.T) {
	repo := &mockFinancialRepository{
		getStudentClearanceStatusFn: func(
			ctx context.Context,
			studentID string,
		) ([]models.StudentClearance, error) {

			return []models.StudentClearance{}, nil
		},
	}

	service := NewClearanceService(repo)

	cleared, err := service.IsStudentCleared(
		context.Background(),
		"STU001",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cleared {
		t.Fatal("expected false")
	}
}

func TestIsStudentCleared_RepositoryError(t *testing.T) {
	expectedErr := errors.New("database failure")

	repo := &mockFinancialRepository{
		getStudentClearanceStatusFn: func(
			ctx context.Context,
			studentID string,
		) ([]models.StudentClearance, error) {

			return nil, expectedErr
		},
	}

	service := NewClearanceService(repo)

	_, err := service.IsStudentCleared(
		context.Background(),
		"STU001",
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}

func TestIsStudentCleared_EmptyStudentID(t *testing.T) {
	service := NewClearanceService(&mockFinancialRepository{})

	_, err := service.IsStudentCleared(
		context.Background(),
		"",
	)

	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "student ID is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}
