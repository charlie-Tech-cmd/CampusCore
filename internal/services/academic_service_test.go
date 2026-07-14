package services

import (
	"testing"
)

func TestNewAcademicService(t *testing.T) {
	service := NewAcademicService(nil)

	if service == nil {
		t.Fatal("expected AcademicService, got nil")
	}

	if service.db != nil {
		t.Fatal("expected nil database")
	}
}

func TestCalculateGradeMetrics(t *testing.T) {
	service := NewAcademicService(nil)

	tests := []struct {
		name  string
		score float64
		grade string
		gpa   float64
	}{
		{"A", 85, "A", 5.0},
		{"B", 65, "B", 4.0},
		{"C", 55, "C", 3.0},
		{"D", 47, "D", 2.0},
		{"E", 40, "E", 1.0},
		{"F", 39, "F", 0.0},
		{"Perfect", 100, "A", 5.0},
		{"Zero", 0, "F", 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grade, gpa := service.CalculateGradeMetrics(tt.score)

			if grade != tt.grade {
				t.Fatalf("expected grade %q, got %q", tt.grade, grade)
			}

			if gpa != tt.gpa {
				t.Fatalf("expected GPA %.1f, got %.1f", tt.gpa, gpa)
			}
		})
	}
}