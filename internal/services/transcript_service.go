package services

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"campuscore/internal/academic"
	"campuscore/internal/models"
)

// TranscriptService generates student transcripts.
type TranscriptService struct {
	users       models.UserRepository
	courses     models.CourseRepository
	results     models.ResultRepository
	departments models.DepartmentRepository
	faculties   models.FacultyRepository
}

// NewTranscriptService creates a transcript service.
func NewTranscriptService(
	users models.UserRepository,
	courses models.CourseRepository,
	results models.ResultRepository,
	departments models.DepartmentRepository,
	faculties models.FacultyRepository,
) *TranscriptService {
	return &TranscriptService{
		users:       users,
		courses:     courses,
		results:     results,
		departments: departments,
		faculties:   faculties,
	}
}

// GenerateTranscript generates a complete transcript for a student.
func (s *TranscriptService) GenerateTranscript(
	studentID string,
) (*models.Transcript, error) {

	// Verify student exists.
	student, err := s.users.FindByID(studentID)
	if err != nil {
		return nil, err
	}

	// Resolve department and faculty information.
	var departmentName string
	var facultyName string

	if student.DepartmentID > 0 {
		department, err := s.departments.FindByID(student.DepartmentID)
		if err != nil {
			return nil, err
		}

		departmentName = department.Name

		faculty, err := s.faculties.FindByID(department.FacultyID)
		if err != nil {
			return nil, err
		}

		facultyName = faculty.Name
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

	cgpa := academic.CalculateGPA(
		totalQualityPoints,
		totalCreditUnits,
	)

	classification := academic.ClassifyDegree(cgpa)

	// Build semester transcript.
	semesters, err := s.buildSemesterTranscripts(results)
	if err != nil {
		return nil, err
	}

	transcript := &models.Transcript{
		StudentID: student.ID,

		StudentName: strings.TrimSpace(
			student.Surname + " " +
				student.FirstName + " " +
				student.MiddleName,
		),

		MatricNumber: student.ID,

		DepartmentName: departmentName,
		FacultyName:    facultyName,

		CGPA:           cgpa,
		Classification: classification,

		Semesters:   semesters,
		GeneratedAt: time.Now(),
	}

	return transcript, nil
}

// buildSemesterTranscripts groups results into semester transcripts.
func (s *TranscriptService) buildSemesterTranscripts(
	results []models.Result,
) ([]models.SemesterTranscript, error) {

	semesterMap := make(map[string]*models.SemesterTranscript)

	for _, result := range results {

		course, err := s.courses.FindByCode(result.CourseCode)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to find course %s: %w",
				result.CourseCode,
				err,
			)
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

	// Ensure transcript semesters have deterministic ordering.
	sort.Slice(semesters, func(i, j int) bool {
		if semesters[i].Level != semesters[j].Level {
			return semesters[i].Level < semesters[j].Level
		}

		if semesters[i].Session != semesters[j].Session {
			return semesters[i].Session < semesters[j].Session
		}

		return semesters[i].Semester < semesters[j].Semester
	})

	return semesters, nil
}
