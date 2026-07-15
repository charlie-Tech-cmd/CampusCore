package api

import "context"

import "campuscore/internal/models"

type mockAcademicService struct {
	registerCourseFunc func(studentID, courseCode, session, semester string) error
}

func (m *mockAcademicService) RegisterCourse(
	studentID,
	courseCode,
	session,
	semester string,
) error {
	if m.registerCourseFunc != nil {
		return m.registerCourseFunc(studentID, courseCode, session, semester)
	}
	return nil
}

type mockTicketService struct {
	submitTicketFunc func(context.Context, *models.SupportTicket) error
}

func (m *mockTicketService) SubmitHelpdeskTicket(
	ctx context.Context,
	ticket *models.SupportTicket,
) error {
	if m.submitTicketFunc != nil {
		return m.submitTicketFunc(ctx, ticket)
	}
	return nil
}