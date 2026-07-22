package services

import (
	"campuscore/internal/models"
	"errors"
)

// TranscriptService generates student transcripts.
type TranscriptService struct {
	students    StudentRepository
	courses     CourseRepository
	results     ResultRepository
	departments DepartmentRepository
	faculties   FacultyRepository
}

// NewTranscriptService creates a transcript service.
func NewTranscriptService(
	students StudentRepository,
	courses CourseRepository,
	results ResultRepository,
	departments DepartmentRepository,
	faculties FacultyRepository,
) *TranscriptService {

	return &TranscriptService{
		students:    students,
		courses:     courses,
		results:     results,
		departments: departments,
		faculties:   faculties,
	}
}

func (s *TranscriptService) GenerateTranscript(studentID string) (*models.Transcript, error) {

	// Verify student exists.
	student, err := s.students.FindByID(studentID)
	if err != nil {
		return nil, err
	}

	// Retrieve student's results.
	results, err := s.results.FindByStudent(studentID)
	if err != nil {
		return nil, err
	}

	// Student must have at least one result.
	if len(results) == 0 {
		return nil, errors.New("no academic results found")
	}

	_ = student
	_ = results

	// TODO:
	// Compute quality points.
	// Compute GPA.
	// Compute CGPA.
	// Determine classification.
	// Build transcript.

	return nil, nil
}
