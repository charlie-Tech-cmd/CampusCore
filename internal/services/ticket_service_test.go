package services

import (
	"context"
	"errors"
	"testing"

	"campuscore/internal/models"
)

func TestSubmitHelpdeskTicket_Success(t *testing.T) {
	repo := &mockFinancialRepository{
		createTicketFn: func(ctx context.Context, ticket *models.SupportTicket) error {
			if ticket.StudentID != "STU001" {
				t.Fatalf("expected STU001, got %q", ticket.StudentID)
			}

			if ticket.Category != "fees" {
				t.Fatalf("expected fees, got %q", ticket.Category)
			}

			if ticket.Subject != "Payment Issue" {
				t.Fatalf("unexpected subject %q", ticket.Subject)
			}

			if ticket.Message != "Payment not reflected" {
				t.Fatalf("unexpected message %q", ticket.Message)
			}

			if ticket.Status != "open" {
				t.Fatalf("expected default status open, got %q", ticket.Status)
			}

			return nil
		},
	}

	service := NewTicketService(repo)

	ticket := &models.SupportTicket{
		StudentID: "STU001",
		Category:  "fees",
		Subject:   "Payment Issue",
		Message:   "Payment not reflected",
	}

	err := service.SubmitHelpdeskTicket(context.Background(), ticket)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestNewTicketService(t *testing.T) {
	repo := &mockFinancialRepository{}

	service := NewTicketService(repo)

	if service == nil {
		t.Fatal("expected TicketService, got nil")
	}

	if service.finRepo != repo {
		t.Fatal("repository was not assigned correctly")
	}
}

func TestSubmitHelpdeskTicket_NilTicket(t *testing.T) {
	repo := &mockFinancialRepository{}
	service := NewTicketService(repo)

	err := service.SubmitHelpdeskTicket(context.Background(), nil)

	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	expected := "ticket is required"

	if err.Error() != expected {
		t.Fatalf("expected %q, got %q", expected, err.Error())
	}
}

func TestSubmitHelpdeskTicket_ValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		ticket  *models.SupportTicket
		wantErr string
	}{
		{
			name: "missing student ID",
			ticket: &models.SupportTicket{
				Category: "fees",
				Subject:  "Payment",
				Message:  "Help",
			},
			wantErr: "student ID is required",
		},
		{
			name: "missing category",
			ticket: &models.SupportTicket{
				StudentID: "STU001",
				Subject:   "Payment",
				Message:   "Help",
			},
			wantErr: "category is required",
		},
		{
			name: "missing subject",
			ticket: &models.SupportTicket{
				StudentID: "STU001",
				Category:  "fees",
				Message:   "Help",
			},
			wantErr: "subject is required",
		},
		{
			name: "missing message",
			ticket: &models.SupportTicket{
				StudentID: "STU001",
				Category:  "fees",
				Subject:   "Payment",
			},
			wantErr: "message is required",
		},
	}

	service := NewTicketService(&mockFinancialRepository{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.SubmitHelpdeskTicket(context.Background(), tt.ticket)

			if err == nil {
				t.Fatal("expected an error, got nil")
			}

			if err.Error() != tt.wantErr {
				t.Fatalf("expected %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

func TestSubmitHelpdeskTicket_DefaultStatus(t *testing.T) {
	repo := &mockFinancialRepository{
		createTicketFn: func(ctx context.Context, ticket *models.SupportTicket) error {
			if ticket.Status != "open" {
				t.Fatalf("expected status 'open', got %q", ticket.Status)
			}
			return nil
		},
	}

	service := NewTicketService(repo)

	ticket := &models.SupportTicket{
		StudentID: "STU001",
		Category:  "fees",
		Subject:   "Payment Issue",
		Message:   "Payment not reflected",
	}

	if err := service.SubmitHelpdeskTicket(context.Background(), ticket); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestSubmitHelpdeskTicket_TrimsInput(t *testing.T) {
	repo := &mockFinancialRepository{
		createTicketFn: func(ctx context.Context, ticket *models.SupportTicket) error {
			if ticket.StudentID != "STU001" {
				t.Fatalf("expected trimmed student ID, got %q", ticket.StudentID)
			}

			if ticket.Category != "fees" {
				t.Fatalf("expected trimmed category, got %q", ticket.Category)
			}

			if ticket.Subject != "Payment Issue" {
				t.Fatalf("expected trimmed subject, got %q", ticket.Subject)
			}

			if ticket.Message != "Payment not reflected" {
				t.Fatalf("expected trimmed message, got %q", ticket.Message)
			}

			return nil
		},
	}

	service := NewTicketService(repo)

	ticket := &models.SupportTicket{
		StudentID: "  STU001  ",
		Category:  "  fees  ",
		Subject:   "  Payment Issue  ",
		Message:   "  Payment not reflected  ",
	}

	if err := service.SubmitHelpdeskTicket(context.Background(), ticket); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestSubmitHelpdeskTicket_CreateTicketError(t *testing.T) {
	repo := &mockFinancialRepository{
		createTicketFn: func(ctx context.Context, ticket *models.SupportTicket) error {
			return errors.New("database unavailable")
		},
	}

	service := NewTicketService(repo)

	ticket := &models.SupportTicket{
		StudentID: "STU001",
		Category:  "fees",
		Subject:   "Payment Issue",
		Message:   "Payment not reflected",
	}

	err := service.SubmitHelpdeskTicket(context.Background(), ticket)

	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	expected := "database unavailable"

	if err.Error() != expected {
		t.Fatalf("expected %q, got %q", expected, err.Error())
	}
}