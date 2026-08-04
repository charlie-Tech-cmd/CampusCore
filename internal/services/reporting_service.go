package services

import (
	"campuscore/internal/models"
)

// ReportingService handles reporting and analytics.
type ReportingService struct {
	repo models.ReportingRepository
}

// NewReportingService creates a reporting service.
func NewReportingService(
	repo models.ReportingRepository,
) *ReportingService {

	return &ReportingService{
		repo: repo,
	}
}

// GetEnrollmentSummary returns enrollment statistics.
func (s *ReportingService) GetEnrollmentSummary() (
	*models.EnrollmentSummary,
	error,
) {
	return s.repo.GetEnrollmentSummary()
}

// GetPaymentSummary returns payment statistics.
func (s *ReportingService) GetPaymentSummary() (
	*models.PaymentSummary,
	error,
) {
	return s.repo.GetPaymentSummary()
}

// GetAcademicPerformanceSummary returns academic statistics.
func (s *ReportingService) GetAcademicPerformanceSummary() (
	*models.AcademicPerformanceSummary,
	error,
) {
	return s.repo.GetAcademicPerformanceSummary()
}

// GetClearanceSummary returns clearance statistics.
func (s *ReportingService) GetClearanceSummary() (
	*models.ClearanceSummary,
	error,
) {
	return s.repo.GetClearanceSummary()
}
