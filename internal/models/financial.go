package models

import (
	"context"
	"time"
)

// ClearanceStatus represents a student's clearance state.
type ClearanceStatus string

const (
	ClearancePending   ClearanceStatus = "pending"
	ClearanceSubmitted ClearanceStatus = "submitted"
	ClearanceCleared   ClearanceStatus = "cleared"
)

// FeeStructure defines the fees for a department and level.
type FeeStructure struct {
	ID             int     `json:"id"`
	DepartmentID   int     `json:"department_id"`
	Level          int     `json:"level"`
	FeeType        string  `json:"fee_type"` // e.g. "school_fees", "faculty_levy"
	AmountRequired float64 `json:"amount_required"`
	Session        string  `json:"session"`
}

// FeePayment stores a student's payment record.
type FeePayment struct {
	ID               int       `json:"id"`
	StudentID        string    `json:"student_id"`
	GatewayReference string    `json:"gateway_reference"` // Unique payment reference.
	AmountPaid       float64   `json:"amount_paid"`
	FeeType          string    `json:"fee_type"`
	Session          string    `json:"session"`
	Status           string    `json:"status"` // "pending", "successful", "failed"
	PaidAt           time.Time `json:"paid_at"`
}

// StudentClearance tracks a student's clearance progress.
type StudentClearance struct {
	ID            int             `json:"id"`
	StudentID     string          `json:"student_id"`
	OfficeID      int             `json:"office_id"`
	OfficeName    string          `json:"office_name,omitempty"` // Filled when joined with offices.
	Status        ClearanceStatus `json:"status"`
	AssignedStaff string          `json:"assigned_staff_id,omitempty"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

// SupportTicket represents a student support request.
type SupportTicket struct {
	ID         int       `json:"id"`
	StudentID  string    `json:"student_id"`
	Category   string    `json:"category"` // "fees", "registration", "biodata", "general"
	Status     string    `json:"status"`   // "open", "resolved"
	Subject    string    `json:"subject"`
	Message    string    `json:"message"`
	ResolvedBy string    `json:"resolved_by,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// FinancialRepository defines financial data operations.
type FinancialRepository interface {
	GetFeeStructure(
		ctx context.Context,
		departmentID, level int,
		feeType, session string,
	) (*FeeStructure, error)

	RecordPayment(
		ctx context.Context,
		payment *FeePayment,
	) error

	CheckPaymentExists(
		ctx context.Context,
		gatewayRef string,
	) (bool, error)

	GetStudentClearanceStatus(
		ctx context.Context,
		studentID string,
	) ([]StudentClearance, error)

	UpdateClearanceStatus(
		ctx context.Context,
		studentID string,
		officeID int,
		status ClearanceStatus,
		staffID string,
	) error

	CreateTicket(
		ctx context.Context,
		ticket *SupportTicket,
	) error
}