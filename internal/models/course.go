package models

// Course represents an academic course offered by the institution.
type Course struct {
	Code            string `json:"code" db:"code"`
	Title           string `json:"title" db:"title"`
	Description     string `json:"description,omitempty" db:"description"`
	CreditUnits     int    `json:"credit_units" db:"credit_units"`
	DepartmentID    int    `json:"department_id" db:"department_id"`
	Level           int    `json:"level" db:"level"`
	Semester        string `json:"semester" db:"semester"`
	MaxCapacity     int    `json:"max_capacity" db:"max_capacity"`
	CurrentEnrolled int    `json:"current_enrolled" db:"current_enrolled"`
	IsActive        bool   `json:"is_active" db:"is_active"`
}

// CourseRepository defines the required contract for course data operations.
type CourseRepository interface {
	Create(course *Course) error
	FindByCode(code string) (*Course, error)
	GetAll() ([]Course, error)
	GetByDepartment(departmentID int) ([]Course, error)
	Update(course *Course) error
	Delete(code string) error
}
