package services

import (
	"context"
	"testing"

	"campuscore/internal/models"
)

func TestNewClearanceService(t *testing.T) {
	repo := &mockFinancialRepository{}

	service := NewClearanceService(repo)

	if service == nil {
		t.Fatal("expected service, got nil")
	}

	if service.repo != repo {
		t.Fatal("repository was not assigned")
	}
}

func TestGetStudentClearance_Success(t *testing.T) {
	expected := []models.StudentClearance{
		{
			StudentID: "STU001",
			OfficeID:  1,
			Status:    models.ClearancePending,
		},
	}

	repo := &mockFinancialRepository{
		getStudentClearanceStatusFn: func(
			ctx context.Context,
			studentID string,
		) ([]models.StudentClearance, error) {

			if studentID != "STU001" {
				t.Fatalf("expected STU001, got %q", studentID)
			}

			return expected, nil
		},
	}

	service := NewClearanceService(repo)

	result, err := service.GetStudentClearance(
		context.Background(),
		"STU001",
	)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 clearance record, got %d", len(result))
	}

	if result[0].StudentID != "STU001" {
		t.Fatalf("unexpected student ID %q", result[0].StudentID)
	}
}