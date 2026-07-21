package models

import "time"

// Timetable represents a scheduled class for a course.
type Timetable struct {
	ID           int    `json:"id" db:"id"`
	CourseCode   string `json:"course_code" db:"course_code"`
	LecturerID   string `json:"lecturer_id" db:"lecturer_id"`
	DepartmentID int    `json:"department_id" db:"department_id"`

	Session  string `json:"session" db:"session"`
	Semester string `json:"semester" db:"semester"`
	Level    int    `json:"level" db:"level"`

	DayOfWeek string `json:"day_of_week" db:"day_of_week"`

	StartTime string `json:"start_time" db:"start_time"`
	EndTime   string `json:"end_time" db:"end_time"`

	Venue string `json:"venue" db:"venue"`

	IsActive bool `json:"is_active" db:"is_active"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// TimetableRepository defines timetable persistence operations.
type TimetableRepository interface {
	Create(entry *Timetable) error
	FindByID(id int) (*Timetable, error)
	GetAll() ([]Timetable, error)
	GetByDepartment(departmentID int) ([]Timetable, error)
	GetByLecturer(lecturerID string) ([]Timetable, error)
	GetByCourse(courseCode string) ([]Timetable, error)
	GetByLevel(level int) ([]Timetable, error)
	Update(entry *Timetable) error
	Delete(id int) error
}
