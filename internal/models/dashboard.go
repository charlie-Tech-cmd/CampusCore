package models

type StudentDashboard struct {
	Profile             *User          `json:"profile"`
	RegisteredCourses   []Enrollment   `json:"registered_courses"`
	OutstandingInvoices []Invoice      `json:"outstanding_invoices"`
	Notifications       []Notification `json:"notifications"`
	RecentResults       []Result       `json:"recent_results"`
	CGPA                float64        `json:"cgpa"`
}
