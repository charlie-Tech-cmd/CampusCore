package repository

import (
	"database/sql"

	"campuscore/internal/models"
)

type PostgresReportRepository struct {
	db *sql.DB
}

func NewPostgresReportRepository(
	db *sql.DB,
) *PostgresReportRepository {

	return &PostgresReportRepository{
		db: db,
	}
}

func (r *PostgresReportRepository) GetDashboardSummary() (
	*models.DashboardSummary,
	error,
) {

	summary := &models.DashboardSummary{}

	err := r.db.QueryRow(`
		SELECT COUNT(*)
		FROM users
		WHERE role = 'student'
	`).Scan(&summary.TotalStudents)

	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(`
	SELECT COUNT(*)
	FROM admission_applications
`).Scan(&summary.TotalApplicants)

	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(`
	SELECT COUNT(*)
	FROM departments
`).Scan(&summary.TotalDepartments)

	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(`
	SELECT COUNT(*)
	FROM courses
`).Scan(&summary.TotalCourses)

	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(`
	SELECT COALESCE(AVG(gp), 0)
	FROM results
	WHERE gp IS NOT NULL
`).Scan(&summary.AverageCGPA)

	if err != nil {
		return nil, err
	}
	return summary, nil

}
