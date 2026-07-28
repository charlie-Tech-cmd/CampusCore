package api

import (
	"context"
	"encoding/json"
	// "errors"
	"net/http"
	"net/http/httptest"
	// "net/url"
	// "strings"
	"testing"

	"campuscore/internal/models"
)

type mockClearanceService struct {
	getStudentClearanceFn func(
		ctx context.Context,
		studentID string,
	) ([]models.StudentClearance, error)

	updateClearanceFn func(
		ctx context.Context,
		studentID string,
		officeID int,
		status models.ClearanceStatus,
		staffID string,
	) error

	isStudentClearedFn func(
		ctx context.Context,
		studentID string,
	) (bool, error)
}

func (m *mockClearanceService) GetStudentClearance(
	ctx context.Context,
	studentID string,
) ([]models.StudentClearance, error) {

	if m.getStudentClearanceFn != nil {
		return m.getStudentClearanceFn(ctx, studentID)
	}

	return nil, nil
}

func (m *mockClearanceService) UpdateClearance(
	ctx context.Context,
	studentID string,
	officeID int,
	status models.ClearanceStatus,
	staffID string,
) error {

	if m.updateClearanceFn != nil {
		return m.updateClearanceFn(
			ctx,
			studentID,
			officeID,
			status,
			staffID,
		)
	}

	return nil
}

func (m *mockClearanceService) IsStudentCleared(
	ctx context.Context,
	studentID string,
) (bool, error) {

	if m.isStudentClearedFn != nil {
		return m.isStudentClearedFn(ctx, studentID)
	}

	return false, nil
}

func TestGetStudentClearance_Success(t *testing.T) {

	mockService := &mockClearanceService{
		getStudentClearanceFn: func(
			ctx context.Context,
			studentID string,
		) ([]models.StudentClearance, error) {

			return []models.StudentClearance{
				{
					StudentID: studentID,
					OfficeID:  1,
					Status:    models.ClearancePending,
				},
			}, nil
		},
	}

	handler := NewClearanceHandler(mockService)

	req := httptest.NewRequest(
		http.MethodGet,
		"/clearance?student_id=STU001",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.GetStudentClearance(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d got %d",
			http.StatusOK,
			rec.Code,
		)
	}

	var response []models.StudentClearance

	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if len(response) != 1 {
		t.Fatalf("expected one clearance record")
	}

	if response[0].StudentID != "STU001" {
		t.Fatalf("unexpected student id")
	}
}
