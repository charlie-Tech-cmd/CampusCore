package governance

import (
	"errors"
	"testing"

	"campuscore/internal/models"
)

type mockGovernanceRepository struct {
	approval       models.Approval
	getErr         error
	updateErr      error
	updatedState   models.ResultStatus
	updatedStaffID string
	updatedRemarks string
	updateCalls    int
}

func (m *mockGovernanceRepository) GetApprovalStatus(courseCode string) (*models.Approval, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	return &m.approval, nil
}

func (m *mockGovernanceRepository) UpdateApprovalState(
	courseCode string,
	nextState models.ResultStatus,
	staffID string,
	remarks string,
) error {
	m.updateCalls++
	m.updatedState = nextState
	m.updatedStaffID = staffID
	m.updatedRemarks = remarks

	return m.updateErr
}

func TestProcessApprovalAdvance(t *testing.T) {
	tests := []struct {
		name      string
		state     models.ResultStatus
		role      models.UserRole
		wantState models.ResultStatus
		wantErr   bool
	}{
		{
			name:      "admin advances submitted",
			state:     models.StatusSubmitted,
			role:      models.RoleAdmin,
			wantState: models.StatusHODApproved,
		},
		{
			name:      "lecturer advances HOD approved",
			state:     models.StatusHODApproved,
			role:      models.RoleLecturer,
			wantState: models.StatusDeanApproved,
		},
		{
			name:      "admin advances HOD approved",
			state:     models.StatusHODApproved,
			role:      models.RoleAdmin,
			wantState: models.StatusDeanApproved,
		},
		{
			name:      "admin advances dean approved",
			state:     models.StatusDeanApproved,
			role:      models.RoleAdmin,
			wantState: models.StatusSenateApproved,
		},
		{
			name:      "admin finalizes senate approved",
			state:     models.StatusSenateApproved,
			role:      models.RoleAdmin,
			wantState: models.StatusFinalized,
		},
		{
			name:    "student cannot advance submitted",
			state:   models.StatusSubmitted,
			role:    models.RoleStudent,
			wantErr: true,
		},
		{
			name:    "lecturer cannot advance submitted",
			state:   models.StatusSubmitted,
			role:    models.RoleLecturer,
			wantErr: true,
		},
		{
			name:    "student cannot advance HOD approved",
			state:   models.StatusHODApproved,
			role:    models.RoleStudent,
			wantErr: true,
		},
		{
			name:    "bursar cannot advance HOD approved",
			state:   models.StatusHODApproved,
			role:    models.RoleBursar,
			wantErr: true,
		},
		{
			name:    "lecturer cannot advance dean approved",
			state:   models.StatusDeanApproved,
			role:    models.RoleLecturer,
			wantErr: true,
		},
		{
			name:    "lecturer cannot finalize",
			state:   models.StatusSenateApproved,
			role:    models.RoleLecturer,
			wantErr: true,
		},
		{
			name:    "finalized record cannot advance",
			state:   models.StatusFinalized,
			role:    models.RoleAdmin,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockGovernanceRepository{
				approval: models.Approval{
					CurrentState: tt.state,
				},
			}

			engine := NewEngine(repo)

			err := engine.ProcessApprovalAdvance(
				"CSC101",
				tt.role,
				"STAFF001",
			)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}

				if repo.updateCalls != 0 {
					t.Fatalf(
						"expected repository not to be updated, got %d calls",
						repo.updateCalls,
					)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if repo.updateCalls != 1 {
				t.Fatalf(
					"expected exactly one repository update, got %d",
					repo.updateCalls,
				)
			}

			if repo.updatedState != tt.wantState {
				t.Fatalf(
					"expected next state %q, got %q",
					tt.wantState,
					repo.updatedState,
				)
			}

			if repo.updatedStaffID != "STAFF001" {
				t.Fatalf(
					"expected staff ID STAFF001, got %q",
					repo.updatedStaffID,
				)
			}
		})
	}
}

func TestProcessApprovalRejection(t *testing.T) {
	tests := []struct {
		name      string
		state     models.ResultStatus
		role      models.UserRole
		remarks   string
		wantState models.ResultStatus
		wantErr   bool
	}{
		{
			name:      "lecturer rejects HOD approved result",
			state:     models.StatusHODApproved,
			role:      models.RoleLecturer,
			remarks:   "Incorrect departmental approval",
			wantState: models.StatusSubmitted,
		},
		{
			name:      "admin rejects HOD approved result",
			state:     models.StatusHODApproved,
			role:      models.RoleAdmin,
			remarks:   "Incorrect departmental approval",
			wantState: models.StatusSubmitted,
		},
		{
			name:      "admin rejects dean approved result",
			state:     models.StatusDeanApproved,
			role:      models.RoleAdmin,
			remarks:   "Dean approval requires correction",
			wantState: models.StatusHODApproved,
		},
		{
			name:      "admin rejects senate approved result",
			state:     models.StatusSenateApproved,
			role:      models.RoleAdmin,
			remarks:   "Senate review identified an issue",
			wantState: models.StatusDeanApproved,
		},
		{
			name:    "student cannot reject HOD approved result",
			state:   models.StatusHODApproved,
			role:    models.RoleStudent,
			remarks: "This is a valid rejection reason",
			wantErr: true,
		},
		{
			name:    "bursar cannot reject HOD approved result",
			state:   models.StatusHODApproved,
			role:    models.RoleBursar,
			remarks: "This is a valid rejection reason",
			wantErr: true,
		},
		{
			name:    "lecturer cannot reject dean approved result",
			state:   models.StatusDeanApproved,
			role:    models.RoleLecturer,
			remarks: "This is a valid rejection reason",
			wantErr: true,
		},
		{
			name:    "lecturer cannot reject senate approved result",
			state:   models.StatusSenateApproved,
			role:    models.RoleLecturer,
			remarks: "This is a valid rejection reason",
			wantErr: true,
		},
		{
			name:    "cannot reject submitted result",
			state:   models.StatusSubmitted,
			role:    models.RoleAdmin,
			remarks: "This is a valid rejection reason",
			wantErr: true,
		},
		{
			name:    "cannot reject finalized result",
			state:   models.StatusFinalized,
			role:    models.RoleAdmin,
			remarks: "This is a valid rejection reason",
			wantErr: true,
		},
		{
			name:    "reject requires minimum remarks",
			state:   models.StatusHODApproved,
			role:    models.RoleAdmin,
			remarks: "short",
			wantErr: true,
		},
		{
			name:    "reject requires non-empty remarks",
			state:   models.StatusHODApproved,
			role:    models.RoleAdmin,
			remarks: "         ",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockGovernanceRepository{
				approval: models.Approval{
					CurrentState: tt.state,
				},
			}

			engine := NewEngine(repo)

			err := engine.ProcessApprovalRejection(
				"CSC101",
				tt.role,
				"STAFF001",
				tt.remarks,
			)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}

				if repo.updateCalls != 0 {
					t.Fatalf(
						"expected repository not to be updated, got %d calls",
						repo.updateCalls,
					)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if repo.updateCalls != 1 {
				t.Fatalf(
					"expected exactly one repository update, got %d",
					repo.updateCalls,
				)
			}

			if repo.updatedState != tt.wantState {
				t.Fatalf(
					"expected rejection target state %q, got %q",
					tt.wantState,
					repo.updatedState,
				)
			}

			if repo.updatedStaffID != "STAFF001" {
				t.Fatalf(
					"expected staff ID STAFF001, got %q",
					repo.updatedStaffID,
				)
			}

			if repo.updatedRemarks != tt.remarks {
				t.Fatalf(
					"expected remarks %q, got %q",
					tt.remarks,
					repo.updatedRemarks,
				)
			}
		})
	}
}

func TestProcessApprovalAdvanceRepositoryError(t *testing.T) {
	repo := &mockGovernanceRepository{
		approval: models.Approval{
			CurrentState: models.StatusSubmitted,
		},
		updateErr: errors.New("database failure"),
	}

	engine := NewEngine(repo)

	err := engine.ProcessApprovalAdvance(
		"CSC101",
		models.RoleAdmin,
		"STAFF001",
	)

	if err == nil {
		t.Fatal("expected repository error, got nil")
	}

	if repo.updateCalls != 1 {
		t.Fatalf("expected one repository update, got %d", repo.updateCalls)
	}
}

func TestProcessApprovalRejectionRepositoryError(t *testing.T) {
	repo := &mockGovernanceRepository{
		approval: models.Approval{
			CurrentState: models.StatusHODApproved,
		},
		updateErr: errors.New("database failure"),
	}

	engine := NewEngine(repo)

	err := engine.ProcessApprovalRejection(
		"CSC101",
		models.RoleAdmin,
		"STAFF001",
		"This is a valid rejection reason",
	)

	if err == nil {
		t.Fatal("expected repository error, got nil")
	}

	if repo.updateCalls != 1 {
		t.Fatalf("expected one repository update, got %d", repo.updateCalls)
	}
}

func TestProcessApprovalAdvanceInvalidState(t *testing.T) {
	repo := &mockGovernanceRepository{
		approval: models.Approval{
			CurrentState: models.ResultStatus("invalid_state"),
		},
	}

	engine := NewEngine(repo)

	err := engine.ProcessApprovalAdvance(
		"CSC101",
		models.RoleAdmin,
		"STAFF001",
	)

	if err == nil {
		t.Fatal("expected invalid state error, got nil")
	}

	if repo.updateCalls != 0 {
		t.Fatalf(
			"expected repository not to be updated, got %d calls",
			repo.updateCalls,
		)
	}
}

func TestProcessApprovalRejectionInvalidState(t *testing.T) {
	repo := &mockGovernanceRepository{
		approval: models.Approval{
			CurrentState: models.ResultStatus("invalid_state"),
		},
	}

	engine := NewEngine(repo)

	err := engine.ProcessApprovalRejection(
		"CSC101",
		models.RoleAdmin,
		"STAFF001",
		"This is a valid rejection reason",
	)

	if err == nil {
		t.Fatal("expected invalid state error, got nil")
	}

	if repo.updateCalls != 0 {
		t.Fatalf(
			"expected repository not to be updated, got %d calls",
			repo.updateCalls,
		)
	}
}
