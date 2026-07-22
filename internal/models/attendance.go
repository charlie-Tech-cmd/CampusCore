package models

import "time"

// Attendance represents a student's attendance for a scheduled class.
type Attendance struct {
	ID         int       `json:"id" db:"id"`
	StudentID  string    `json:"student_id" db:"student_id"`
	CourseCode string    `json:"course_code" db:"course_code"`
	LecturerID string    `json:"lecturer_id" db:"lecturer_id"`
	Session    string    `json:"session" db:"session"`
	Semester   string    `json:"semester" db:"semester"`
	ClassDate  time.Time `json:"class_date" db:"class_date"`
	Status     string    `json:"status" db:"status"` // present, absent, excused
	MarkedAt   time.Time `json:"marked_at" db:"marked_at"`
}

// AttendanceRepository defines attendance persistence operations.
type AttendanceRepository interface {
	Create(record *Attendance) error
	FindByID(id int) (*Attendance, error)
	GetAll() ([]Attendance, error)
	GetByStudent(studentID string) ([]Attendance, error)
	GetByCourse(courseCode string) ([]Attendance, error)
	GetByLecturer(lecturerID string) ([]Attendance, error)
	Update(record *Attendance) error
	Delete(id int) error
}
