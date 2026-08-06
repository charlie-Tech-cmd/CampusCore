package services

import (
	"errors"
	"strconv"
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
	courses     models.CourseRepository
	results     ResultRepository
	departments DepartmentRepository
	faculties   FacultyRepository
}

// NewTranscriptService creates a transcript service.
func NewTranscriptService(
	users UserRepository,
	courses models.CourseRepository,
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

	semesters := s.buildSemesterTranscripts(results)

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

		Semesters:   semesters,
		GeneratedAt: time.Now(),
	}

	return transcript, nil

}

func (s *TranscriptService) buildSemesterTranscripts(
	results []models.Result,
) []models.SemesterTranscript {

	semesterMap := make(map[string]*models.SemesterTranscript)

	for _, result := range results {

		course, err := s.courses.FindByCode(result.CourseCode)
		if err != nil {
			// Skip results whose course cannot be found.
			continue
		}

		key := result.Session + "-" +
			result.Semester + "-" +
			strconv.Itoa(course.Level)

		semester, exists := semesterMap[key]
		if !exists {
			semester = &models.SemesterTranscript{
				Session:  result.Session,
				Semester: result.Semester,
				Level:    course.Level,
			}

			semesterMap[key] = semester
		}

		entry := models.TranscriptEntry{
			CourseCode:  result.CourseCode,
			CourseTitle: course.Title,
			CreditUnits: result.CreditUnits,
			Score:       result.Score,
			Grade:       result.Grade,
			GradePoint:  result.GradePoint,
			Session:     result.Session,
			Semester:    result.Semester,
			Level:       course.Level,
		}

		semester.Courses = append(
			semester.Courses,
			entry,
		)

		semester.TotalUnits += result.CreditUnits

		semester.QualityPoints += academic.CalculateQualityPoints(
			result.GradePoint,
			result.CreditUnits,
		)
	}

	var semesters []models.SemesterTranscript

	for _, semester := range semesterMap {
		semester.GPA = academic.CalculateGPA(
			semester.QualityPoints,
			semester.TotalUnits,
		)

		semesters = append(
			semesters,
			*semester,
		)
	}

	return semesters
}
