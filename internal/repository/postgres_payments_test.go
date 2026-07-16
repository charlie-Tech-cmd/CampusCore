package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"errors"
	"time"

	"campuscore/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestNewPostgresFinancialRepository(t *testing.T) {
	repo := NewPostgresFinancialRepository(nil)

	if repo == nil {
		t.Fatal("expected repository, got nil")
	}

	if repo.db != nil {
		t.Fatal("expected nil db")
	}
}

func TestGetFeeStructure_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresFinancialRepository(db)

	query := regexp.QuoteMeta(`
                SELECT
                        id,
                        department_id,
                        level,
                        fee_type,
                        amount_required,
                        session
                FROM fee_structures
                WHERE department_id = $1
                        AND level = $2
                        AND fee_type = $3
                        AND session = $4
        `)

	rows := sqlmock.NewRows([]string{
		"id",
		"department_id",
		"level",
		"fee_type",
		"amount_required",
		"session",
	}).AddRow(
		1,
		2,
		300,
		"school_fees",
		50000.00,
		"2025/2026",
	)

	mock.ExpectQuery(query).
		WithArgs(2, 300, "school_fees", "2025/2026").
		WillReturnRows(rows)

	fee, err := repo.GetFeeStructure(
		context.Background(),
		2,
		300,
		"school_fees",
		"2025/2026",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := &models.FeeStructure{
		ID:             1,
		DepartmentID:   2,
		Level:          300,
		FeeType:        "school_fees",
		AmountRequired: 50000,
		Session:        "2025/2026",
	}

	if *fee != *expected {
		t.Fatalf("expected %+v, got %+v", *expected, *fee)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestGetFeeStructure_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresFinancialRepository(db)

	query := regexp.QuoteMeta(`
                SELECT
                        id,
                        department_id,
                        level,
                        fee_type,
                        amount_required,
                        session
                FROM fee_structures
                WHERE department_id = $1
                        AND level = $2
                        AND fee_type = $3
                        AND session = $4
        `)

	mock.ExpectQuery(query).
		WithArgs(2, 300, "school_fees", "2025/2026").
		WillReturnError(sql.ErrNoRows)

	fee, err := repo.GetFeeStructure(
		context.Background(),
		2,
		300,
		"school_fees",
		"2025/2026",
	)

	if err != sql.ErrNoRows {
		t.Fatalf("expected sql.ErrNoRows, got %v", err)
	}

	if fee != nil {
		t.Fatal("expected nil fee")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestGetFeeStructure_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresFinancialRepository(db)

	query := regexp.QuoteMeta(`
                SELECT
                        id,
                        department_id,
                        level,
                        fee_type,
                        amount_required,
                        session
                FROM fee_structures
                WHERE department_id = $1
                        AND level = $2
                        AND fee_type = $3
                        AND session = $4
        `)

	rows := sqlmock.NewRows([]string{
		"id",
		"department_id",
		"level",
		"fee_type",
		"amount_required",
		"session",
	}).AddRow(
		"invalid", // ID should be an integer
		2,
		300,
		"school_fees",
		50000.00,
		"2025/2026",
	)

	mock.ExpectQuery(query).
		WithArgs(2, 300, "school_fees", "2025/2026").
		WillReturnRows(rows)

	fee, err := repo.GetFeeStructure(
		context.Background(),
		2,
		300,
		"school_fees",
		"2025/2026",
	)

	if err == nil {
		t.Fatal("expected scan error, got nil")
	}

	if fee != nil {
		t.Fatal("expected nil fee")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestRecordPayment_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresFinancialRepository(db)

	payment := &models.FeePayment{
		StudentID:        "STU001",
		GatewayReference: "PAY123",
		AmountPaid:       50000,
		FeeType:          "school_fees",
		Session:          "2025/2026",
		Status:           "successful",
	}

	query := regexp.QuoteMeta(`
                INSERT INTO fee_payments
                        (student_id,
                         gateway_reference,
                         amount_paid,
                         fee_type,
                         session,
                         status)
                VALUES ($1, $2, $3, $4, $5, $6)
        `)

	mock.ExpectExec(query).
		WithArgs(
			payment.StudentID,
			payment.GatewayReference,
			payment.AmountPaid,
			payment.FeeType,
			payment.Session,
			payment.Status,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.RecordPayment(context.Background(), payment)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestRecordPayment_ExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresFinancialRepository(db)

	payment := &models.FeePayment{
		StudentID:        "STU001",
		GatewayReference: "PAY123",
		AmountPaid:       50000,
		FeeType:          "school_fees",
		Session:          "2025/2026",
		Status:           "successful",
	}

	query := regexp.QuoteMeta(`
                INSERT INTO fee_payments
                        (student_id,
                         gateway_reference,
                         amount_paid,
                         fee_type,
                         session,
                         status)
                VALUES ($1, $2, $3, $4, $5, $6)
        `)

	mock.ExpectExec(query).
		WithArgs(
			payment.StudentID,
			payment.GatewayReference,
			payment.AmountPaid,
			payment.FeeType,
			payment.Session,
			payment.Status,
		).
		WillReturnError(errors.New("insert failed"))

	err = repo.RecordPayment(context.Background(), payment)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "insert failed" {
		t.Fatalf("expected insert failed, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestCheckPaymentExists_True(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresFinancialRepository(db)

	query := regexp.QuoteMeta(`
                SELECT EXISTS(
                        SELECT 1
                        FROM fee_payments
                        WHERE gateway_reference = $1
                )
        `)

	mock.ExpectQuery(query).
		WithArgs("PAY123").
		WillReturnRows(
			sqlmock.NewRows([]string{"exists"}).
				AddRow(true),
		)

	exists, err := repo.CheckPaymentExists(context.Background(), "PAY123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !exists {
		t.Fatal("expected payment to exist")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestCheckPaymentExists_False(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresFinancialRepository(db)

	query := regexp.QuoteMeta(`
                SELECT EXISTS(
                        SELECT 1
                        FROM fee_payments
                        WHERE gateway_reference = $1
                )
        `)

	mock.ExpectQuery(query).
		WithArgs("PAY999").
		WillReturnRows(
			sqlmock.NewRows([]string{"exists"}).
				AddRow(false),
		)

	exists, err := repo.CheckPaymentExists(context.Background(), "PAY999")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if exists {
		t.Fatal("expected payment not to exist")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestCheckPaymentExists_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresFinancialRepository(db)

	query := regexp.QuoteMeta(`
                SELECT EXISTS(
                        SELECT 1
                        FROM fee_payments
                        WHERE gateway_reference = $1
                )
        `)

	mock.ExpectQuery(query).
		WithArgs("PAY123").
		WillReturnError(errors.New("database error"))

	exists, err := repo.CheckPaymentExists(context.Background(), "PAY123")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if exists {
		t.Fatal("expected exists to be false")
	}

	if err.Error() != "database error" {
		t.Fatalf("expected database error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestGetStudentClearanceStatus_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresFinancialRepository(db)

	query := regexp.QuoteMeta(`
                SELECT
                        sc.id,
                        sc.student_id,
                        sc.office_id,
                        co.office_name,
                        sc.status,
                        sc.assigned_staff_id,
                        sc.updated_at
                FROM student_clearances sc
                JOIN clearance_offices co
                        ON sc.office_id = co.id
                WHERE sc.student_id = $1
        `)

	rows := sqlmock.NewRows([]string{
		"id",
		"student_id",
		"office_id",
		"office_name",
		"status",
		"assigned_staff_id",
		"updated_at",
	}).AddRow(
		1,
		"STU001",
		2,
		"Bursary",
		"approved",
		"STAFF001",
		time.Now(),
	)

	mock.ExpectQuery(query).
		WithArgs("STU001").
		WillReturnRows(rows)

	clearances, err := repo.GetStudentClearanceStatus(
		context.Background(),
		"STU001",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(clearances) != 1 {
		t.Fatalf("expected 1 clearance, got %d", len(clearances))
	}

	if clearances[0].StudentID != "STU001" {
		t.Fatalf("expected STU001, got %s", clearances[0].StudentID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestGetStudentClearanceStatus_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresFinancialRepository(db)

	query := regexp.QuoteMeta(`
                SELECT
                        sc.id,
                        sc.student_id,
                        sc.office_id,
                        co.office_name,
                        sc.status,
                        sc.assigned_staff_id,
                        sc.updated_at
                FROM student_clearances sc
                JOIN clearance_offices co
                        ON sc.office_id = co.id
                WHERE sc.student_id = $1
        `)

	mock.ExpectQuery(query).
		WithArgs("STU001").
		WillReturnError(errors.New("query failed"))

	_, err = repo.GetStudentClearanceStatus(context.Background(), "STU001")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "query failed" {
		t.Fatalf("expected query failed, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestGetStudentClearanceStatus_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresFinancialRepository(db)

	query := regexp.QuoteMeta(`
                SELECT
                        sc.id,
                        sc.student_id,
                        sc.office_id,
                        co.office_name,
                        sc.status,
                        sc.assigned_staff_id,
                        sc.updated_at
                FROM student_clearances sc
                JOIN clearance_offices co
                        ON sc.office_id = co.id
                WHERE sc.student_id = $1
        `)

	rows := sqlmock.NewRows([]string{
		"id",
		"student_id",
		"office_id",
		"office_name",
		"status",
		"assigned_staff_id",
		"updated_at",
	}).AddRow(
		"invalid",
		"STU001",
		2,
		"Bursary",
		"approved",
		"STAFF001",
		time.Now(),
	)

	mock.ExpectQuery(query).
		WithArgs("STU001").
		WillReturnRows(rows)

	_, err = repo.GetStudentClearanceStatus(context.Background(), "STU001")

	if err == nil {
		t.Fatal("expected scan error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestGetStudentClearanceStatus_EmptyResult(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresFinancialRepository(db)

	query := regexp.QuoteMeta(`
                SELECT
                        sc.id,
                        sc.student_id,
                        sc.office_id,
                        co.office_name,
                        sc.status,
                        sc.assigned_staff_id,
                        sc.updated_at
                FROM student_clearances sc
                JOIN clearance_offices co
                        ON sc.office_id = co.id
                WHERE sc.student_id = $1
        `)

	rows := sqlmock.NewRows([]string{
		"id",
		"student_id",
		"office_id",
		"office_name",
		"status",
		"assigned_staff_id",
		"updated_at",
	})

	mock.ExpectQuery(query).
		WithArgs("STU001").
		WillReturnRows(rows)

	clearances, err := repo.GetStudentClearanceStatus(
		context.Background(),
		"STU001",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(clearances) != 0 {
		t.Fatalf("expected 0 records, got %d", len(clearances))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestGetStudentClearanceStatus_RowsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresFinancialRepository(db)

	query := regexp.QuoteMeta(`
                SELECT
                        sc.id,
                        sc.student_id,
                        sc.office_id,
                        co.office_name,
                        sc.status,
                        sc.assigned_staff_id,
                        sc.updated_at
                FROM student_clearances sc
                JOIN clearance_offices co
                        ON sc.office_id = co.id
                WHERE sc.student_id = $1
        `)

	rows := sqlmock.NewRows([]string{
		"id",
		"student_id",
		"office_id",
		"office_name",
		"status",
		"assigned_staff_id",
		"updated_at",
	}).
		AddRow(
			1,
			"STU001",
			2,
			"Bursary",
			"approved",
			"STAFF001",
			time.Now(),
		).
		RowError(0, errors.New("row iteration error"))

	mock.ExpectQuery(query).
		WithArgs("STU001").
		WillReturnRows(rows)

	_, err = repo.GetStudentClearanceStatus(
		context.Background(),
		"STU001",
	)

	if err == nil {
		t.Fatal("expected row iteration error")
	}

	if err.Error() != "row iteration error" {
		t.Fatalf("expected row iteration error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateClearanceStatus_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresFinancialRepository(db)

	query := regexp.QuoteMeta(`
                UPDATE student_clearances
                SET
                        status = $1,
                        assigned_staff_id = $2,
                        updated_at = NOW()
                WHERE student_id = $3
                        AND office_id = $4
        `)

	mock.ExpectExec(query).
		WithArgs(
			models.ClearanceStatus("approved"),
			"STAFF001",
			"STU001",
			2,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.UpdateClearanceStatus(
		context.Background(),
		"STU001",
		2,
		models.ClearanceStatus("approved"),
		"STAFF001",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateClearanceStatus_ExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresFinancialRepository(db)

	query := regexp.QuoteMeta(`
                UPDATE student_clearances
                SET
                        status = $1,
                        assigned_staff_id = $2,
                        updated_at = NOW()
                WHERE student_id = $3
                        AND office_id = $4
        `)

	mock.ExpectExec(query).
		WithArgs(
			models.ClearanceStatus("approved"),
			"STAFF001",
			"STU001",
			2,
		).
		WillReturnError(errors.New("update failed"))

	err = repo.UpdateClearanceStatus(
		context.Background(),
		"STU001",
		2,
		models.ClearanceStatus("approved"),
		"STAFF001",
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "update failed" {
		t.Fatalf("expected update failed, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestCreateTicket_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresFinancialRepository(db)

	query := regexp.QuoteMeta(`
                INSERT INTO support_tickets
                        (student_id,
                         category,
                         status,
                         subject,
                         message)
                VALUES ($1, $2, $3, $4, $5)
        `)

	ticket := &models.SupportTicket{
		StudentID: "STU001",
		Category:  "Academic",
		Status:    "Open",
		Subject:   "Course Registration",
		Message:   "Unable to register courses.",
	}

	mock.ExpectExec(query).
		WithArgs(
			ticket.StudentID,
			ticket.Category,
			ticket.Status,
			ticket.Subject,
			ticket.Message,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateTicket(context.Background(), ticket)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestCreateTicket_ExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresFinancialRepository(db)

	query := regexp.QuoteMeta(`
                INSERT INTO support_tickets
                        (student_id,
                         category,
                         status,
                         subject,
                         message)
                VALUES ($1, $2, $3, $4, $5)
        `)

	ticket := &models.SupportTicket{
		StudentID: "STU001",
		Category:  "Academic",
		Status:    "Open",
		Subject:   "Course Registration",
		Message:   "Unable to register courses.",
	}

	mock.ExpectExec(query).
		WithArgs(
			ticket.StudentID,
			ticket.Category,
			ticket.Status,
			ticket.Subject,
			ticket.Message,
		).
		WillReturnError(errors.New("insert failed"))

	err = repo.CreateTicket(context.Background(), ticket)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "insert failed" {
		t.Fatalf("expected insert failed, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestNewPostgresUserRepository(t *testing.T) {
	repo := NewPostgresUserRepository(nil)

	if repo == nil {
		t.Fatal("expected repository, got nil")
	}

	if repo.db != nil {
		t.Fatal("expected nil db")
	}
}

func TestCreate_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresUserRepository(db)

	query := regexp.QuoteMeta(`
                INSERT INTO users (id, surname, first_name, middle_name, email, phone, password_hash, role, department_id, level)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`)

	user := &models.User{
		ID:           "STU001",
		Surname:      "Doe",
		FirstName:    "John",
		MiddleName:   "A",
		Email:        "john@example.com",
		Phone:        "08012345678",
		PasswordHash: "hashed-password",
		Role:         "student",
		DepartmentID: 2,
		Level:        300,
	}

	mock.ExpectExec(query).
		WithArgs(
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
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(user)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestCreate_ExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresUserRepository(db)

	query := regexp.QuoteMeta(`
                INSERT INTO users (id, surname, first_name, middle_name, email, phone, password_hash, role, department_id, level)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`)

	user := &models.User{
		ID:           "STU001",
		Surname:      "Doe",
		FirstName:    "John",
		MiddleName:   "A",
		Email:        "john@example.com",
		Phone:        "08012345678",
		PasswordHash: "hashed-password",
		Role:         "student",
		DepartmentID: 2,
		Level:        300,
	}

	mock.ExpectExec(query).
		WithArgs(
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
		).
		WillReturnError(errors.New("insert failed"))

	err = repo.Create(user)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestFindByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresUserRepository(db)

	query := regexp.QuoteMeta(`
                SELECT id, surname, first_name, middle_name, email, phone, password_hash, role, department_id, level, last_login, created_at
                FROM users
                WHERE id = $1
                LIMIT 1;`)

	now := time.Now()
	created := now.Add(-24 * time.Hour)

	rows := sqlmock.NewRows([]string{
		"id",
		"surname",
		"first_name",
		"middle_name",
		"email",
		"phone",
		"password_hash",
		"role",
		"department_id",
		"level",
		"last_login",
		"created_at",
	}).AddRow(
		"STU001",
		"Doe",
		"John",
		"A",
		"john@example.com",
		"08012345678",
		"hash",
		"student",
		2,
		300,
		now,
		created,
	)

	mock.ExpectQuery(query).
		WithArgs("STU001").
		WillReturnRows(rows)

	user, err := repo.FindByID("STU001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user.ID != "STU001" {
		t.Fatalf("expected STU001, got %s", user.ID)
	}

	if user.DepartmentID != 2 {
		t.Fatalf("expected department 2")
	}

	if user.LastLogin.IsZero() {
		t.Fatal("expected last login to be populated")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestFindByID_NullFields(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresUserRepository(db)

	query := regexp.QuoteMeta(`
                SELECT id, surname, first_name, middle_name, email, phone, password_hash, role, department_id, level, last_login, created_at
                FROM users
                WHERE id = $1
                LIMIT 1;`)

	created := time.Now()

	rows := sqlmock.NewRows([]string{
		"id",
		"surname",
		"first_name",
		"middle_name",
		"email",
		"phone",
		"password_hash",
		"role",
		"department_id",
		"level",
		"last_login",
		"created_at",
	}).AddRow(
		"STU001",
		"Doe",
		"John",
		"A",
		"john@example.com",
		"08012345678",
		"hash",
		"student",
		nil,
		300,
		nil,
		created,
	)

	mock.ExpectQuery(query).
		WithArgs("STU001").
		WillReturnRows(rows)

	user, err := repo.FindByID("STU001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user.DepartmentID != 0 {
		t.Fatalf("expected DepartmentID to remain 0, got %d", user.DepartmentID)
	}

	if !user.LastLogin.IsZero() {
		t.Fatal("expected zero LastLogin")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestFindByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresUserRepository(db)

	query := regexp.QuoteMeta(`
                SELECT id, surname, first_name, middle_name, email, phone, password_hash, role, department_id, level, last_login, created_at
                FROM users
                WHERE id = $1
                LIMIT 1;`)

	mock.ExpectQuery(query).
		WithArgs("UNKNOWN").
		WillReturnError(sql.ErrNoRows)

	_, err = repo.FindByID("UNKNOWN")

	if err == nil {
		t.Fatal("expected error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestFindByID_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresUserRepository(db)

	query := regexp.QuoteMeta(`
                SELECT id, surname, first_name, middle_name, email, phone, password_hash, role, department_id, level, last_login, created_at
                FROM users
                WHERE id = $1
                LIMIT 1;`)

	rows := sqlmock.NewRows([]string{
	"id",
	"surname",
	"first_name",
	"middle_name",
	"email",
	"phone",
	"password_hash",
	"role",
	"department_id",
	"level",
	"last_login",
	"created_at",
}).AddRow(
	"STU001",
	"Doe",
	"John",
		"A",
		"john@example.com",
		"08012345678",
		"hash",
		"student",
		"invalid_department", // should fail scanning into sql.NullInt32
		300,
		time.Now(),
		time.Now(),
)

	mock.ExpectQuery(query).
		WithArgs("STU001").
		WillReturnRows(rows)

	_, err = repo.FindByID("STU001")

	if err == nil {
		t.Fatal("expected scan error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestFindByEmail_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresUserRepository(db)

	query := regexp.QuoteMeta(`
                SELECT id, surname, first_name, middle_name, email, phone, password_hash, role, department_id, level, last_login, created_at
                FROM users
                WHERE email = $1
                LIMIT 1;`)

	now := time.Now()
	created := now.Add(-24 * time.Hour)

	rows := sqlmock.NewRows([]string{
		"id",
		"surname",
		"first_name",
		"middle_name",
		"email",
		"phone",
		"password_hash",
		"role",
		"department_id",
		"level",
		"last_login",
		"created_at",
	}).AddRow(
		"STU001",
		"Doe",
		"John",
		"A",
		"john@example.com",
		"08012345678",
		"hash",
		"student",
		2,
		300,
		now,
		created,
	)

	mock.ExpectQuery(query).
		WithArgs("john@example.com").
		WillReturnRows(rows)

	user, err := repo.FindByEmail("john@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user.Email != "john@example.com" {
		t.Fatalf("expected john@example.com, got %s", user.Email)
	}

	if user.DepartmentID != 2 {
		t.Fatalf("expected department 2, got %d", user.DepartmentID)
	}

	if user.LastLogin.IsZero() {
		t.Fatal("expected LastLogin to be populated")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestFindByEmail_NullFields(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresUserRepository(db)

	query := regexp.QuoteMeta(`
                SELECT id, surname, first_name, middle_name, email, phone, password_hash, role, department_id, level, last_login, created_at
                FROM users
                WHERE email = $1
                LIMIT 1;`)

	created := time.Now()

	rows := sqlmock.NewRows([]string{
		"id",
		"surname",
		"first_name",
		"middle_name",
		"email",
		"phone",
		"password_hash",
		"role",
		"department_id",
		"level",
		"last_login",
		"created_at",
	}).AddRow(
		"STU001",
		"Doe",
		"John",
		"A",
		"john@example.com",
		"08012345678",
		"hash",
		"student",
		nil,
		300,
		nil,
		created,
	)

	mock.ExpectQuery(query).
		WithArgs("john@example.com").
		WillReturnRows(rows)

	user, err := repo.FindByEmail("john@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user.DepartmentID != 0 {
		t.Fatalf("expected DepartmentID 0, got %d", user.DepartmentID)
	}

	if !user.LastLogin.IsZero() {
		t.Fatal("expected zero LastLogin")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestFindByEmail_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresUserRepository(db)

	query := regexp.QuoteMeta(`
                SELECT id, surname, first_name, middle_name, email, phone, password_hash, role, department_id, level, last_login, created_at
                FROM users
                WHERE email = $1
                LIMIT 1;`)

	mock.ExpectQuery(query).
		WithArgs("missing@example.com").
		WillReturnError(sql.ErrNoRows)

	_, err = repo.FindByEmail("missing@example.com")
	if err == nil {
		t.Fatal("expected error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestFindByEmail_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresUserRepository(db)

	query := regexp.QuoteMeta(`
                SELECT id, surname, first_name, middle_name, email, phone, password_hash, role, department_id, level, last_login, created_at
                FROM users
                WHERE email = $1
                LIMIT 1;`)

	rows := sqlmock.NewRows([]string{
		"id",
		"surname",
		"first_name",
		"middle_name",
		"email",
		"phone",
		"password_hash",
		"role",
		"department_id",
		"level",
		"last_login",
		"created_at",
	}).AddRow(
		"STU001",
		"Doe",
		"John",
		"A",
		"john@example.com",
		"08012345678",
		"hash",
		"student",
		"invalid_department",
		300,
		time.Now(),
		time.Now(),
	)

	mock.ExpectQuery(query).
		WithArgs("john@example.com").
		WillReturnRows(rows)

	_, err = repo.FindByEmail("john@example.com")
	if err == nil {
		t.Fatal("expected scan error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateLastLogin_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresUserRepository(db)

	query := regexp.QuoteMeta(`
UPDATE users SET last_login = CURRENT_TIMESTAMP WHERE id = $1;`)

	mock.ExpectExec(query).
		WithArgs("STU001").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.UpdateLastLogin("STU001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateLastLogin_ExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresUserRepository(db)

	query := regexp.QuoteMeta(`
UPDATE users SET last_login = CURRENT_TIMESTAMP WHERE id = $1;`)

	mock.ExpectExec(query).
		WithArgs("STU001").
		WillReturnError(errors.New("update failed"))

	err = repo.UpdateLastLogin("STU001")
	if err == nil {
		t.Fatal("expected error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestNewPostgresGovernanceRepository(t *testing.T) {
	repo := NewPostgresGovernanceRepository(nil)

	if repo == nil {
		t.Fatal("expected repository")
	}

	if repo.db != nil {
		t.Fatal("expected nil db")
	}
}

func TestGetApprovalStatus_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewPostgresGovernanceRepository(db)

	query := regexp.QuoteMeta(`
                SELECT id, course_code, session, semester, current_state, action_by, remarks, updated_at
                FROM approvals
                WHERE course_code = $1
                LIMIT 1;`)

	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id",
		"course_code",
		"session",
		"semester",
		"current_state",
		"action_by",
		"remarks",
		"updated_at",
	}).AddRow(
		1,
		"CSC401",
		"2026/2027",
		"First",
		models.StatusHODApproved,
		"STAFF001",
		"Approved",
		now,
	)

	mock.ExpectQuery(query).
		WithArgs("CSC401").
		WillReturnRows(rows)

	approval, err := repo.GetApprovalStatus("CSC401")
	if err != nil {
		t.Fatal(err)
	}

	if approval.CourseCode != "CSC401" {
		t.Fatal("wrong course")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestGetApprovalStatus_NoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewPostgresGovernanceRepository(db)

	query := regexp.QuoteMeta(`
                SELECT id, course_code, session, semester, current_state, action_by, remarks, updated_at
                FROM approvals
                WHERE course_code = $1
                LIMIT 1;`)

	mock.ExpectQuery(query).
		WithArgs("CSC401").
		WillReturnError(sql.ErrNoRows)

	approval, err := repo.GetApprovalStatus("CSC401")
	if err != nil {
		t.Fatal(err)
	}

	if approval.CurrentState != models.StatusSubmitted {
		t.Fatal("expected default submitted state")
	}

	if approval.CourseCode != "CSC401" {
		t.Fatal("wrong course")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestGetApprovalStatus_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewPostgresGovernanceRepository(db)

	query := regexp.QuoteMeta(`
                SELECT id, course_code, session, semester, current_state, action_by, remarks, updated_at
                FROM approvals
                WHERE course_code = $1
                LIMIT 1;`)

	rows := sqlmock.NewRows([]string{
		"id",
		"course_code",
		"session",
		"semester",
		"current_state",
		"action_by",
		"remarks",
		"updated_at",
	}).AddRow(
		"invalid",
		"CSC401",
		"2026/2027",
		"First",
		models.StatusHODApproved,
		"STAFF001",
		"Approved",
		time.Now(),
	)

	mock.ExpectQuery(query).
		WithArgs("CSC401").
		WillReturnRows(rows)

	_, err = repo.GetApprovalStatus("CSC401")
	if err == nil {
		t.Fatal("expected scan error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateApprovalState_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewPostgresGovernanceRepository(db)

	query := regexp.QuoteMeta(`
                INSERT INTO approvals (course_code, session, semester, current_state, action_by, remarks, updated_at)
                VALUES ($1, '2026/2027', 'First', $2, $3, $4, CURRENT_TIMESTAMP)
                ON CONFLICT (course_code, session, semester)
                DO UPDATE SET
                        current_state = EXCLUDED.current_state,
                        action_by = EXCLUDED.action_by,
                        remarks = EXCLUDED.remarks,
                        updated_at = CURRENT_TIMESTAMP;`)

	mock.ExpectExec(query).
		WithArgs(
			"CSC401",
			models.StatusHODApproved,
			"STAFF001",
			"Approved",
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateApprovalState(
		"CSC401",
		models.StatusHODApproved,
		"STAFF001",
		"Approved",
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateApprovalState_ExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewPostgresGovernanceRepository(db)

	query := regexp.QuoteMeta(`
                INSERT INTO approvals (course_code, session, semester, current_state, action_by, remarks, updated_at)
                VALUES ($1, '2026/2027', 'First', $2, $3, $4, CURRENT_TIMESTAMP)
                ON CONFLICT (course_code, session, semester)
                DO UPDATE SET
                        current_state = EXCLUDED.current_state,
                        action_by = EXCLUDED.action_by,
                        remarks = EXCLUDED.remarks,
                        updated_at = CURRENT_TIMESTAMP;`)

	mock.ExpectExec(query).
		WithArgs(
			"CSC401",
			models.StatusHODApproved,
			"STAFF001",
			"Approved",
		).
		WillReturnError(errors.New("exec failed"))

	err = repo.UpdateApprovalState(
		"CSC401",
		models.StatusHODApproved,
		"STAFF001",
		"Approved",
	)

	if err == nil {
		t.Fatal("expected error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
