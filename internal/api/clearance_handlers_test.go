package api

import (
	"context"
	"testing"

	"campuscore/internal/models"
)

func TestNewClearanceService(t *testing.T) {
	repo := &mockFinancialRepository{}

	service := NewClearanceService(repo)

	if service == nil {
		t.Fatal("expected service")
	}

	if service.repo == nil {
		t.Fatal("expected repository to be assigned")
	}
}


func TestGetChecklistStatus_EmptyStudentID(t *testing.T) {
	service := NewClearanceService(&mockFinancialRepository{})

	_, err := service.GetChecklistStatus(context.Background(), "")

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetChecklistStatus_Success(t *testing.T) {
	expected := []models.StudentClearance{
		{
			StudentID: "STU001",
		},
	}

	repo := &mockFinancialRepository{
		getStudentClearanceStatusFunc: func(
			ctx context.Context,
			studentID string,
		) ([]models.StudentClearance, error) {

			if studentID != "STU001" {
				t.Fatal("wrong student id")
			}

			return expected, nil
		},
	}

	service := NewClearanceService(repo)

	result, err := service.GetChecklistStatus(
		context.Background(),
		"STU001",
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatal("unexpected result")
	}
}

func TestProcessOfficeSignOff_Success(t *testing.T) {
	repo := &mockFinancialRepository{
		updateClearanceStatusFunc: func(
			ctx context.Context,
			studentID string,
			officeID int,
			status models.ClearanceStatus,
			staffID string,
		) error {

			if studentID != "STU001" {
				t.Fatalf("unexpected studentID: %s", studentID)
			}

			if officeID != 1 {
				t.Fatalf("unexpected officeID: %d", officeID)
			}

			if status != models.ClearanceCleared {
				t.Fatalf("unexpected status: %v", status)
			}

			if staffID != "STAFF001" {
				t.Fatalf("unexpected staffID: %s", staffID)
			}

			return nil
		},
	}

	service := NewClearanceService(repo)

	err := service.ProcessOfficeSignOff(
		context.Background(),
		"STU001",
		1,
		string(models.ClearanceCleared),
		"STAFF001",
	)

	if err != nil {
		t.Fatal(err)
	}
}

func TestProcessOfficeSignOff_EmptyStudentID(t *testing.T) {
	service := NewClearanceService(&mockFinancialRepository{})

	err := service.ProcessOfficeSignOff(
		context.Background(),
		"",
		1,
		string(models.ClearancePending),
		"STAFF001",
	)

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestProcessOfficeSignOff_EmptyStaffID(t *testing.T) {
	service := NewClearanceService(&mockFinancialRepository{})

	err := service.ProcessOfficeSignOff(
		context.Background(),
		"STU001",
		1,
		string(models.ClearancePending),
		"",
	)

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestProcessOfficeSignOff_InvalidStatus(t *testing.T) {
	service := NewClearanceService(&mockFinancialRepository{})

	err := service.ProcessOfficeSignOff(
		context.Background(),
		"STU001",
		1,
		"INVALID_STATUS",
		"STAFF001",
	)

	if err == nil {
		t.Fatal("expected error")
	}
}