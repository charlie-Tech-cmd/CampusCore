package repository

import (
	"context"
	"database/sql"

	"campuscore/internal/models"
)

// PostgresFinancialRepository stores financial data in PostgreSQL.
type PostgresFinancialRepository struct {
	db *sql.DB
}

// NewPostgresFinancialRepository creates a new financial repository.
func NewPostgresFinancialRepository(db *sql.DB) *PostgresFinancialRepository {
	return &PostgresFinancialRepository{
		db: db,
	}
}

// GetFeeStructure returns the fee structure for a department and level.
func (r *PostgresFinancialRepository) GetFeeStructure(
	ctx context.Context,
	departmentID, level int,
	feeType, session string,
) (*models.FeeStructure, error) {

	query := `
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
	`

	var fee models.FeeStructure

	err := r.db.QueryRowContext(
		ctx,
		query,
		departmentID,
		level,
		feeType,
		session,
	).Scan(
		&fee.ID,
		&fee.DepartmentID,
		&fee.Level,
		&fee.FeeType,
		&fee.AmountRequired,
		&fee.Session,
	)

	if err != nil {
		return nil, err
	}

	return &fee, nil
}

// RecordPayment saves a payment.
func (r *PostgresFinancialRepository) RecordPayment(
	ctx context.Context,
	payment *models.FeePayment,
) error {

	query := `
		INSERT INTO fee_payments
			(student_id,
			 gateway_reference,
			 amount_paid,
			 fee_type,
			 session,
			 status)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		payment.StudentID,
		payment.GatewayReference,
		payment.AmountPaid,
		payment.FeeType,
		payment.Session,
		payment.Status,
	)

	return err
}

// CheckPaymentExists checks whether a payment already exists.
func (r *PostgresFinancialRepository) CheckPaymentExists(
	ctx context.Context,
	gatewayRef string,
) (bool, error) {

	query := `
		SELECT EXISTS(
			SELECT 1
			FROM fee_payments
			WHERE gateway_reference = $1
		)
	`

	var exists bool

	err := r.db.QueryRowContext(ctx, query, gatewayRef).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// GetStudentClearanceStatus returns a student's clearance records.
func (r *PostgresFinancialRepository) GetStudentClearanceStatus(
	ctx context.Context,
	studentID string,
) ([]models.StudentClearance, error) {

	query := `
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
	`

	rows, err := r.db.QueryContext(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clearances []models.StudentClearance

	for rows.Next() {
		var clearance models.StudentClearance

		err := rows.Scan(
			&clearance.ID,
			&clearance.StudentID,
			&clearance.OfficeID,
			&clearance.OfficeName,
			&clearance.Status,
			&clearance.AssignedStaff,
			&clearance.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		clearances = append(clearances, clearance)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return clearances, nil
}

// UpdateClearanceStatus updates one clearance record.
func (r *PostgresFinancialRepository) UpdateClearanceStatus(
	ctx context.Context,
	studentID string,
	officeID int,
	status models.ClearanceStatus,
	staffID string,
) error {

	query := `
		UPDATE student_clearances
		SET
			status = $1,
			assigned_staff_id = $2,
			updated_at = NOW()
		WHERE student_id = $3
			AND office_id = $4
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		status,
		staffID,
		studentID,
		officeID,
	)

	return err
}

// CreateTicket saves a support ticket.
func (r *PostgresFinancialRepository) CreateTicket(
	ctx context.Context,
	ticket *models.SupportTicket,
) error {

	query := `
		INSERT INTO support_tickets
			(student_id,
			 category,
			 status,
			 subject,
			 message)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		ticket.StudentID,
		ticket.Category,
		ticket.Status,
		ticket.Subject,
		ticket.Message,
	)

	return err
}