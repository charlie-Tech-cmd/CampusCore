package services

import "campuscore/internal/models"

type DashboardService struct {
	users         models.UserRepository
	enrollments   models.EnrollmentRepository
	billing       models.BillingRepository
	notifications models.NotificationRepository
	results       models.ResultRepository
}

func NewDashboardService(
	users models.UserRepository,
	enrollments models.EnrollmentRepository,
	billing models.BillingRepository,
	notifications models.NotificationRepository,
	results models.ResultRepository,
) *DashboardService {

	return &DashboardService{
		users:         users,
		enrollments:   enrollments,
		billing:       billing,
		notifications: notifications,
		results:       results,
	}
}

func (s *DashboardService) GetStudentDashboard(
	studentID string,
) (*models.StudentDashboard, error) {

	profile, err := s.users.GetProfile(studentID)
	if err != nil {
		return nil, err
	}

	courses, err := s.enrollments.FindByStudent(studentID)
	if err != nil {
		return nil, err
	}

	invoices, err := s.billing.FindOutstandingInvoices(
		studentID,
		"student",
	)
	if err != nil {
		return nil, err
	}

	notifications, err := s.notifications.FindByUser(studentID)
	if err != nil {
		return nil, err
	}

	results, err := s.results.FindByStudent(studentID)
	if err != nil {
		return nil, err
	}

	var totalPoints float64
	var totalUnits int

	for _, r := range results {
		totalPoints += float64(r.CreditUnits) * r.GradePoint
		totalUnits += r.CreditUnits
	}

	var cgpa float64
	if totalUnits > 0 {
		cgpa = totalPoints / float64(totalUnits)
	}

	return &models.StudentDashboard{
		Profile:             profile,
		RegisteredCourses:   courses,
		OutstandingInvoices: invoices,
		Notifications:       notifications,
		RecentResults:       results,
		CGPA:                cgpa,
	}, nil
}
