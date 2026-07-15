package api

import (
	"context"

	"campuscore/internal/models"

)


type AcademicService interface {
    RegisterCourse(studentID, courseCode, session, semester string) error
}

type TicketService interface {
    SubmitHelpdeskTicket(ctx context.Context, ticket *models.SupportTicket) error
}

type PaymentService interface {
    ProcessPayment(
        ctx context.Context,
        studentID,
        reference string,
        amount float64,
        feeType,
        session string,
    ) error
}