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

	return nil
}

func (r *PostgresBillingRepository) FindInvoiceByID(
	id int,
) (*models.Invoice, error) {

	return nil, nil
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

	return nil
}

func (r *PostgresBillingRepository) MarkInvoicePaid(
	id int,
) error {

	return nil
}

func (r *PostgresBillingRepository) FindInvoiceByNumber(
	invoiceNumber string,
) (*models.Invoice, error) {

	return nil, nil
}

func (r *PostgresBillingRepository) FindInvoicesByOwner(
	ownerID string,
	ownerType string,
) ([]models.Invoice, error) {

	return nil, nil
}

func (r *PostgresBillingRepository) FindOutstandingInvoices(
	ownerID string,
	ownerType string,
) ([]models.Invoice, error) {

	return nil, nil
}
