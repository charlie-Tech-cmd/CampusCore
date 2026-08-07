package services

import (
	"errors"

	"campuscore/internal/models"
)

const MaxSemesterCredits = 24

type RegistrationService struct {
	students    models.UserRepository
	courses     models.CourseRepository
	enrollments models.EnrollmentRepository
}

// NewRegistrationService creates a registration service.
func NewRegistrationService(
	students models.UserRepository,
	courses models.CourseRepository,
	enrollments models.EnrollmentRepository,
) *RegistrationService {
	return &RegistrationService{
		students:    students,
		courses:     courses,
		enrollments: enrollments,
	}
}

// RegisterCourse registers a student for a course.
func (s *RegistrationService) RegisterCourse(
	studentID string,
	courseCode string,
	session string,
	semester string,
) error {

	// Verify student exists.
	_, err := s.students.FindByID(studentID)
	if err != nil {
		return err
	}

	// Verify course exists.
	course, err := s.courses.FindByCode(courseCode)
	if err != nil {
		return err
	}

	// Verify course is active.
	if !course.IsActive {
		return errors.New("course is not active")
	}

	// Check if the student has already registered this course.
	enrollments, err := s.enrollments.FindByStudent(studentID)
	if err != nil {
		return err
	}

	for _, enrollment := range enrollments {
		if enrollment.CourseCode == courseCode &&
			enrollment.Session == session &&
			enrollment.Semester == semester {
			return errors.New("student is already registered for this course")
		}
	}

	// Verify course capacity.
	if course.CurrentEnrolled >= course.MaxCapacity {
		return errors.New("course has reached maximum capacity")
	}

	// Create the enrollment record.
	enrollment := &models.Enrollment{
		StudentID:   studentID,
		CourseCode:  courseCode,
		CreditUnits: course.CreditUnits,
		Session:     session,
		Semester:    semester,
		Status:      "registered",
	}

	if err := s.enrollments.Register(enrollment); err != nil {
		return err
	}

	// Update course enrollment count.
	course.CurrentEnrolled++

	if err := s.courses.Update(course); err != nil {
		return err
	}

	return nil

}

func (s *RegistrationService) GetStudentCourses(
	studentID string,
) ([]models.Enrollment, error) {

	return s.enrollments.FindByStudent(studentID)
}
