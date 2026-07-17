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

type mockPaymentService struct {
	processPaymentFunc func(
		context.Context,
		string,
		string,
		float64,
		string,
		string,
	) error
}

func (m *mockPaymentService) ProcessPayment(
	ctx context.Context,
	studentID,
	reference string,
	amount float64,
	feeType,
	session string,
) error {

	if m.processPaymentFunc != nil {
		return m.processPaymentFunc(
			ctx,
			studentID,
			reference,
			amount,
			feeType,
			session,
		)
	}

	return nil
}
