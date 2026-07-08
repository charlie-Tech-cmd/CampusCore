package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"campuscore/internal/models"
)

// FinancialService contains payment and clearance logic.
type FinancialService struct {
	repo models.FinancialRepository
	db   *sql.DB
}

// NewFinancialService creates a FinancialService.
func NewFinancialService(repo models.FinancialRepository, db *sql.DB) *FinancialService {
	return &FinancialService{
		repo: repo,
		db:   db,
	}
}

// VerifyTuitionClearance checks whether a student has paid the required fees.
func (s *FinancialService) VerifyTuitionClearance(
	ctx context.Context,
	studentID,
	session string,
) (bool, error) {

	var departmentID, level int

	query := `
		SELECT department_id, level
		FROM users
		WHERE id = $1
		AND role = 'student'
	`

	err := s.db.QueryRowContext(ctx, query, studentID).Scan(&departmentID, &level)
	if err != nil {
		return false, fmt.Errorf("failed to load student profile: %w", err)
	}

	fee, err := s.repo.GetFeeStructure(
		ctx,
		departmentID,
		level,
		"school_fees",
		session,
	)
	if err != nil {
		return false, err
	}

	var totalPaid float64

	query = `
		SELECT COALESCE(SUM(amount_paid),0)
		FROM fee_payments
		WHERE student_id=$1
		AND fee_type='school_fees'
		AND session=$2
		AND status='successful'
	`

	err = s.db.QueryRowContext(ctx, query, studentID, session).Scan(&totalPaid)
	if err != nil {
		return false, err
	}

	return totalPaid >= fee.AmountRequired, nil
}

// ProcessIncomingWebhook records a successful payment.
func (s *FinancialService) ProcessIncomingWebhook(
	ctx context.Context,
	studentID,
	reference string,
	amount float64,
	feeType,
	session string,
) error {

	exists, err := s.repo.CheckPaymentExists(ctx, reference)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("payment %s already exists", reference)
	}

	payment := &models.FeePayment{
		StudentID:        studentID,
		GatewayReference: reference,
		AmountPaid:       amount,
		FeeType:          feeType,
		Session:          session,
		Status:           "successful",
	}

	return s.repo.RecordPayment(ctx, payment)
}

// EvaluateGraduationEligibility checks whether all clearance stages are complete.
func (s *FinancialService) EvaluateGraduationEligibility(
	ctx context.Context,
	studentID string,
) (bool, error) {

	clearances, err := s.repo.GetStudentClearanceStatus(ctx, studentID)
	if err != nil {
		return false, err
	}

	if len(clearances) == 0 {
		return false, errors.New("no clearance records found")
	}

	for _, clearance := range clearances {
		if clearance.Status != models.ClearanceCleared {
			return false, nil
		}
	}

	return true, nil
}