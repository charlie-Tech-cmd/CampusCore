package models

// EnrollmentSummary represents student enrollment statistics.
type EnrollmentSummary struct {
	TotalStudents        int `json:"total_students"`
	RegisteredStudents   int `json:"registered_students"`
	PendingRegistrations int `json:"pending_registrations"`
}

// PaymentSummary represents payment statistics.
type PaymentSummary struct {
	TotalPayments int     `json:"total_payments"`
	TotalRevenue  float64 `json:"total_revenue"`
	Outstanding   float64 `json:"outstanding"`
}

// AcademicPerformanceSummary represents overall academic performance.
type AcademicPerformanceSummary struct {
	AverageGPA     float64 `json:"average_gpa"`
	HighestGPA     float64 `json:"highest_gpa"`
	LowestGPA      float64 `json:"lowest_gpa"`
	GraduatingCGPA float64 `json:"graduating_cgpa"`
}

// ClearanceSummary represents clearance statistics.
type ClearanceSummary struct {
	TotalStudents int `json:"total_students"`
	Cleared       int `json:"cleared"`
	Pending       int `json:"pending"`
}

type ReportingRepository interface {
	GetEnrollmentSummary() (*EnrollmentSummary, error)
	GetPaymentSummary() (*PaymentSummary, error)
	GetAcademicPerformanceSummary() (*AcademicPerformanceSummary, error)
	GetClearanceSummary() (*ClearanceSummary, error)
}
