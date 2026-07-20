package models

import "time"

// Department represents an academic department.
type Department struct {
	ID          int       `json:"id" db:"id"`
	Code        string    `json:"code" db:"code"`
	Name        string    `json:"name" db:"name"`
	FacultyID   int       `json:"faculty_id" db:"faculty_id"`
	Description string    `json:"description" db:"description"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// DepartmentRepository defines department persistence operations.
type DepartmentRepository interface {
	Create(department *Department) error
	FindByID(id int) (*Department, error)
	FindByCode(code string) (*Department, error)
	List() ([]Department, error)
	Update(department *Department) error
	Delete(id int) error
}
