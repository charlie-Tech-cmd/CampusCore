package services

import (
	"context"
	"errors"
	"strings"

	"campuscore/internal/models"
)

// TicketService handles support tickets.
type TicketService struct {
	finRepo models.FinancialRepository
}

// NewTicketService creates a TicketService.
func NewTicketService(finRepo models.FinancialRepository) *TicketService {
	return &TicketService{
		finRepo: finRepo,
	}
}

// SubmitHelpdeskTicket validates and saves a support ticket.
func (s *TicketService) SubmitHelpdeskTicket(
	ctx context.Context,
	ticket *models.SupportTicket,
) error {

	if ticket == nil {
		return errors.New("ticket is required")
	}

	// Clean the input.
	ticket.StudentID = strings.TrimSpace(ticket.StudentID)
	ticket.Category = strings.TrimSpace(ticket.Category)
	ticket.Subject = strings.TrimSpace(ticket.Subject)
	ticket.Message = strings.TrimSpace(ticket.Message)

	switch {
	case ticket.StudentID == "":
		return errors.New("student ID is required")
	case ticket.Category == "":
		return errors.New("category is required")
	case ticket.Subject == "":
		return errors.New("subject is required")
	case ticket.Message == "":
		return errors.New("message is required")
	}

	if ticket.Status == "" {
		ticket.Status = "open"
	}

	return s.finRepo.CreateTicket(ctx, ticket)
}
