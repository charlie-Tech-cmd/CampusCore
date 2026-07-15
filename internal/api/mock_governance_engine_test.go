package api

import "campuscore/internal/models"

type mockGovernanceEngine struct {
	processApprovalAdvanceFunc   func(string, models.UserRole, string) error
	processApprovalRejectionFunc func(string, models.UserRole, string, string) error
}

func (m *mockGovernanceEngine) ProcessApprovalAdvance(
	courseCode string,
	role models.UserRole,
	staffID string,
) error {
	if m.processApprovalAdvanceFunc != nil {
		return m.processApprovalAdvanceFunc(courseCode, role, staffID)
	}
	return nil
}

func (m *mockGovernanceEngine) ProcessApprovalRejection(
	courseCode string,
	role models.UserRole,
	staffID string,
	remarks string,
) error {
	if m.processApprovalRejectionFunc != nil {
		return m.processApprovalRejectionFunc(courseCode, role, staffID, remarks)
	}
	return nil
}