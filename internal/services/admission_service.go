package services

import (
	"errors"
	"strings"

	"campuscore/internal/models"
)

type AdmissionService struct {
	repo         models.AdmissionRepository
	notification *NotificationService
}

// NewAdmissionService creates a new admission service.
func NewAdmissionService(
	repo models.AdmissionRepository,
	notification *NotificationService,
) *AdmissionService {

	return &AdmissionService{
		repo:         repo,
		notification: notification,
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
