package services

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"campuscore/internal/models"
)

func TestProcessIncomingWebhook_Success(t *testing.T) {
	repo := &mockFinancialRepository{
		checkPaymentExistsFn: func(ctx context.Context, ref string) (bool, error) {
			if ref != "REF123" {
				t.Fatalf("expected REF123, got %q", ref)
			}
			return false, nil
		},
		recordPaymentFn: func(ctx context.Context, payment *models.FeePayment) error {
			if payment.StudentID != "STU001" {
				t.Fatalf("unexpected student ID")
			}

			if payment.Status != "successful" {
				t.Fatalf("expected successful status")
			}

			return nil
		},
	}

	service := NewFinancialService(repo, nil)

	err := service.ProcessIncomingWebhook(
		context.Background(),
		"STU001",
		"REF123",
		50000,
		"school_fees",
		"2025/2026",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}


func TestNewFinancialService(t *testing.T) {
	repo := &mockFinancialRepository{}

	service := NewFinancialService(repo, nil)

	if service == nil {
		t.Fatal("expected FinancialService, got nil")
	}

	if service.repo != repo {
		t.Fatal("repository was not assigned correctly")
	}

	if service.db != nil {
		t.Fatal("expected db to be nil")
	}
}

func TestProcessIncomingWebhook_DuplicatePayment(t *testing.T) {
	repo := &mockFinancialRepository{
		checkPaymentExistsFn: func(ctx context.Context, ref string) (bool, error) {
			return true, nil
		},
	}

	service := NewFinancialService(repo, nil)

	err := service.ProcessIncomingWebhook(
		context.Background(),
		"STU001",
		"REF123",
		50000,
		"school_fees",
		"2025/2026",
	)

	if err == nil {
		t.Fatal("expected duplicate payment error")
	}
}

func TestProcessIncomingWebhook_CheckPaymentExistsError(t *testing.T) {
	repo := &mockFinancialRepository{
		checkPaymentExistsFn: func(ctx context.Context, ref string) (bool, error) {
			return false, errors.New("database unavailable")
		},
	}

	service := NewFinancialService(repo, nil)

	err := service.ProcessIncomingWebhook(
		context.Background(),
		"STU001",
		"REF123",
		50000,
		"school_fees",
		"2025/2026",
	)

	if err == nil {
		t.Fatal("expected repository error")
	}
}

func TestProcessIncomingWebhook_RecordPaymentError(t *testing.T) {
	repo := &mockFinancialRepository{
		checkPaymentExistsFn: func(ctx context.Context, ref string) (bool, error) {
			return false, nil
		},
		recordPaymentFn: func(ctx context.Context, payment *models.FeePayment) error {
			return errors.New("failed to save payment")
		},
	}

	service := NewFinancialService(repo, nil)

	err := service.ProcessIncomingWebhook(
		context.Background(),
		"STU001",
		"REF123",
		50000,
		"school_fees",
		"2025/2026",
	)

	if err == nil {
		t.Fatal("expected record payment error")
	}
}

func TestEvaluateGraduationEligibility_Success(t *testing.T) {
	repo := &mockFinancialRepository{
		getStudentClearanceStatusFn: func(ctx context.Context, studentID string) ([]models.StudentClearance, error) {
			return []models.StudentClearance{
				{Status: models.ClearanceCleared},
				{Status: models.ClearanceCleared},
				{Status: models.ClearanceCleared},
			}, nil
		},
	}

	service := NewFinancialService(repo, nil)

	ok, err := service.EvaluateGraduationEligibility(
		context.Background(),
		"STU001",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !ok {
		t.Fatal("expected student to be eligible for graduation")
	}
}

func TestEvaluateGraduationEligibility_NotCleared(t *testing.T) {
	repo := &mockFinancialRepository{
		getStudentClearanceStatusFn: func(ctx context.Context, studentID string) ([]models.StudentClearance, error) {
			return []models.StudentClearance{
				{Status: models.ClearanceCleared},
				{Status: models.ClearancePending},
			}, nil
		},
	}

	service := NewFinancialService(repo, nil)

	ok, err := service.EvaluateGraduationEligibility(
		context.Background(),
		"STU001",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ok {
		t.Fatal("expected student not to be eligible")
	}
}

func TestEvaluateGraduationEligibility_NoRecords(t *testing.T) {
	repo := &mockFinancialRepository{
		getStudentClearanceStatusFn: func(ctx context.Context, studentID string) ([]models.StudentClearance, error) {
			return []models.StudentClearance{}, nil
		},
	}

	service := NewFinancialService(repo, nil)

	ok, err := service.EvaluateGraduationEligibility(
		context.Background(),
		"STU001",
	)

	if err == nil {
		t.Fatal("expected error")
	}

	if ok {
		t.Fatal("expected false")
	}
}

func TestEvaluateGraduationEligibility_RepositoryError(t *testing.T) {
	repo := &mockFinancialRepository{
		getStudentClearanceStatusFn: func(ctx context.Context, studentID string) ([]models.StudentClearance, error) {
			return nil, errors.New("database unavailable")
		},
	}

	service := NewFinancialService(repo, nil)

	ok, err := service.EvaluateGraduationEligibility(
		context.Background(),
		"STU001",
	)

	if err == nil {
		t.Fatal("expected repository error")
	}

	if ok {
		t.Fatal("expected false")
	}
}

func TestVerifyTuitionClearance_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := &mockFinancialRepository{
		getFeeStructureFn: func(
			ctx context.Context,
			departmentID, level int,
			feeType, session string,
		) (*models.FeeStructure, error) {

			if departmentID != 1 {
				t.Fatalf("expected departmentID 1, got %d", departmentID)
			}

			if level != 300 {
				t.Fatalf("expected level 300, got %d", level)
			}

			return &models.FeeStructure{
				AmountRequired: 50000,
			}, nil
		},
	}

	service := NewFinancialService(repo, db)

	mock.ExpectQuery("SELECT department_id, level").
		WithArgs("STU001").
		WillReturnRows(
			sqlmock.NewRows([]string{"department_id", "level"}).
				AddRow(1, 300),
		)

	mock.ExpectQuery("SELECT COALESCE").
		WithArgs("STU001", "2025/2026").
		WillReturnRows(
			sqlmock.NewRows([]string{"sum"}).
				AddRow(50000),
		)

	ok, err := service.VerifyTuitionClearance(
		context.Background(),
		"STU001",
		"2025/2026",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !ok {
		t.Fatal("expected tuition clearance")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestVerifyTuitionClearance_InsufficientPayment(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := &mockFinancialRepository{
		getFeeStructureFn: func(ctx context.Context, departmentID, level int, feeType, session string) (*models.FeeStructure, error) {
			return &models.FeeStructure{
				AmountRequired: 50000,
			}, nil
		},
	}

	service := NewFinancialService(repo, db)

	mock.ExpectQuery("SELECT department_id, level").
		WithArgs("STU001").
		WillReturnRows(
			sqlmock.NewRows([]string{"department_id", "level"}).
				AddRow(1, 300),
		)

	mock.ExpectQuery("SELECT COALESCE").
		WithArgs("STU001", "2025/2026").
		WillReturnRows(
			sqlmock.NewRows([]string{"sum"}).
				AddRow(20000),
		)

	ok, err := service.VerifyTuitionClearance(
		context.Background(),
		"STU001",
		"2025/2026",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ok {
		t.Fatal("expected tuition clearance to be false")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestVerifyTuitionClearance_StudentQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := &mockFinancialRepository{}
	service := NewFinancialService(repo, db)

	mock.ExpectQuery("SELECT department_id, level").
		WithArgs("STU001").
		WillReturnError(errors.New("database error"))

	ok, err := service.VerifyTuitionClearance(
		context.Background(),
		"STU001",
		"2025/2026",
	)

	if err == nil {
		t.Fatal("expected error")
	}

	if ok {
		t.Fatal("expected false")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestVerifyTuitionClearance_FeeStructureError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := &mockFinancialRepository{
		getFeeStructureFn: func(
			ctx context.Context,
			departmentID, level int,
			feeType, session string,
		) (*models.FeeStructure, error) {
			return nil, errors.New("fee structure not found")
		},
	}

	service := NewFinancialService(repo, db)

	mock.ExpectQuery("SELECT department_id, level").
		WithArgs("STU001").
		WillReturnRows(
			sqlmock.NewRows([]string{"department_id", "level"}).
				AddRow(1, 300),
		)

	ok, err := service.VerifyTuitionClearance(
		context.Background(),
		"STU001",
		"2025/2026",
	)

	if err == nil {
		t.Fatal("expected error")
	}

	if ok {
		t.Fatal("expected false")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestVerifyTuitionClearance_PaymentQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := &mockFinancialRepository{
		getFeeStructureFn: func(
			ctx context.Context,
			departmentID, level int,
			feeType, session string,
		) (*models.FeeStructure, error) {
			return &models.FeeStructure{
				AmountRequired: 50000,
			}, nil
		},
	}

	service := NewFinancialService(repo, db)

	mock.ExpectQuery("SELECT department_id, level").
		WithArgs("STU001").
		WillReturnRows(
			sqlmock.NewRows([]string{"department_id", "level"}).
				AddRow(1, 300),
		)

	mock.ExpectQuery("SELECT COALESCE").
		WithArgs("STU001", "2025/2026").
		WillReturnError(errors.New("payment query failed"))

	ok, err := service.VerifyTuitionClearance(
		context.Background(),
		"STU001",
		"2025/2026",
	)

	if err == nil {
		t.Fatal("expected error")
	}

	if ok {
		t.Fatal("expected false")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}