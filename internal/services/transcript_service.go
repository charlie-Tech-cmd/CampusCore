package services

import (
	"errors"
	"strings"
	"time"

	"campuscore/internal/academic"
	"campuscore/internal/models"
)

type UserRepository interface {
	FindByID(id string) (*models.User, error)
}

// TranscriptService generates student transcripts.
type TranscriptService struct {
	users       UserRepository
	courses     CourseRepository
	results     ResultRepository
	departments DepartmentRepository
	faculties   FacultyRepository
}

// NewTranscriptService creates a transcript service.
func NewTranscriptService(
	users UserRepository,
	courses CourseRepository,
	results ResultRepository,
	departments DepartmentRepository,
	faculties FacultyRepository,
) *TranscriptService {
	return &TranscriptService{
		users:       users,
		courses:     courses,
		results:     results,
		departments: departments,
		faculties:   faculties,
	}
}

func (s *TranscriptService) GenerateTranscript(studentID string) (*models.Transcript, error) {

	// Verify student exists.
	student, err := s.users.FindByID(studentID)
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

	var (
		totalQualityPoints float64
		totalCreditUnits   int
	)

	for _, result := range results {
		totalQualityPoints += academic.CalculateQualityPoints(
			result.GradePoint,
			result.CreditUnits,
		)

		totalCreditUnits += result.CreditUnits
	}

	gpa := academic.CalculateGPA(
		totalQualityPoints,
		totalCreditUnits,
	)

	classification := academic.ClassifyDegree(gpa)

	transcript := &models.Transcript{
		StudentID: student.ID,

		StudentName: strings.TrimSpace(
			student.Surname + " " +
				student.FirstName + " " +
				student.MiddleName,
		),

		MatricNumber: student.ID,

		CGPA:           gpa,
		Classification: classification,

		Results: results,

		GeneratedAt: time.Now(),
	}

	return transcript, nil

}
