package repository

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"

	"campuscore/internal/models"
)

func TestEnrollmentRegisterSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewPostgresEnrollmentRepository(db)
	enrollment := &models.Enrollment{
		StudentID:  "STU001",
		CourseCode: "CSC401",
		Session:    "2025/2026",
		Semester:   "First",
		Status:     "approved",
	}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO student_courses")).
		WithArgs("STU001", "CSC401", "2025/2026", "First", "approved").
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := repo.Register(enrollment); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestEnrollmentRegisterClosedSession(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewPostgresEnrollmentRepository(db)
	enrollment := &models.Enrollment{
		StudentID:  "STU001",
		CourseCode: "CSC401",
		Session:    "2025/2026",
		Semester:   "First",
		Status:     "approved",
	}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO student_courses")).
		WithArgs("STU001", "CSC401", "2025/2026", "First", "approved").
		WillReturnResult(sqlmock.NewResult(0, 0))

	if err := repo.Register(enrollment); err == nil {
		t.Fatal("expected an error for a closed or inactive session")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestEnrollmentRegisterDuplicate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewPostgresEnrollmentRepository(db)
	enrollment := &models.Enrollment{
		StudentID:  "STU001",
		CourseCode: "CSC401",
		Session:    "2025/2026",
		Semester:   "First",
		Status:     "approved",
	}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO student_courses")).
		WithArgs("STU001", "CSC401", "2025/2026", "First", "approved").
		WillReturnError(&pq.Error{Code: "23505"})

	if err := repo.Register(enrollment); !errors.Is(err, ErrDuplicateEnrollment) {
		t.Fatalf("expected ErrDuplicateEnrollment, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestEnrollmentUpdateStatusSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewPostgresEnrollmentRepository(db)

	mock.ExpectExec(regexp.QuoteMeta("UPDATE student_courses")).
		WithArgs("STU001", "CSC401", "2025/2026", "First", "dropped").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.UpdateStatus("STU001", "CSC401", "2025/2026", "First", "dropped")
	if err != nil {
		t.Fatalf("UpdateStatus() error = %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestEnrollmentUpdateStatusNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewPostgresEnrollmentRepository(db)

	mock.ExpectExec(regexp.QuoteMeta("UPDATE student_courses")).
		WithArgs("STU001", "CSC401", "2025/2026", "First", "dropped").
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.UpdateStatus("STU001", "CSC401", "2025/2026", "First", "dropped")
	if err == nil || err.Error() != "enrollment not found" {
		t.Fatalf("expected enrollment not found, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestEnrollmentDeleteSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewPostgresEnrollmentRepository(db)

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM student_courses")).
		WithArgs("STU001", "CSC401", "2025/2026", "First").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete("STU001", "CSC401", "2025/2026", "First")
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestEnrollmentDeleteNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewPostgresEnrollmentRepository(db)

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM student_courses")).
		WithArgs("STU001", "CSC401", "2025/2026", "First").
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.Delete("STU001", "CSC401", "2025/2026", "First")
	if err == nil || err.Error() != "enrollment not found" {
		t.Fatalf("expected enrollment not found, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestEnrollmentFindByStudent(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewPostgresEnrollmentRepository(db)
	registeredAt := time.Date(2026, 10, 10, 7, 0, 0, 0, time.UTC)

	rows := sqlmock.NewRows([]string{
		"id",
		"student_id",
		"course_code",
		"credit_units",
		"session_id",
		"session_name",
		"semester",
		"status",
		"registered_at",
	}).AddRow(
		1,
		"STU001",
		"CSC401",
		3,
		7,
		"2025/2026",
		"First",
		"approved",
		registeredAt,
	)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT sc.id")).
		WithArgs("STU001").
		WillReturnRows(rows)

	enrollments, err := repo.FindByStudent("STU001")
	if err != nil {
		t.Fatalf("FindByStudent() error = %v", err)
	}

	if len(enrollments) != 1 {
		t.Fatalf("expected 1 enrollment, got %d", len(enrollments))
	}

	got := enrollments[0]
	if got.StudentID != "STU001" ||
		got.CourseCode != "CSC401" ||
		got.CreditUnits != 3 ||
		got.SessionID != 7 ||
		got.Session != "2025/2026" ||
		got.Semester != "First" ||
		got.Status != "approved" ||
		!got.CreatedAt.Equal(registeredAt) {
		t.Fatalf("unexpected enrollment: %+v", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
