package models

import "time"

// Result represents a student's academic performance in a course.
type Result struct {
	ID          int       `json:"id" db:"id"`
	StudentID   string    `json:"student_id" db:"student_id"`
	CourseCode  string    `json:"course_code" db:"course_code"`
	Session     string    `json:"session" db:"session"`
	Semester    string    `json:"semester" db:"semester"`
	Score       float64   `json:"score" db:"score"`
	Grade       string    `json:"grade" db:"grade"`
	GradePoint  float64   `json:"grade_point" db:"grade_point"`
	CreditUnits int       `json:"credit_units" db:"credit_units"`
	Approved    bool      `json:"approved" db:"approved"`
	ApprovedBy  string    `json:"approved_by,omitempty" db:"approved_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// ResultRepository defines the contract for academic result operations.
type ResultRepository interface {
	Submit(result *Result) error
	FindByStudent(studentID string) ([]Result, error)
	FindByCourse(courseCode string) ([]Result, error)
	Update(result *Result) error
	Delete(id int) error
}
