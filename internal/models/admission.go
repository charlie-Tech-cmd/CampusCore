package models

import "time"

// AdmissionStatus represents an application's current state.
type AdmissionStatus string

const (
	AdmissionPending  AdmissionStatus = "pending"
	AdmissionApproved AdmissionStatus = "approved"
	AdmissionRejected AdmissionStatus = "rejected"
)

// AdmissionApplication represents a student's admission application.
type AdmissionApplication struct {
	ID                   int             `json:"id"`
	ApplicationNo        string          `json:"application_no"`
	ApplicantName        string          `json:"applicant_name"`
	Email                string          `json:"email"`
	Phone                string          `json:"phone"`
	FacultyID            int             `json:"faculty_id"`
	DepartmentID         int             `json:"department_id"`
	Programme            string          `json:"programme"`
	Session              string          `json:"session"`
	Status               AdmissionStatus `json:"status"`
	AdmissionDate        *time.Time      `json:"admission_date,omitempty"`
	AcceptancePaid       bool            `json:"acceptance_paid"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
	OfferLetterGenerated bool            `json:"offer_letter_generated"`
	AcceptanceFeeAmount  float64         `json:"acceptance_fee_amount"`
	AcceptancePaidAt     *time.Time      `json:"acceptance_paid_at,omitempty"`
}

// AdmissionRepository defines admission operations.
type AdmissionRepository interface {
	SubmitApplication(*AdmissionApplication) error
	FindByApplicationNo(string) (*AdmissionApplication, error)
	ListApplications() ([]AdmissionApplication, error)
	ApproveApplication(string) error
	RejectApplication(string) error
}
