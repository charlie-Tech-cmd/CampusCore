package models

import "time"

// ResultStatus represents our strict database state machine enum values
type ResultStatus string

const (
	StatusSubmitted      ResultStatus = "submitted"
	StatusHODApproved    ResultStatus = "hod_approved"
	StatusDeanApproved   ResultStatus = "dean_approved"
	StatusSenateApproved ResultStatus = "senate_approved"
	StatusFinalized      ResultStatus = "finalized"
)

// Approval represents the structural tracking row for a batch of grades
type Approval struct {
	ID           int          `json:"id"`
	CourseCode   string       `json:"course_code"`
	Session      string       `json:"session"`
	Semester     string       `json:"semester"`
	CurrentState ResultStatus `json:"current_state"`
	ActionBy     string       `json:"action_by"` // Staff ID who updated this state
	Remarks      string       `json:"remarks"`   // Mandatory reason field for rejections
	UpdatedAt    time.Time    `json:"updated_at"`
}

// GovernanceRepository outlines the dynamic, decoupled database interfaces
type GovernanceRepository interface {
	GetApprovalStatus(courseCode string) (*Approval, error)
	UpdateApprovalState(courseCode string, newState ResultStatus, staffID string, remarks string) error
}