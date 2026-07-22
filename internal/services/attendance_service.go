package services

import (
	"errors"
	"time"

	"campuscore/internal/models"
)

// AttendanceRepository defines attendance persistence operations.
type AttendanceRepository interface {
	Create(record *models.Attendance) error
	FindByID(id int) (*models.Attendance, error)
	GetAll() ([]models.Attendance, error)
	GetByStudent(studentID string) ([]models.Attendance, error)
	GetByCourse(courseCode string) ([]models.Attendance, error)
	GetByLecturer(lecturerID string) ([]models.Attendance, error)
	Update(record *models.Attendance) error
	Delete(id int) error
}

// AttendanceService contains attendance business logic.
type AttendanceService struct {
	attendance AttendanceRepository
}

// NewAttendanceService creates a new AttendanceService.
func NewAttendanceService(
	attendance AttendanceRepository,
) *AttendanceService {
	return &AttendanceService{
		attendance: attendance,
	}
}

// MarkAttendance records attendance for a student.
func (s *AttendanceService) MarkAttendance(
	record *models.Attendance,
) error {

	// Validate attendance status.
	switch record.Status {
	case "present", "absent", "excused":
		// valid
	default:
		return errors.New("invalid attendance status")
	}

	// Automatically timestamp when attendance is marked.
	record.MarkedAt = time.Now()

	// Persist the attendance record.
	if err := s.attendance.Create(record); err != nil {
		return err
	}

	return nil
}

// GetAttendance retrieves an attendance record by ID.
func (s *AttendanceService) GetAttendance(
	id int,
) (*models.Attendance, error) {

	return s.attendance.FindByID(id)
}

// ListAttendance returns all attendance records.
func (s *AttendanceService) ListAttendance() ([]models.Attendance, error) {
	return s.attendance.GetAll()
}

// ListStudentAttendance returns attendance for one student.
func (s *AttendanceService) ListStudentAttendance(
	studentID string,
) ([]models.Attendance, error) {

	return s.attendance.GetByStudent(studentID)
}

// ListCourseAttendance returns attendance for one course.
func (s *AttendanceService) ListCourseAttendance(
	courseCode string,
) ([]models.Attendance, error) {

	return s.attendance.GetByCourse(courseCode)
}

// ListLecturerAttendance returns attendance recorded by a lecturer.
func (s *AttendanceService) ListLecturerAttendance(
	lecturerID string,
) ([]models.Attendance, error) {

	return s.attendance.GetByLecturer(lecturerID)
}

// UpdateAttendance updates an attendance record.
func (s *AttendanceService) UpdateAttendance(
	record *models.Attendance,
) error {

	// Validate attendance status.
	switch record.Status {
	case "present", "absent", "excused":
		// valid
	default:
		return errors.New("invalid attendance status")
	}

	return s.attendance.Update(record)
}

// DeleteAttendance removes an attendance record.
func (s *AttendanceService) DeleteAttendance(
	id int,
) error {

	return s.attendance.Delete(id)
}
