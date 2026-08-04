package repository

import (
	"database/sql"

	"campuscore/internal/models"
)

// PostgresReportingRepository implements ReportingRepository.
type PostgresReportingRepository struct {
	db *sql.DB
}

// NewPostgresReportingRepository creates a reporting repository.
func NewPostgresReportingRepository(
	db *sql.DB,
) *PostgresReportingRepository {

	return &PostgresReportingRepository{
		db: db,
	}
}

// GetEnrollmentSummary returns enrollment statistics.
func (r *PostgresReportingRepository) GetEnrollmentSummary() (
	*models.EnrollmentSummary,
	error,
) {

	return &models.EnrollmentSummary{}, nil
}

// GetPaymentSummary returns payment statistics.
func (r *PostgresReportingRepository) GetPaymentSummary() (
	*models.PaymentSummary,
	error,
) {

	return &models.PaymentSummary{}, nil
}

// GetAcademicPerformanceSummary returns academic statistics.
func (r *PostgresReportingRepository) GetAcademicPerformanceSummary() (
	*models.AcademicPerformanceSummary,
	error,
) {

	return &models.AcademicPerformanceSummary{}, nil
}

// GetClearanceSummary returns clearance statistics.
func (r *PostgresReportingRepository) GetClearanceSummary() (
	*models.ClearanceSummary,
	error,
) {

	return &models.ClearanceSummary{}, nil
}
