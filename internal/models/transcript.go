package models

import (
	"time"
)

// TranscriptEntry represents one course on a student's transcript.
type TranscriptEntry struct {
	CourseCode  string  `json:"course_code"`
	CourseTitle string  `json:"course_title"`
	CreditUnits int     `json:"credit_units"`
	Score       float64 `json:"score"`
	Grade       string  `json:"grade"`
	GradePoint  float64 `json:"grade_point"`
	Session     string  `json:"session"`
	Semester    string  `json:"semester"`
	Level       int     `json:"level"`
}

// SemesterTranscript represents a semester summary.
type SemesterTranscript struct {
	Session       string            `json:"session"`
	Semester      string            `json:"semester"`
	Level         int               `json:"level"`
	Courses       []TranscriptEntry `json:"courses"`
	TotalUnits    int               `json:"total_units"`
	QualityPoints float64           `json:"quality_points"`
	GPA           float64           `json:"gpa"`
}

// Transcript represents a student's academic transcript.
type Transcript struct {
	StudentID    string `json:"student_id"`
	StudentName  string `json:"student_name"`
	MatricNumber string `json:"matric_number"`

	DepartmentName string `json:"department_name"`
	FacultyName    string `json:"faculty_name"`

	CGPA           float64 `json:"cgpa"`
	Classification string  `json:"classification"`

	Results []Result `json:"results"`

	GeneratedAt time.Time `json:"generated_at"`
}
