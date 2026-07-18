package services

import (
	"campuscore/internal/models"
	"database/sql"
	"errors"
	"fmt"
)

// AcademicService manages business logic validation for enrollment and grading
type AcademicService struct {
	db *sql.DB
}

// NewAcademicService instantiates our business logic controller with a database connection handle
func NewAcademicService(db *sql.DB) *AcademicService {
	return &AcademicService{db: db}
}

// RegisterCourse validates and processes a student course enrollment request safely
func (s *AcademicService) RegisterCourse(studentID string, courseCode string, session string, semester string) error {
	// Start an isolated database transaction to ensure atomicity across all dependency checks
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to initialize transaction context: %w", err)
	}
	defer tx.Rollback() // Safe recovery cascade: automatically rolls back adjustments if execution panics or fails

	// 1. Fetch target course structural metrics and verify room availability
	var creditUnits, level, maxCapacity, currentEnrolled int
	courseQuery := `SELECT credit_units, level, max_capacity, current_enrolled FROM courses WHERE code = $1 FOR UPDATE;`
	err = tx.QueryRow(courseQuery, courseCode).Scan(&creditUnits, &level, &maxCapacity, &currentEnrolled)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("academic rule violation: requested course code does not exist in curriculum record")
		}
		return err
	}

	if currentEnrolled >= maxCapacity {
		return fmt.Errorf("enrollment capacity exceeded: course %s has reached its maximum limit of %d students", courseCode, maxCapacity)
	}

	// 2. Fetch target student metadata profile values to evaluate registration bounds
	var studentLevel int
	studentQuery := `SELECT level FROM users WHERE id = $1 AND role = 'student';`
	err = tx.QueryRow(studentQuery, studentID).Scan(&studentLevel)
	if err != nil {
		return errors.New("access denied: student account configuration missing or invalid")
	}

	if studentLevel < level {
		return fmt.Errorf("academic rule violation: course %s is reserved for %d-level students (current tier: %d-level)", courseCode, level, studentLevel)
	}

	// 3. Evaluate total academic load threshold constraints for this semester (Max 24 units allowed)
	var currentTotalUnits int
	loadQuery := `
		SELECT COALESCE(SUM(c.credit_units), 0) 
		FROM student_courses sc 
		JOIN courses c ON sc.course_code = c.code 
		WHERE sc.student_id = $1 AND sc.session = $2 AND sc.semester = $3 AND sc.status = 'approved';`

	err = tx.QueryRow(loadQuery, studentID, session, semester).Scan(&currentTotalUnits)
	if err != nil {
		return err
	}

	if currentTotalUnits+creditUnits > 24 {
		return fmt.Errorf("credit load limit exceeded: adding this course (%d units) pushes total load to %d units (maximum limit: 24 units)", creditUnits, currentTotalUnits+creditUnits)
	}

	// 4. Verify prerequisite requirements (Self-Referencing junction validation)
	prereqQuery := `SELECT prerequisite_code FROM course_prerequisites WHERE course_code = $1;`
	rows, err := tx.Query(prereqQuery, courseCode)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var prereqCode string
		if err := rows.Scan(&prereqCode); err != nil {
			return err
		}

		// Double-check if the student has a passing grade in the required prerequisite course
		var passed bool
		checkPassedQuery := `SELECT EXISTS(SELECT 1 FROM results WHERE student_id = $1 AND course_code = $2 AND score >= 40.00);`
		err = tx.QueryRow(checkPassedQuery, studentID, prereqCode).Scan(&passed)
		if err != nil || !passed {
			return fmt.Errorf("prerequisite requirement failed: you must pass course %s before attempting %s", prereqCode, courseCode)
		}
	}

	// 5. Commit record registry parameters to the database logs
	insertQuery := `
		INSERT INTO student_courses (student_id, course_code, session, semester) 
		VALUES ($1, $2, $3, $4);`
	_, err = tx.Exec(insertQuery, studentID, courseCode, session, semester)
	if err != nil {
		return fmt.Errorf("failed to complete course registry insertion: %w", err)
	}

	// Increment current enrollment tracking counters dynamically
	updateCapQuery := `UPDATE courses SET current_enrolled = current_enrolled + 1 WHERE code = $1;`
	_, err = tx.Exec(updateCapQuery, courseCode)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// CalculateGradeMetrics processes a raw assessment point to output structural grading results
func (s *AcademicService) CalculateGradeMetrics(score float64) (string, float64) {
	switch {
	case score >= 70:
		return "A", 5.0
	case score >= 60:
		return "B", 4.0
	case score >= 50:
		return "C", 3.0
	case score >= 45:
		return "D", 2.0
	case score >= 40:
		return "E", 1.0
	default:
		return "F", 0.0
	}
}

// GetStudentProfile retrieves a student's profile information.
func (s *AcademicService) GetStudentProfile(studentID string) (*models.User, error) {
	query := `
		SELECT
			id,
			surname,
			first_name,
			middle_name,
			email,
			phone,
			role,
			department_id,
			level,
			last_login,
			created_at
		FROM users
		WHERE id = $1
		  AND role = 'student'
		LIMIT 1;
	`

	var profile models.User
	var dept sql.NullInt32
	var lastLogin sql.NullTime

	err := s.db.QueryRow(query, studentID).Scan(
		&profile.ID,
		&profile.Surname,
		&profile.FirstName,
		&profile.MiddleName,
		&profile.Email,
		&profile.Phone,
		&profile.Role,
		&dept,
		&profile.Level,
		&lastLogin,
		&profile.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("student profile not found")
		}
		return nil, err
	}

	if dept.Valid {
		profile.DepartmentID = int(dept.Int32)
	}

	if lastLogin.Valid {
		profile.LastLogin = lastLogin.Time
	}

	return &profile, nil
}

// UpdateStudentProfile updates editable student profile fields.
func (s *AcademicService) UpdateStudentProfile(profile *models.User) error {
	query := `
		UPDATE users
		SET
			surname = $2,
			first_name = $3,
			middle_name = $4,
			email = $5,
			phone = $6
		WHERE id = $1
		  AND role = 'student';
	`

	result, err := s.db.Exec(
		query,
		profile.ID,
		profile.Surname,
		profile.FirstName,
		profile.MiddleName,
		profile.Email,
		profile.Phone,
	)

	if err != nil {
		return fmt.Errorf("failed to update student profile: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("student profile not found")
	}

	return nil
}
