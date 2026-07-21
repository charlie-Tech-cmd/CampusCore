package models

import "time"

// Transcript represents a student's academic transcript.
type Transcript struct {
	ID           int       `json:"id" db:"id"`
	StudentID    string    `json:"student_id" db:"student_id"`
	Session      string    `json:"session" db:"session"`
	Semester     string    `json:"semester" db:"semester"`
	CGPA         float64   `json:"cgpa" db:"cgpa"`
	TotalCredits int       `json:"total_credits" db:"total_credits"`
	Remarks      string    `json:"remarks" db:"remarks"`
	GeneratedAt  time.Time `json:"generated_at" db:"generated_at"`
	GeneratedBy  string    `json:"generated_by" db:"generated_by"`
}

// TranscriptRepository defines transcript persistence operations.
type TranscriptRepository interface {
	Create(transcript *Transcript) error
	FindByID(id int) (*Transcript, error)
	FindByStudent(studentID string) ([]Transcript, error)
	List() ([]Transcript, error)
	Update(transcript *Transcript) error
	Delete(id int) error
}
