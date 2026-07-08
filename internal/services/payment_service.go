package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"campuscore/internal/models"
)

// PaymentService handles payment operations.
type PaymentService struct {
	repo models.FinancialRepository
}

// NewPaymentService creates a PaymentService.
func NewPaymentService(repo models.FinancialRepository) *PaymentService {
	return &PaymentService{
		repo: repo,
	}
}

// ProcessPayment validates and records a payment.
func (s *PaymentService) ProcessPayment(
	ctx context.Context,
	studentID,
	gatewayRef string,
	amount float64,
	feeType,
	session string,
) error {

	studentID = strings.TrimSpace(studentID)
	gatewayRef = strings.TrimSpace(gatewayRef)
	feeType = strings.TrimSpace(strings.ToLower(feeType))
	session = strings.TrimSpace(session)

	switch {
	case studentID == "":
		return errors.New("student ID is required")
	case gatewayRef == "":
		return errors.New("payment reference is required")
	case amount <= 0:
		return errors.New("amount must be greater than zero")
	case feeType == "":
		return errors.New("fee type is required")
	case session == "":
		return errors.New("session is required")
	}

	exists, err := s.repo.CheckPaymentExists(ctx, gatewayRef)
	if err != nil {
		return fmt.Errorf("failed to verify payment: %w", err)
	}

	if exists {
		return fmt.Errorf("payment reference %q already exists", gatewayRef)
	}

	payment := &models.FeePayment{
		StudentID:        studentID,
		GatewayReference: gatewayRef,
		AmountPaid:       amount,
		FeeType:          feeType,
		Session:          session,
		Status:           "successful",
	}

	if err := s.repo.RecordPayment(ctx, payment); err != nil {
		return fmt.Errorf("failed to save payment: %w", err)
	}

	return nil
}