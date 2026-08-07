package repository

import (
	"database/sql"

	"campuscore/internal/models"
)

type PostgresBillingRepository struct {
	db *sql.DB
}

func NewPostgresBillingRepository(
	db *sql.DB,
) *PostgresBillingRepository {

	return &PostgresBillingRepository{
		db: db,
	}
}

func (r *PostgresBillingRepository) CreateInvoice(
	invoice *models.Invoice,
) error {

	query := `
		INSERT INTO invoices (
			invoice_number,
			owner_id,
			owner_type,
			fee_type,
			session,
			amount,
			status,
			due_date,
			created_at,
			updated_at
		)
		VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10
		)
	`

	_, err := r.db.Exec(
		query,
		invoice.InvoiceNumber,
		invoice.OwnerID,
		invoice.OwnerType,
		invoice.FeeType,
		invoice.Session,
		invoice.Amount,
		invoice.Status,
		invoice.DueDate,
		invoice.CreatedAt,
		invoice.UpdatedAt,
	)

	return err
}

func (r *PostgresBillingRepository) FindInvoiceByID(
	id int,
) (*models.Invoice, error) {

	query := `
		SELECT
			id,
			invoice_number,
			owner_id,
			owner_type,
			fee_type,
			session,
			amount,
			status,
			due_date,
			created_at,
			updated_at
		FROM invoices
		WHERE id = $1
	`

	invoice := &models.Invoice{}

	err := r.db.QueryRow(
		query,
		id,
	).Scan(
		&invoice.ID,
		&invoice.InvoiceNumber,
		&invoice.OwnerID,
		&invoice.OwnerType,
		&invoice.FeeType,
		&invoice.Session,
		&invoice.Amount,
		&invoice.Status,
		&invoice.DueDate,
		&invoice.CreatedAt,
		&invoice.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return invoice, nil
}

func (r *PostgresBillingRepository) FindOwnerInvoices(
	ownerID string,
	ownerType string,
) ([]models.Invoice, error) {

	return nil, nil
}

func (r *PostgresBillingRepository) UpdateInvoice(
	invoice *models.Invoice,
) error {

	query := `
		UPDATE invoices
		SET
			owner_id = $1,
			owner_type = $2,
			fee_type = $3,
			session = $4,
			amount = $5,
			status = $6,
			due_date = $7,
			updated_at = $8
		WHERE id = $9
	`

	_, err := r.db.Exec(
		query,
		invoice.OwnerID,
		invoice.OwnerType,
		invoice.FeeType,
		invoice.Session,
		invoice.Amount,
		invoice.Status,
		invoice.DueDate,
		invoice.UpdatedAt,
		invoice.ID,
	)

	return err
}

func (r *PostgresBillingRepository) MarkInvoicePaid(
	id int,
) error {

	query := `
		UPDATE invoices
		SET
			status = $1,
			updated_at = NOW()
		WHERE id = $2
	`

	_, err := r.db.Exec(
		query,
		models.InvoicePaid,
		id,
	)

	return err
}

func (r *PostgresBillingRepository) FindInvoiceByNumber(
	invoiceNumber string,
) (*models.Invoice, error) {

	query := `
		SELECT
			id,
			invoice_number,
			owner_id,
			owner_type,
			fee_type,
			session,
			amount,
			status,
			due_date,
			created_at,
			updated_at
		FROM invoices
		WHERE invoice_number = $1
	`

	invoice := &models.Invoice{}

	err := r.db.QueryRow(
		query,
		invoiceNumber,
	).Scan(
		&invoice.ID,
		&invoice.InvoiceNumber,
		&invoice.OwnerID,
		&invoice.OwnerType,
		&invoice.FeeType,
		&invoice.Session,
		&invoice.Amount,
		&invoice.Status,
		&invoice.DueDate,
		&invoice.CreatedAt,
		&invoice.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return invoice, nil
}

func (r *PostgresBillingRepository) FindInvoicesByOwner(
	ownerID string,
	ownerType string,
) ([]models.Invoice, error) {

	query := `
		SELECT
			id,
			invoice_number,
			owner_id,
			owner_type,
			fee_type,
			session,
			amount,
			status,
			due_date,
			created_at,
			updated_at
		FROM invoices
		WHERE owner_id = $1
		AND owner_type = $2
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(
		query,
		ownerID,
		ownerType,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []models.Invoice

	for rows.Next() {

		var invoice models.Invoice

		err := rows.Scan(
			&invoice.ID,
			&invoice.InvoiceNumber,
			&invoice.OwnerID,
			&invoice.OwnerType,
			&invoice.FeeType,
			&invoice.Session,
			&invoice.Amount,
			&invoice.Status,
			&invoice.DueDate,
			&invoice.CreatedAt,
			&invoice.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		invoices = append(invoices, invoice)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return invoices, nil
}

func (r *PostgresBillingRepository) FindOutstandingInvoices(
	ownerID string,
	ownerType string,
) ([]models.Invoice, error) {

	query := `
		SELECT
			id,
			invoice_number,
			owner_id,
			owner_type,
			fee_type,
			session,
			amount,
			status,
			due_date,
			created_at,
			updated_at
		FROM invoices
		WHERE owner_id = $1
		AND owner_type = $2
		AND status = $3
		ORDER BY due_date ASC
	`

	rows, err := r.db.Query(
		query,
		ownerID,
		ownerType,
		models.InvoicePending,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []models.Invoice

	for rows.Next() {

		var invoice models.Invoice

		err := rows.Scan(
			&invoice.ID,
			&invoice.InvoiceNumber,
			&invoice.OwnerID,
			&invoice.OwnerType,
			&invoice.FeeType,
			&invoice.Session,
			&invoice.Amount,
			&invoice.Status,
			&invoice.DueDate,
			&invoice.CreatedAt,
			&invoice.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		invoices = append(invoices, invoice)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return invoices, nil
}
