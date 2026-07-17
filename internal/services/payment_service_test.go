package services

import (
	"context"
	"errors"
	"testing"

	"campuscore/internal/models"
)

func TestNewPaymentService(t *testing.T) {
	repo := &mockFinancialRepository{}

	service := NewPaymentService(repo)

	if service == nil {
		t.Fatal("expected PaymentService, got nil")
	}

	if service.repo != repo {
		t.Fatal("repository was not assigned correctly")
	}
}

func TestProcessPayment_Success(t *testing.T) {
	repo := &mockFinancialRepository{
		checkPaymentExistsFn: func(ctx context.Context, ref string) (bool, error) {
			if ref != "PAY123" {
				t.Fatalf("expected payment reference PAY123, got %q", ref)
			}
			return false, nil
		},
		recordPaymentFn: func(ctx context.Context, payment *models.FeePayment) error {
			if payment.StudentID != "STU001" {
				t.Fatalf("expected student ID STU001, got %q", payment.StudentID)
			}

			if payment.GatewayReference != "PAY123" {
				t.Fatalf("expected gateway reference PAY123, got %q", payment.GatewayReference)
			}

			if payment.AmountPaid != 50000 {
				t.Fatalf("expected amount 50000, got %f", payment.AmountPaid)
			}

			if payment.FeeType != "school_fees" {
				t.Fatalf("expected fee type school_fees, got %q", payment.FeeType)
			}

			if payment.Session != "2025/2026" {
				t.Fatalf("expected session 2025/2026, got %q", payment.Session)
			}

			if payment.Status != "successful" {
				t.Fatalf("expected status successful, got %q", payment.Status)
			}

			return nil
		},
	}

	service := NewPaymentService(repo)

	err := service.ProcessPayment(
		context.Background(),
		"STU001",
		"PAY123",
		50000,
		"school_fees",
		"2025/2026",
	)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestProcessPayment_ValidationErrors(t *testing.T) {
	tests := []struct {
		name      string
		studentID string
		ref       string
		amount    float64
		feeType   string
		session   string
		wantErr   string
	}{
		{
			name:    "missing student ID",
			ref:     "PAY123",
			amount:  50000,
			feeType: "school_fees",
			session: "2025/2026",
			wantErr: "student ID is required",
		},
		{
			name:      "missing payment reference",
			studentID: "STU001",
			amount:    50000,
			feeType:   "school_fees",
			session:   "2025/2026",
			wantErr:   "payment reference is required",
		},
		{
			name:      "invalid amount",
			studentID: "STU001",
			ref:       "PAY123",
			amount:    0,
			feeType:   "school_fees",
			session:   "2025/2026",
			wantErr:   "amount must be greater than zero",
		},
		{
			name:      "missing fee type",
			studentID: "STU001",
			ref:       "PAY123",
			amount:    50000,
			session:   "2025/2026",
			wantErr:   "fee type is required",
		},
		{
			name:      "missing session",
			studentID: "STU001",
			ref:       "PAY123",
			amount:    50000,
			feeType:   "school_fees",
			wantErr:   "session is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewPaymentService(&mockFinancialRepository{})

			err := service.ProcessPayment(
				context.Background(),
				tt.studentID,
				tt.ref,
				tt.amount,
				tt.feeType,
				tt.session,
			)

			if err == nil {
				t.Fatal("expected an error, got nil")
			}

			if err.Error() != tt.wantErr {
				t.Fatalf("expected %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

func TestProcessPayment_DuplicateReference(t *testing.T) {
	repo := &mockFinancialRepository{
		checkPaymentExistsFn: func(ctx context.Context, ref string) (bool, error) {
			return true, nil
		},
	}

	service := NewPaymentService(repo)

	err := service.ProcessPayment(
		context.Background(),
		"STU001",
		"PAY123",
		50000,
		"school_fees",
		"2025/2026",
	)

	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	expected := `payment reference "PAY123" already exists`

	if err.Error() != expected {
		t.Fatalf("expected %q, got %q", expected, err.Error())
	}
}

func TestProcessPayment_CheckPaymentExistsError(t *testing.T) {
	repo := &mockFinancialRepository{
		checkPaymentExistsFn: func(ctx context.Context, ref string) (bool, error) {
			return false, errors.New("database unavailable")
		},
	}

	service := NewPaymentService(repo)

	err := service.ProcessPayment(
		context.Background(),
		"STU001",
		"PAY123",
		50000,
		"school_fees",
		"2025/2026",
	)

	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	expected := "failed to verify payment: database unavailable"

	if err.Error() != expected {
		t.Fatalf("expected %q, got %q", expected, err.Error())
	}
}

func TestProcessPayment_RecordPaymentError(t *testing.T) {
	repo := &mockFinancialRepository{
		checkPaymentExistsFn: func(ctx context.Context, ref string) (bool, error) {
			return false, nil
		},
		recordPaymentFn: func(ctx context.Context, payment *models.FeePayment) error {
			return errors.New("insert failed")
		},
	}

	service := NewPaymentService(repo)

	err := service.ProcessPayment(
		context.Background(),
		"STU001",
		"PAY123",
		50000,
		"school_fees",
		"2025/2026",
	)

	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	expected := "failed to save payment: insert failed"

	if err.Error() != expected {
		t.Fatalf("expected %q, got %q", expected, err.Error())
	}
}

func TestProcessPayment_TrimsInput(t *testing.T) {
	repo := &mockFinancialRepository{
		checkPaymentExistsFn: func(ctx context.Context, ref string) (bool, error) {
			if ref != "PAY123" {
				t.Fatalf("expected trimmed reference PAY123, got %q", ref)
			}
			return false, nil
		},
		recordPaymentFn: func(ctx context.Context, payment *models.FeePayment) error {
			if payment.StudentID != "STU001" {
				t.Fatalf("expected trimmed student ID STU001, got %q", payment.StudentID)
			}

			if payment.GatewayReference != "PAY123" {
				t.Fatalf("expected trimmed gateway reference PAY123, got %q", payment.GatewayReference)
			}

			if payment.FeeType != "school_fees" {
				t.Fatalf("expected normalized fee type school_fees, got %q", payment.FeeType)
			}

			if payment.Session != "2025/2026" {
				t.Fatalf("expected trimmed session 2025/2026, got %q", payment.Session)
			}

			return nil
		},
	}

	service := NewPaymentService(repo)

	err := service.ProcessPayment(
		context.Background(),
		"  STU001  ",
		"  PAY123  ",
		50000,
		"  SCHOOL_FEES  ",
		"  2025/2026  ",
	)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}
