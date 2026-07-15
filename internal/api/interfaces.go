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