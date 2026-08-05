package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"campuscore/internal/models"
)

type BillingService struct {
	repo models.BillingRepository
}

func NewBillingService(
	repo models.BillingRepository,
) *BillingService {

	return &BillingService{
		repo: repo,
	}
}

// CreateInvoice creates a new invoice.
func (s *BillingService) GenerateInvoice(
	ownerID string,
	ownerType string,
	feeType string,
	session string,
	amount float64,
	dueDate time.Time,
) error {

	ownerID = strings.TrimSpace(ownerID)
	ownerType = strings.TrimSpace(strings.ToLower(ownerType))
	feeType = strings.TrimSpace(strings.ToLower(feeType))
	session = strings.TrimSpace(session)

	switch {

	case ownerID == "":
		return errors.New("owner ID is required")

	case ownerType == "":
		return errors.New("owner type is required")

	case feeType == "":
		return errors.New("fee type is required")

	case session == "":
		return errors.New("session is required")

	case amount <= 0:
		return errors.New("amount must be greater than zero")
	}

	invoiceNumber := s.generateInvoiceNumber()
	invoice := &models.Invoice{
		InvoiceNumber: invoiceNumber,

		OwnerID:   ownerID,
		OwnerType: ownerType,
		FeeType:   feeType,
		Session:   session,
		Amount:    amount,

		Status: models.InvoicePending,

		DueDate:   dueDate,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return s.repo.CreateInvoice(invoice)
}

func (s *BillingService) generateInvoiceNumber() string {
	return fmt.Sprintf(
		"INV-%d-%06d",
		time.Now().Year(),
		time.Now().Unix()%1000000,
	)
}
