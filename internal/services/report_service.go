package services

import "campuscore/internal/models"

type ReportService struct {
	repo models.ReportRepository
}

func NewReportService(
	repo models.ReportRepository,
) *ReportService {

	return &ReportService{
		repo: repo,
	}
}

func (s *ReportService) DashboardSummary() (
	*models.DashboardSummary,
	error,
) {

	return s.repo.GetDashboardSummary()
}
