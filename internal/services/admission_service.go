package services

import (
	"errors"
	"strings"
	"time"

	"campuscore/internal/models"
)

type BillingManager interface {
	GenerateInvoice(
		ownerID string,
		ownerType string,
		feeType string,
		session string,
		amount float64,
		dueDate time.Time,
	) error
}

type AdmissionService struct {
	repo         models.AdmissionRepository
	notification *NotificationService
	billing      BillingManager
}

// NewAdmissionService creates a new admission service.
func NewAdmissionService(
	repo models.AdmissionRepository,
	notification *NotificationService,
	billing BillingManager,
) *AdmissionService {

	return &AdmissionService{
		repo:         repo,
		notification: notification,
		billing:      billing,
	}
}

// SubmitApplication validates and submits an admission application.
func (s *AdmissionService) SubmitApplication(
	application *models.AdmissionApplication,
) error {

	if application == nil {
		return errors.New("application is required")
	}

	application.ApplicationNo = strings.TrimSpace(application.ApplicationNo)
	application.ApplicantName = strings.TrimSpace(application.ApplicantName)
	application.Email = strings.TrimSpace(application.Email)
	application.Phone = strings.TrimSpace(application.Phone)
	application.Programme = strings.TrimSpace(application.Programme)
	application.Session = strings.TrimSpace(application.Session)

	switch {
	case application.ApplicationNo == "":
		return errors.New("application number is required")

	case application.ApplicantName == "":
		return errors.New("applicant name is required")

	case application.Email == "":
		return errors.New("email is required")

	case application.Programme == "":
		return errors.New("programme is required")

	case application.Session == "":
		return errors.New("session is required")
	}

	return s.repo.SubmitApplication(application)
}

// GetApplication retrieves an application.
func (s *AdmissionService) GetApplication(
	applicationNo string,
) (*models.AdmissionApplication, error) {

	applicationNo = strings.TrimSpace(applicationNo)

	if applicationNo == "" {
		return nil, errors.New("application number is required")
	}

	return s.repo.FindByApplicationNo(applicationNo)
}

// ListApplications returns all applications.
func (s *AdmissionService) ListApplications() (
	[]models.AdmissionApplication,
	error,
) {

	return s.repo.ListApplications()
}

// ApproveApplication approves an application.
func (s *AdmissionService) ApproveApplication(
	applicationNo string,
) error {

	applicationNo = strings.TrimSpace(applicationNo)

	if applicationNo == "" {
		return errors.New("application number is required")
	}

	err := s.repo.ApproveApplication(applicationNo)
	if err != nil {
		return err
	}

	application, err := s.repo.FindByApplicationNo(applicationNo)
	if err != nil {
		return err
	}

	if s.billing != nil {
		err := s.billing.GenerateInvoice(
			application.ApplicationNo,
			"applicant",
			"acceptance_fee",
			application.Session,
			application.AcceptanceFeeAmount,
			time.Now().AddDate(0, 0, 14),
		)
		if err != nil {
			return err
		}
	}

	if s.notification != nil {

		if s.notification != nil {
			if err := s.notification.NotifyAdmissionApproved(
				application.Email,
			); err != nil {
				return err
			}
		}
	}

	// Notification integration will come next.
	return nil
}

// RejectApplication rejects an application.
func (s *AdmissionService) RejectApplication(
	applicationNo string,
) error {

	applicationNo = strings.TrimSpace(applicationNo)

	if applicationNo == "" {
		return errors.New("application number is required")
	}

	return s.repo.RejectApplication(applicationNo)
}
