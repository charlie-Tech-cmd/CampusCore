package api

import (
	"context"

	"campuscore/internal/models"
)

type AcademicService interface {
	RegisterCourse(
		studentID string,
		courseCode string,
		session string,
		semester string,
	) error
}

type TicketService interface {
	SubmitHelpdeskTicket(
		ctx context.Context,
		ticket *models.SupportTicket,
	) error
}

type PaymentService interface {
	ProcessPayment(
		ctx context.Context,
		studentID string,
		reference string,
		amount float64,
		feeType string,
		session string,
	) error
}

type GovernanceEngine interface {
	ProcessApprovalAdvance(
		courseCode string,
		role models.UserRole,
		staffID string,
	) error

	ProcessApprovalRejection(
		courseCode string,
		role models.UserRole,
		staffID string,
		remarks string,
	) error
}
