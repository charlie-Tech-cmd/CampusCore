package models

import "time"

// UserRole represents the explicit authorization level within CampusCore
type UserRole string

const (
	RoleStudent   UserRole = "student"
	RoleLecturer  UserRole = "lecturer"
	RoleAdmin     UserRole = "admin"
	RoleBursar    UserRole = "bursar"
	RoleLibrarian UserRole = "librarian"
)

// User represents the universal core registry record for any account
type User struct {
	ID           string    `json:"id" db:"id"` // Matric No, Staff ID, or Admin UUID
	Surname      string    `json:"surname" db:"surname"`
	FirstName    string    `json:"first_name" db:"first_name"`
	MiddleName   string    `json:"middle_name,omitempty" db:"middle_name"`
	Email        string    `json:"email" db:"email"`
	Phone        string    `json:"phone" db:"phone"`
	PasswordHash string    `json:"-" db:"password_hash"` // Hidden from JSON serialization for security
	Role         UserRole  `json:"role" db:"role"`
	DepartmentID int       `json:"department_id,omitempty" db:"department_id"`
	Level        int       `json:"level" db:"level"` // e.g., 100, 200, 300, 400, 500
	LastLogin    time.Time `json:"last_login,omitempty" db:"last_login"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// UserOnboarding acts as the data transfer object when capturing sign-up forms
type UserOnboarding struct {
	ID              string   `json:"id"`
	Surname         string   `json:"surname"`
	FirstName       string   `json:"first_name"`
	MiddleName      string   `json:"middle_name"`
	Email           string   `json:"email"`
	Phone           string   `json:"phone"`
	Password        string   `json:"password"`
	ConfirmPassword string   `json:"confirm_password"`
	Role            UserRole `json:"role"`
	DepartmentID    int      `json:"department_id"`
	Level           int      `json:"level"`
}

// UserRepository defines the required contract for user data manipulations
type UserRepository interface {
	Create(user *User) error
	FindByID(id string) (*User, error)
	FindByEmail(email string) (*User, error)
	UpdateLastLogin(id string) error
}
