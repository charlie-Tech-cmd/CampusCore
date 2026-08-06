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

	GetAdmissionReport() (*AdmissionReport, error)
}

type AdmissionReport struct {
	TotalApplications int `json:"total_applications"`

	PendingApplications int `json:"pending_applications"`

	ApprovedApplications int `json:"approved_applications"`

	RejectedApplications int `json:"rejected_applications"`

	ApprovalRate float64 `json:"approval_rate"`
}
