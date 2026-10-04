package repository

import (
	"campuscore/internal/models"
	"database/sql"
	"errors"
	"fmt"
)

// PostgresUserRepository implements the models.UserRepository interface contract
type PostgresUserRepository struct {
	db *sql.DB
}

// NewPostgresUserRepository instantiates our user account data access handler
func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

// Create inserts a brand new account registry item into the users database table
func (r *PostgresUserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (id, surname, first_name, middle_name, email, phone, password_hash, role, department_id, level)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`

	_, err := r.db.Exec(query,
		user.ID,
		user.Surname,
		user.FirstName,
		user.MiddleName,
		user.Email,
		user.Phone,
		user.PasswordHash,
		user.Role,
		user.DepartmentID,
		user.Level,
	)
	if err != nil {
		return fmt.Errorf("failed to write record registry row to postgres: %w", err)
	}
	return nil
}

// FindByID performs a fast indexed B-Tree evaluation to match an account by ID code
func (r *PostgresUserRepository) FindByID(id string) (*models.User, error) {
	query := `
		SELECT id, surname, first_name, middle_name, email, phone, password_hash, role, department_id, level, last_login, created_at
		FROM users 
		WHERE id = $1 
		LIMIT 1;`

	row := r.db.QueryRow(query, id)

	var u models.User
	var middleNameNull sql.NullString
	var phoneNull sql.NullString
	var deptNull sql.NullInt32
	var levelNull sql.NullInt32
	var lastLoginNull sql.NullTime

	err := row.Scan(
		&u.ID,
		&u.Surname,
		&u.FirstName,
		&middleNameNull,
		&u.Email,
		&phoneNull,
		&u.PasswordHash,
		&u.Role,
		&deptNull,
		&levelNull,
		&lastLoginNull,
		&u.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user account with ID %s not found: %w", id, err)
		}
		return nil, fmt.Errorf("failed to process find by id query: %w", err)
	}

	if middleNameNull.Valid {
		u.MiddleName = middleNameNull.String
	}
	if phoneNull.Valid {
		u.Phone = phoneNull.String
	}
	if deptNull.Valid {
		u.DepartmentID = int(deptNull.Int32)
	}
	if levelNull.Valid {
		u.Level = int(levelNull.Int32)
	}
	if lastLoginNull.Valid {
		u.LastLogin = lastLoginNull.Time
	}

	return &u, nil
}

// FindByEmail searches for an account record matching an explicit login email address
func (r *PostgresUserRepository) FindByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, surname, first_name, middle_name, email, phone, password_hash, role, department_id, level, last_login, created_at
		FROM users 
		WHERE email = $1 
		LIMIT 1;`

	row := r.db.QueryRow(query, email)

	var u models.User
	var middleNameNull sql.NullString
	var phoneNull sql.NullString
	var deptNull sql.NullInt32
	var levelNull sql.NullInt32
	var lastLoginNull sql.NullTime

	err := row.Scan(
		&u.ID,
		&u.Surname,
		&u.FirstName,
		&middleNameNull,
		&u.Email,
		&phoneNull,
		&u.PasswordHash,
		&u.Role,
		&deptNull,
		&levelNull,
		&lastLoginNull,
		&u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user account with email %s not found: %w", email, err)
		}
		return nil, fmt.Errorf("failed to process find by email query: %w", err)
	}

	if middleNameNull.Valid {
		u.MiddleName = middleNameNull.String
	}
	if phoneNull.Valid {
		u.Phone = phoneNull.String
	}
	if deptNull.Valid {
		u.DepartmentID = int(deptNull.Int32)
	}
	if levelNull.Valid {
		u.Level = int(levelNull.Int32)
	}
	if lastLoginNull.Valid {
		u.LastLogin = lastLoginNull.Time
	}

	return &u, nil
}

// UpdateLastLogin patches the access timestamp row when an active authorization executes successfully
func (r *PostgresUserRepository) UpdateLastLogin(id string) error {
	query := `UPDATE users SET last_login = CURRENT_TIMESTAMP WHERE id = $1;`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to update account tracking log status: %w", err)
	}
	return nil
}

// GetProfile retrieves a student's profile using their unique ID.
func (r *PostgresUserRepository) GetProfile(id string) (*models.User, error) {
	query := `
		SELECT
			id,
			surname,
			first_name,
			middle_name,
			email,
			phone,
			password_hash,
			role,
			department_id,
			level,
			last_login,
			created_at
		FROM users
		WHERE id = $1
		LIMIT 1;
	`

	row := r.db.QueryRow(query, id)

	var u models.User
	var middleNameNull sql.NullString
	var phoneNull sql.NullString
	var deptNull sql.NullInt32
	var levelNull sql.NullInt32
	var lastLoginNull sql.NullTime

	err := row.Scan(
		&u.ID,
		&u.Surname,
		&u.FirstName,
		&middleNameNull,
		&u.Email,
		&phoneNull,
		&u.PasswordHash,
		&u.Role,
		&deptNull,
		&levelNull,
		&lastLoginNull,
		&u.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("profile not found for user %s: %w", id, err)
		}
		return nil, fmt.Errorf("failed to retrieve user profile: %w", err)
	}

	if middleNameNull.Valid {
		u.MiddleName = middleNameNull.String
	}
	if phoneNull.Valid {
		u.Phone = phoneNull.String
	}
	if deptNull.Valid {
		u.DepartmentID = int(deptNull.Int32)
	}
	if levelNull.Valid {
		u.Level = int(levelNull.Int32)
	}
	if lastLoginNull.Valid {
		u.LastLogin = lastLoginNull.Time
	}

	return &u, nil
}

// UpdateProfile updates editable profile information.
func (r *PostgresUserRepository) UpdateProfile(user *models.User) error {
	query := `
		UPDATE users
		SET
			surname = $2,
			first_name = $3,
			middle_name = $4,
			email = $5,
			phone = $6,
			department_id = $7,
			level = $8
		WHERE id = $1;
	`

	result, err := r.db.Exec(
		query,
		user.ID,
		user.Surname,
		user.FirstName,
		user.MiddleName,
		user.Email,
		user.Phone,
		user.DepartmentID,
		user.Level,
	)

	if err != nil {
		return fmt.Errorf("failed to update user profile: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to determine affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user profile not found")
	}

	return nil
}
