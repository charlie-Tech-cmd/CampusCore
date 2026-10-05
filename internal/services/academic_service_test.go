package services

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
)

func TestNewAcademicService(t *testing.T) {
	service := NewAcademicService(nil)

	if service == nil {
		t.Fatal("expected AcademicService, got nil")
	}

	if service.db != nil {
		t.Fatal("expected nil database")
	}
}

func TestCalculateGradeMetrics(t *testing.T) {
	service := NewAcademicService(nil)

	tests := []struct {
		name  string
		score float64
		grade string
		gpa   float64
	}{
		{"A", 85, "A", 5.0},
		{"B", 65, "B", 4.0},
		{"C", 55, "C", 3.0},
		{"D", 47, "D", 2.0},
		{"E", 40, "E", 1.0},
		{"F", 39, "F", 0.0},
		{"Perfect", 100, "A", 5.0},
		{"Zero", 0, "F", 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grade, gpa := service.CalculateGradeMetrics(tt.score)

			if grade != tt.grade {
				t.Fatalf("expected grade %q, got %q", tt.grade, grade)
			}

			if gpa != tt.gpa {
				t.Fatalf("expected GPA %.1f, got %.1f", tt.gpa, gpa)
			}
		})
	}
}

func TestRegisterCourse_BeginTransactionError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	// Close the database so Begin() will fail.
	db.Close()

	service := NewAcademicService(db)

	err = service.RegisterCourse(
		"STU001",
		"CSC401",
		"2025/2026",
		"First",
	)

	if err == nil {
		t.Fatal("expected transaction initialization error")
	}
}

func TestRegisterCourse_CourseNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	service := NewAcademicService(db)

	mock.ExpectBegin()

	mock.ExpectQuery(`SELECT\s+credit_units,\s+level,\s+max_capacity,\s+current_enrolled,\s+semester\s+FROM courses`).WithArgs("CSC401").
		WillReturnError(sql.ErrNoRows)

	mock.ExpectRollback()

	err = service.RegisterCourse(
		"STU001",
		"CSC401",
		"2025/2026",
		"First",
	)

	if err == nil {
		t.Fatal("expected error")
	}

	expected := "academic rule violation: requested course code does not exist in curriculum record"
	if err.Error() != expected {
		t.Fatalf("expected %q, got %q", expected, err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestRegisterCourse_CourseFull(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	service := NewAcademicService(db)

	mock.ExpectBegin()

	rows := sqlmock.NewRows([]string{
		"credit_units", "level", "max_capacity", "current_enrolled", "semester"}).AddRow(3, 400, 100, 100, "First")

	mock.ExpectQuery(`SELECT\s+credit_units,\s+level,\s+max_capacity,\s+current_enrolled,\s+semester\s+FROM courses`).WithArgs("CSC401").
		WillReturnRows(rows)

	mock.ExpectRollback()

	err = service.RegisterCourse(
		"STU001",
		"CSC401",
		"2025/2026",
		"First",
	)

	if err == nil {
		t.Fatal("expected course capacity error")
	}

	expected := "enrollment capacity exceeded: course CSC401 has reached its maximum limit of 100 students"
	if err.Error() != expected {
		t.Fatalf("expected %q, got %q", expected, err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestRegisterCourse_InvalidStudent(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	service := NewAcademicService(db)

	mock.ExpectBegin()

	courseRows := sqlmock.NewRows([]string{
		"credit_units", "level", "max_capacity", "current_enrolled", "semester"}).AddRow(3, 400, 100, 50, "First")

	mock.ExpectQuery(`SELECT\s+credit_units,\s+level,\s+max_capacity,\s+current_enrolled,\s+semester\s+FROM courses`).WithArgs("CSC401").
		WillReturnRows(courseRows)

	mock.ExpectQuery(
		"SELECT level FROM users WHERE id = \\$1 AND role = 'student'",
	).
		WithArgs("STU001").
		WillReturnError(sql.ErrNoRows)

	mock.ExpectRollback()

	err = service.RegisterCourse(
		"STU001",
		"CSC401",
		"2025/2026",
		"First",
	)

	if err == nil {
		t.Fatal("expected error")
	}

	expected := "access denied: student account configuration missing or invalid"
	if err.Error() != expected {
		t.Fatalf("expected %q, got %q", expected, err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestRegisterCourse_LevelRestriction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	service := NewAcademicService(db)

	mock.ExpectBegin()

	courseRows := sqlmock.NewRows([]string{
		"credit_units", "level", "max_capacity", "current_enrolled", "semester"}).AddRow(3, 400, 100, 50, "First")

	mock.ExpectQuery(`SELECT\s+credit_units,\s+level,\s+max_capacity,\s+current_enrolled,\s+semester\s+FROM courses`).
		WithArgs("CSC401").
		WillReturnRows(courseRows)

	studentRows := sqlmock.NewRows([]string{"level"}).
		AddRow(300)

	mock.ExpectQuery(
		"SELECT level FROM users WHERE id = \\$1 AND role = 'student'",
	).
		WithArgs("STU001").
		WillReturnRows(studentRows)

	mock.ExpectRollback()

	err = service.RegisterCourse(
		"STU001",
		"CSC401",
		"2025/2026",
		"First",
	)

	if err == nil {
		t.Fatal("expected level restriction error")
	}

	expected := "academic rule violation: course CSC401 is reserved for 400-level students (current tier: 300-level)"
	if err.Error() != expected {
		t.Fatalf("expected %q, got %q", expected, err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestRegisterCourse_CreditLimitExceeded(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	service := NewAcademicService(db)

	mock.ExpectBegin()

	// Course exists
	courseRows := sqlmock.NewRows([]string{
		"credit_units", "level", "max_capacity", "current_enrolled", "semester"}).AddRow(3, 400, 100, 50, "First")

	mock.ExpectQuery(`SELECT\s+credit_units,\s+level,\s+max_capacity,\s+current_enrolled,\s+semester\s+FROM courses`).
		WithArgs("CSC401").
		WillReturnRows(courseRows)

	// Student exists
	studentRows := sqlmock.NewRows([]string{"level"}).
		AddRow(400)

	mock.ExpectQuery(
		"SELECT level FROM users WHERE id = \\$1 AND role = 'student'",
	).
		WithArgs("STU001").
		WillReturnRows(studentRows)

	// Student already has 22 units
	loadRows := sqlmock.NewRows([]string{"coalesce"}).
		AddRow(22)

	mock.ExpectQuery(
		"SELECT COALESCE\\(SUM\\(c.credit_units\\), 0\\)",
	).
		WithArgs("STU001", "2025/2026", "First").
		WillReturnRows(loadRows)

	mock.ExpectRollback()

	err = service.RegisterCourse(
		"STU001",
		"CSC401",
		"2025/2026",
		"First",
	)

	if err == nil {
		t.Fatal("expected credit limit error")
	}

	expected := "credit load limit exceeded: adding this course (3 units) pushes total load to 25 units (maximum limit: 24 units)"
	if err.Error() != expected {
		t.Fatalf("expected %q, got %q", expected, err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestRegisterCourse_PrerequisiteFailed(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	service := NewAcademicService(db)

	mock.ExpectBegin()

	// Course lookup
	courseRows := sqlmock.NewRows([]string{
		"credit_units", "level", "max_capacity", "current_enrolled", "semester"}).AddRow(3, 400, 100, 50, "First")

	mock.ExpectQuery(`SELECT\s+credit_units,\s+level,\s+max_capacity,\s+current_enrolled,\s+semester\s+FROM courses`).
		WithArgs("CSC401").
		WillReturnRows(courseRows)

	// Student lookup
	studentRows := sqlmock.NewRows([]string{"level"}).
		AddRow(400)

	mock.ExpectQuery(
		"SELECT level FROM users WHERE id = \\$1 AND role = 'student'",
	).
		WithArgs("STU001").
		WillReturnRows(studentRows)

	// Current load
	loadRows := sqlmock.NewRows([]string{"coalesce"}).
		AddRow(18)

	mock.ExpectQuery(
		"SELECT COALESCE\\(SUM\\(c.credit_units\\), 0\\)",
	).
		WithArgs("STU001", "2025/2026", "First").
		WillReturnRows(loadRows)

	// One prerequisite exists
	prereqRows := sqlmock.NewRows([]string{"prerequisite_code"}).
		AddRow("CSC301")

	mock.ExpectQuery(
		"SELECT prerequisite_code FROM course_prerequisites",
	).
		WithArgs("CSC401").
		WillReturnRows(prereqRows)

	// Student did NOT pass prerequisite
	passedRows := sqlmock.NewRows([]string{"exists"}).
		AddRow(false)

	mock.ExpectQuery(
		"SELECT EXISTS",
	).
		WithArgs("STU001", "CSC301").
		WillReturnRows(passedRows)

	mock.ExpectRollback()

	err = service.RegisterCourse(
		"STU001",
		"CSC401",
		"2025/2026",
		"First",
	)

	if err == nil {
		t.Fatal("expected prerequisite error")
	}

	expected := "prerequisite requirement failed: you must pass course CSC301 before attempting CSC401"

	if err.Error() != expected {
		t.Fatalf("expected %q, got %q", expected, err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestRegisterCourse_InsertError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	service := NewAcademicService(db)

	mock.ExpectBegin()

	// Course lookup
	courseRows := sqlmock.NewRows([]string{
		"credit_units", "level", "max_capacity", "current_enrolled", "semester"}).AddRow(3, 400, 100, 50, "First")

	mock.ExpectQuery(`SELECT\s+credit_units,\s+level,\s+max_capacity,\s+current_enrolled,\s+semester\s+FROM courses`).
		WithArgs("CSC401").
		WillReturnRows(courseRows)

	// Student lookup
	studentRows := sqlmock.NewRows([]string{"level"}).
		AddRow(400)

	mock.ExpectQuery(
		"SELECT level FROM users WHERE id = \\$1 AND role = 'student'",
	).
		WithArgs("STU001").
		WillReturnRows(studentRows)

	// Credit load
	loadRows := sqlmock.NewRows([]string{"coalesce"}).
		AddRow(18)

	mock.ExpectQuery(
		"SELECT COALESCE\\(SUM\\(c.credit_units\\), 0\\)",
	).
		WithArgs("STU001", "2025/2026", "First").
		WillReturnRows(loadRows)

	// No prerequisites
	prereqRows := sqlmock.NewRows([]string{"prerequisite_code"})

	mock.ExpectQuery(
		"SELECT prerequisite_code FROM course_prerequisites",
	).
		WithArgs("CSC401").
		WillReturnRows(prereqRows)

	// Insert fails
	mock.ExpectExec(
		"INSERT INTO student_courses",
	).
		WithArgs(
			"STU001",
			"CSC401",
			"2025/2026",
			"First",
		).
		WillReturnError(errors.New("insert failed"))

	mock.ExpectRollback()

	err = service.RegisterCourse(
		"STU001",
		"CSC401",
		"2025/2026",
		"First",
	)

	if err == nil {
		t.Fatal("expected insert error")
	}

	expected := "failed to complete course registration: insert failed"

	if err.Error() != expected {
		t.Fatalf("expected %q, got %q", expected, err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestRegisterCourse_UpdateEnrollmentError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	service := NewAcademicService(db)

	mock.ExpectBegin()

	// Course lookup
	courseRows := sqlmock.NewRows([]string{
		"credit_units", "level", "max_capacity", "current_enrolled", "semester"}).AddRow(3, 400, 100, 50, "First")

	mock.ExpectQuery(`SELECT\s+credit_units,\s+level,\s+max_capacity,\s+current_enrolled,\s+semester\s+FROM courses`).
		WithArgs("CSC401").
		WillReturnRows(courseRows)

	// Student lookup
	studentRows := sqlmock.NewRows([]string{"level"}).
		AddRow(400)

	mock.ExpectQuery(
		"SELECT level FROM users WHERE id = \\$1 AND role = 'student'",
	).
		WithArgs("STU001").
		WillReturnRows(studentRows)

	// Current credit load
	loadRows := sqlmock.NewRows([]string{"coalesce"}).
		AddRow(18)

	mock.ExpectQuery(
		"SELECT COALESCE\\(SUM\\(c.credit_units\\), 0\\)",
	).
		WithArgs("STU001", "2025/2026", "First").
		WillReturnRows(loadRows)

	// No prerequisites
	prereqRows := sqlmock.NewRows([]string{"prerequisite_code"})

	mock.ExpectQuery(
		"SELECT prerequisite_code FROM course_prerequisites",
	).
		WithArgs("CSC401").
		WillReturnRows(prereqRows)

	// Registration insert succeeds
	mock.ExpectExec(
		"INSERT INTO student_courses",
	).
		WithArgs(
			"STU001",
			"CSC401",
			"2025/2026",
			"First",
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Enrollment update fails
	mock.ExpectExec(
		"UPDATE courses SET current_enrolled = current_enrolled \\+ 1 WHERE code = \\$1",
	).
		WithArgs("CSC401").
		WillReturnError(errors.New("update failed"))

	mock.ExpectRollback()

	err = service.RegisterCourse(
		"STU001",
		"CSC401",
		"2025/2026",
		"First",
	)

	if err == nil {
		t.Fatal("expected update error")
	}

	expected := "failed to update course enrollment count: update failed"

	if err.Error() != expected {
		t.Fatalf("expected %q, got %q", expected, err.Error())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestRegisterCourse_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	service := NewAcademicService(db)

	mock.ExpectBegin()

	// Course lookup
	courseRows := sqlmock.NewRows([]string{
		"credit_units", "level", "max_capacity", "current_enrolled", "semester"}).AddRow(3, 400, 100, 50, "First")

	mock.ExpectQuery(
		`SELECT\s+credit_units,\s+level,\s+max_capacity,\s+current_enrolled,\s+semester\s+FROM courses`,
	).
		WithArgs("CSC401").
		WillReturnRows(courseRows)
	// Student lookup
	studentRows := sqlmock.NewRows([]string{"level"}).
		AddRow(400)

	mock.ExpectQuery(
		"SELECT level FROM users WHERE id = \\$1 AND role = 'student'",
	).
		WithArgs("STU001").
		WillReturnRows(studentRows)

	// Current credit load
	loadRows := sqlmock.NewRows([]string{"coalesce"}).
		AddRow(18)

	mock.ExpectQuery(
		"SELECT COALESCE\\(SUM\\(c.credit_units\\), 0\\)",
	).
		WithArgs("STU001", "2025/2026", "First").
		WillReturnRows(loadRows)

	// No prerequisites
	prereqRows := sqlmock.NewRows([]string{"prerequisite_code"})

	mock.ExpectQuery(
		"SELECT prerequisite_code FROM course_prerequisites",
	).
		WithArgs("CSC401").
		WillReturnRows(prereqRows)

	// Insert succeeds
	mock.ExpectExec(
		"INSERT INTO student_courses",
	).
		WithArgs(
			"STU001",
			"CSC401",
			"2025/2026",
			"First",
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Update enrollment succeeds
	mock.ExpectExec(
		"UPDATE courses SET current_enrolled = current_enrolled \\+ 1 WHERE code = \\$1",
	).
		WithArgs("CSC401").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err = service.RegisterCourse(
		"STU001",
		"CSC401",
		"2025/2026",
		"First",
	)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestRegisterCourse_CommitError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	service := NewAcademicService(db)

	mock.ExpectBegin()

	// Course lookup
	courseRows := sqlmock.NewRows([]string{
		"credit_units", "level", "max_capacity", "current_enrolled", "semester"}).AddRow(3, 400, 100, 50, "First")

	mock.ExpectQuery(`SELECT\s+credit_units,\s+level,\s+max_capacity,\s+current_enrolled,\s+semester\s+FROM courses`).WithArgs("CSC401").
		WillReturnRows(courseRows)

	// Student lookup
	studentRows := sqlmock.NewRows([]string{"level"}).
		AddRow(400)

	mock.ExpectQuery(
		"SELECT level FROM users WHERE id = \\$1 AND role = 'student'",
	).
		WithArgs("STU001").
		WillReturnRows(studentRows)

	// Current load
	loadRows := sqlmock.NewRows([]string{"coalesce"}).
		AddRow(18)

	mock.ExpectQuery(
		"SELECT COALESCE\\(SUM\\(c.credit_units\\), 0\\)",
	).
		WithArgs("STU001", "2025/2026", "First").
		WillReturnRows(loadRows)

	// No prerequisites
	prereqRows := sqlmock.NewRows([]string{"prerequisite_code"})

	mock.ExpectQuery(
		"SELECT prerequisite_code FROM course_prerequisites",
	).
		WithArgs("CSC401").
		WillReturnRows(prereqRows)

	// Insert succeeds
	mock.ExpectExec(
		"INSERT INTO student_courses",
	).
		WithArgs(
			"STU001",
			"CSC401",
			"2025/2026",
			"First",
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Update succeeds
	mock.ExpectExec(
		"UPDATE courses SET current_enrolled = current_enrolled \\+ 1 WHERE code = \\$1",
	).
		WithArgs("CSC401").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Commit fails
	mock.ExpectCommit().
		WillReturnError(errors.New("commit failed"))

	err = service.RegisterCourse(
		"STU001",
		"CSC401",
		"2025/2026",
		"First",
	)

	if err == nil {
		t.Fatal("expected commit error")
	}

	if err.Error() != "commit failed" {
		t.Fatalf("expected %q, got %q", "commit failed", err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestRegisterCourse_CourseQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	service := NewAcademicService(db)

	mock.ExpectBegin()

	mock.ExpectQuery(`SELECT\s+credit_units,\s+level,\s+max_capacity,\s+current_enrolled,\s+semester\s+FROM courses`).WithArgs("CSC401").
		WillReturnError(errors.New("database unavailable"))

	mock.ExpectRollback()

	err = service.RegisterCourse(
		"STU001",
		"CSC401",
		"2025/2026",
		"First",
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "database unavailable" {
		t.Fatalf("expected database unavailable, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestRegisterCourse_LoadQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	service := NewAcademicService(db)

	mock.ExpectBegin()

	// Course lookup succeeds
	mock.ExpectQuery(`SELECT\s+credit_units,\s+level,\s+max_capacity,\s+current_enrolled,\s+semester\s+FROM courses`).WithArgs("CSC401").
		WillReturnRows(
			sqlmock.NewRows([]string{
				"credit_units", "level", "max_capacity", "current_enrolled", "semester"}).AddRow(3, 400, 100, 10, "First"),
		)

	// Student lookup succeeds
	mock.ExpectQuery("SELECT level FROM users").
		WithArgs("STU001").
		WillReturnRows(
			sqlmock.NewRows([]string{"level"}).AddRow(400),
		)

	// Load query fails
	mock.ExpectQuery("SELECT COALESCE\\(SUM\\(c.credit_units\\), 0\\)").
		WithArgs("STU001", "2025/2026", "First").
		WillReturnError(errors.New("load query failed"))

	mock.ExpectRollback()

	err = service.RegisterCourse(
		"STU001",
		"CSC401",
		"2025/2026",
		"First",
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "load query failed" {
		t.Fatalf("expected load query failed, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestRegisterCourse_PrerequisiteQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	service := NewAcademicService(db)

	mock.ExpectBegin()

	// Course lookup succeeds
	mock.ExpectQuery(`SELECT\s+credit_units,\s+level,\s+max_capacity,\s+current_enrolled,\s+semester\s+FROM courses`).WithArgs("CSC401").
		WillReturnRows(sqlmock.NewRows([]string{
			"credit_units", "level", "max_capacity", "current_enrolled", "semester"}).AddRow(3, 400, 100, 10, "First"))

	// Student lookup succeeds
	mock.ExpectQuery("SELECT level FROM users").
		WithArgs("STU001").
		WillReturnRows(sqlmock.NewRows([]string{
			"level",
		}).AddRow(400))

	// Current load succeeds
	mock.ExpectQuery("SELECT COALESCE\\(SUM\\(c.credit_units\\), 0\\)").
		WithArgs("STU001", "2025/2026", "First").
		WillReturnRows(sqlmock.NewRows([]string{
			"sum",
		}).AddRow(12))

	// Prerequisite query fails
	mock.ExpectQuery("SELECT prerequisite_code FROM course_prerequisites").
		WithArgs("CSC401").
		WillReturnError(errors.New("prerequisite query failed"))

	mock.ExpectRollback()

	err = service.RegisterCourse(
		"STU001",
		"CSC401",
		"2025/2026",
		"First",
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "prerequisite query failed" {
		t.Fatalf("expected prerequisite query failed, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestRegisterCourse_PrerequisiteScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	service := NewAcademicService(db)

	mock.ExpectBegin()

	// Course lookup
	mock.ExpectQuery(`SELECT\s+credit_units,\s+level,\s+max_capacity,\s+current_enrolled,\s+semester\s+FROM courses`).
		WithArgs("CS101").
		WillReturnRows(
			sqlmock.NewRows([]string{
				"credit_units", "level", "max_capacity", "current_enrolled", "semester"}).AddRow(3, 100, 100, 20, "First"),
		)

	// Student lookup
	mock.ExpectQuery(`SELECT level FROM users`).
		WithArgs("STU001").
		WillReturnRows(
			sqlmock.NewRows([]string{"level"}).
				AddRow(100),
		)

	// Current load
	mock.ExpectQuery(`SELECT COALESCE\(SUM\(c.credit_units\), 0\)`).
		WithArgs("STU001", "2025/2026", "First").
		WillReturnRows(
			sqlmock.NewRows([]string{"total"}).
				AddRow(12),
		)

	// Prerequisite query returns an invalid type for Scan(&string)
	mock.ExpectQuery(`SELECT prerequisite_code FROM course_prerequisites`).
		WithArgs("CS101").
		WillReturnRows(
			sqlmock.NewRows([]string{"prerequisite_code"}).
				AddRow(123),
		)

	mock.ExpectRollback()

	err = service.RegisterCourse(
		"STU001",
		"CS101",
		"2025/2026",
		"First",
	)

	if err == nil {
		t.Fatal("expected prerequisite scan error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestRegisterCourse_PrerequisiteCheckError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	service := NewAcademicService(db)

	mock.ExpectBegin()

	// Course lookup
	mock.ExpectQuery(`SELECT\s+credit_units,\s+level,\s+max_capacity,\s+current_enrolled,\s+semester\s+FROM courses`).
		WithArgs("CS101").
		WillReturnRows(
			sqlmock.NewRows([]string{
				"credit_units", "level", "max_capacity", "current_enrolled", "semester"}).AddRow(3, 100, 100, 20, "First"),
		)

	// Student lookup
	mock.ExpectQuery(`SELECT level FROM users`).
		WithArgs("STU001").
		WillReturnRows(
			sqlmock.NewRows([]string{"level"}).
				AddRow(100),
		)

	// Current load
	mock.ExpectQuery(`SELECT COALESCE\(SUM\(c.credit_units\), 0\)`).
		WithArgs("STU001", "2025/2026", "First").
		WillReturnRows(
			sqlmock.NewRows([]string{"total"}).
				AddRow(12),
		)

	// One prerequisite
	mock.ExpectQuery(`SELECT prerequisite_code FROM course_prerequisites`).
		WithArgs("CS101").
		WillReturnRows(
			sqlmock.NewRows([]string{"prerequisite_code"}).
				AddRow("CSC100"),
		)

	// Force prerequisite lookup error
	mock.ExpectQuery(`SELECT EXISTS`).
		WithArgs("STU001", "CSC100").
		WillReturnError(errors.New("database error"))
	mock.ExpectRollback()

	err = service.RegisterCourse(
		"STU001",
		"CS101",
		"2025/2026",
		"First",
	)

	if err == nil {
		t.Fatal("expected error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestRegisterCourse_DuplicateRegistration(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	service := NewAcademicService(db)

	mock.ExpectBegin()

	// 1. Course lookup.
	mock.ExpectQuery(`SELECT\s+credit_units,\s+level,\s+max_capacity,\s+current_enrolled,\s+semester\s+FROM courses`).WithArgs("CSC401").
		WillReturnRows(
			sqlmock.NewRows([]string{
				"credit_units", "level", "max_capacity", "current_enrolled", "semester"}).AddRow(3, 400, 100, 50, "First"),
		)

	// 2. Student lookup.
	mock.ExpectQuery(
		"SELECT level FROM users",
	).
		WithArgs("STU001").
		WillReturnRows(
			sqlmock.NewRows([]string{
				"level",
			}).AddRow(400),
		)

	// 3. Current credit load.
	mock.ExpectQuery(
		"SELECT COALESCE\\(SUM\\(c.credit_units\\), 0\\)",
	).
		WithArgs(
			"STU001",
			"2025/2026",
			"First",
		).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"coalesce",
			}).AddRow(0),
		)

	// 4. Prerequisite lookup.
	mock.ExpectQuery(
		"SELECT prerequisite_code FROM course_prerequisites",
	).
		WithArgs("CSC401").
		WillReturnRows(
			sqlmock.NewRows([]string{
				"prerequisite_code",
			}),
		)

	// 5. Registration INSERT fails because the
	// student is already registered.
	mock.ExpectExec(
		"INSERT INTO student_courses",
	).
		WithArgs(
			"STU001",
			"CSC401",
			"2025/2026",
			"First",
		).
		WillReturnError(&pq.Error{
			Code: "23505",
		})

	mock.ExpectRollback()

	err = service.RegisterCourse(
		"STU001",
		"CSC401",
		"2025/2026",
		"First",
	)

	if err == nil {
		t.Fatal("expected duplicate registration error")
	}

	expected := "course registration conflict: you are already registered for CSC401 for the 2025/2026 session"

	if err.Error() != expected {
		t.Fatalf(
			"expected %q, got %q",
			expected,
			err.Error(),
		)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(
			"unmet sqlmock expectations: %v",
			err,
		)
	}
}

func TestRegisterCourse_NullCourseSemester(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	service := NewAcademicService(db)

	mock.ExpectBegin()

	mock.ExpectQuery(
		`SELECT\s+credit_units,\s+level,\s+max_capacity,\s+current_enrolled,\s+semester\s+FROM courses`,
	).
		WithArgs("CSC401").
		WillReturnRows(
			sqlmock.NewRows([]string{
				"credit_units",
				"level",
				"max_capacity",
				"current_enrolled",
				"semester",
			}).AddRow(3, 400, 100, 50, nil),
		)

	mock.ExpectRollback()

	err = service.RegisterCourse(
		"STU001",
		"CSC401",
		"2025/2026",
		"First",
	)

	if err == nil {
		t.Fatal("expected error for course with NULL semester")
	}

	expected := "academic rule violation: course CSC401 is not offered in the First semester"

	if err.Error() != expected {
		t.Fatalf(
			"expected %q, got %q",
			expected,
			err.Error(),
		)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}
