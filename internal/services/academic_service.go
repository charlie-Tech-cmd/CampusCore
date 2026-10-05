package services

import (
	"database/sql"
	"errors"
	"fmt"

	"campuscore/internal/models"

	"github.com/lib/pq"
)

// AcademicService manages business logic validation for enrollment and grading
type AcademicService struct {
	db *sql.DB
}

// NewAcademicService instantiates our business logic controller with a database connection handle
func NewAcademicService(db *sql.DB) *AcademicService {
	return &AcademicService{db: db}
}

// RegisterCourse validates and processes a student course enrollment request safely.
func (s *AcademicService) RegisterCourse(
	studentID string,
	courseCode string,
	session string,
	semester string,
) error {
	// Start a transaction so all validation and registration operations
	// succeed or fail together.
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to initialize transaction context: %w", err)
	}
	defer tx.Rollback()

	// 1. Fetch course details and lock the course row.
	var creditUnits, level, maxCapacity, currentEnrolled int
	var courseSemester sql.NullString
	courseQuery := `
    SELECT
        credit_units,
        level,
        max_capacity,
        current_enrolled,
        semester
    FROM courses
    WHERE code = $1
    FOR UPDATE;
`

	err = tx.QueryRow(courseQuery, courseCode).Scan(
		&creditUnits,
		&level,
		&maxCapacity,
		&currentEnrolled,
		&courseSemester,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New(
				"academic rule violation: requested course code does not exist in curriculum record",
			)
		}

		return err
	}

	if !courseSemester.Valid || courseSemester.String != semester {
		return fmt.Errorf(
			"academic rule violation: course %s is not offered in the %s semester",
			courseCode,
			semester,
		)
	}

	// Check course capacity.
	if currentEnrolled >= maxCapacity {
		return fmt.Errorf(
			"enrollment capacity exceeded: course %s has reached its maximum limit of %d students",
			courseCode,
			maxCapacity,
		)
	}

	// 2. Verify student and determine student level.
	var studentLevel int

	studentQuery := `
		SELECT level
		FROM users
		WHERE id = $1
		  AND role = 'student';
	`

	err = tx.QueryRow(studentQuery, studentID).Scan(&studentLevel)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New(
				"access denied: student account configuration missing or invalid",
			)
		}

		return err
	}

	// Check level restriction.
	if studentLevel < level {
		return fmt.Errorf(
			"academic rule violation: course %s is reserved for %d-level students (current tier: %d-level)",
			courseCode,
			level,
			studentLevel,
		)
	}

	// 3. Check student's current credit load.
	var currentTotalUnits int

	loadQuery := `
		SELECT COALESCE(SUM(c.credit_units), 0)
		FROM student_courses sc
		JOIN courses c ON sc.course_code = c.code
		JOIN academic_sessions a ON sc.session_id = a.id
		WHERE sc.student_id = $1
		  AND a.session_name = $2
		  AND a.semester = $3
		  AND sc.status = 'approved';
	`

	err = tx.QueryRow(
		loadQuery,
		studentID,
		session,
		semester,
	).Scan(&currentTotalUnits)

	if err != nil {
		return err
	}

	const maxCreditUnits = 24

	if currentTotalUnits+creditUnits > maxCreditUnits {
		return fmt.Errorf(
			"credit load limit exceeded: adding this course (%d units) pushes total load to %d units (maximum limit: %d units)",
			creditUnits,
			currentTotalUnits+creditUnits,
			maxCreditUnits,
		)
	}

	// 4. Check course prerequisites.
	prereqQuery := `
        SELECT prerequisite_code
        FROM course_prerequisites
        WHERE course_code = $1;
`

	rows, err := tx.Query(prereqQuery, courseCode)
	if err != nil {
		return err
	}
	var prerequisites []string

	for rows.Next() {
		var prereqCode string

		if err := rows.Scan(&prereqCode); err != nil {
			rows.Close()
			return fmt.Errorf("failed to read prerequisite: %w", err)
		}

		prerequisites = append(prerequisites, prereqCode)
	}

	if err := rows.Err(); err != nil {
		rows.Close()
		return fmt.Errorf("failed to read prerequisites: %w", err)
	}

	rows.Close()

	// Close prerequisites before checking results.
	checkPassedQuery := `
        SELECT EXISTS(
                SELECT 1
                FROM results
                WHERE student_id = $1
                  AND course_code = $2
                  AND score >= 40.00
        );
`

	for _, prereqCode := range prerequisites {
		var passed bool

		err := tx.QueryRow(
			checkPassedQuery,
			studentID,
			prereqCode,
		).Scan(&passed)

		if err != nil {
			return fmt.Errorf(
				"failed to verify prerequisite %s: %w",
				prereqCode,
				err,
			)
		}

		if !passed {
			return fmt.Errorf(
				"prerequisite requirement failed: you must pass course %s before attempting %s",
				prereqCode,
				courseCode,
			)
		}
	}
	// 5. Register the course.
	insertQuery := `
		INSERT INTO student_courses (
			student_id,
			course_code,
			session_id,
			status
		)
		SELECT
			$1,
			$2,
			id,
			'approved'
		FROM academic_sessions
		WHERE session_name = $3
		  AND semester = $4
		  AND registration_open = true
		  AND is_active = true;
	`

	result, err := tx.Exec(
		insertQuery,
		studentID,
		courseCode,
		session,
		semester,
	)

	if err != nil {
		// PostgreSQL unique violation.
		// Constraint:
		// student_courses_student_id_course_code_session_id_key
		var pqErr *pq.Error

		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return fmt.Errorf(
				"course registration conflict: you are already registered for %s for the %s session",
				courseCode,
				session,
			)
		}

		return fmt.Errorf(
			"failed to complete course registration: %w",
			err,
		)
	}

	// Make sure a valid, open academic session was found.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(
			"failed to verify course registration: %w",
			err,
		)
	}

	if rowsAffected == 0 {
		return fmt.Errorf(
			"academic session %s / %s is not active or registration is closed",
			session,
			semester,
		)
	}

	// 6. Increment course enrollment count.
	updateCapQuery := `
		UPDATE courses
		SET current_enrolled = current_enrolled + 1
		WHERE code = $1;
	`

	_, err = tx.Exec(updateCapQuery, courseCode)
	if err != nil {
		return fmt.Errorf(
			"failed to update course enrollment count: %w",
			err,
		)
	}

	// 7. Commit the complete registration.
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
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
