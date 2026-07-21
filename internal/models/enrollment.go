package models

import "time"

// Enrollment represents a student's registration for a course.
type Enrollment struct {
	ID          int       `json:"id" db:"id"`
	StudentID   string    `json:"student_id" db:"student_id"`
	CourseCode  string    `json:"course_code" db:"course_code"`
	CreditUnits int       `json:"credit_units" db:"credit_units"`
	Session     string    `json:"session" db:"session"`
	Semester    string    `json:"semester" db:"semester"`
	Status      string    `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// EnrollmentRepository defines the contract for student course enrollment operations.
type EnrollmentRepository interface {
	Register(enrollment *Enrollment) error
	FindByStudent(studentID string) ([]Enrollment, error)
	FindByCourse(courseCode string) ([]Enrollment, error)
	UpdateStatus(studentID, courseCode, status string) error
	Delete(studentID, courseCode string) error
}
