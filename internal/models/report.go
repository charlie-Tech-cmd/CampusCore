package models

type DashboardSummary struct {
	TotalStudents    int `json:"total_students"`
	TotalApplicants  int `json:"total_applicants"`
	TotalDepartments int `json:"total_departments"`
	TotalCourses     int `json:"total_courses"`

	OutstandingInvoices int     `json:"outstanding_invoices"`
	RevenueCollected    float64 `json:"revenue_collected"`

	AverageCGPA float64 `json:"average_cgpa"`
}

type ReportRepository interface {
	GetDashboardSummary() (*DashboardSummary, error)
}
