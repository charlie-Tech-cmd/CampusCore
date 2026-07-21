package models

import "time"

// Faculty represents a university faculty.
type Faculty struct {
	ID          int       `json:"id" db:"id"`
	Code        string    `json:"code" db:"code"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// FacultyRepository defines persistence operations.
type FacultyRepository interface {
	Create(*Faculty) error
	FindByID(int) (*Faculty, error)
	FindByCode(string) (*Faculty, error)
	List() ([]Faculty, error)
	Update(*Faculty) error
	Delete(int) error
}
