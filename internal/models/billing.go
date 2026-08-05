package models

import "time"

type InvoiceStatus string

const (
	InvoicePending InvoiceStatus = "pending"
	InvoicePaid    InvoiceStatus = "paid"
	InvoiceExpired InvoiceStatus = "expired"
	InvoiceVoided  InvoiceStatus = "voided"
)

type Invoice struct {
	ID int `json:"id"`

	InvoiceNumber string `json:"invoice_number"`

	OwnerID string `json:"owner_id"`

	OwnerType string `json:"owner_type"`

	FeeType string `json:"fee_type"`

	Session string `json:"session"`

	Amount float64 `json:"amount"`

	Status InvoiceStatus `json:"status"`

	DueDate time.Time `json:"due_date"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`
}

type BillingRepository interface {
	CreateInvoice(
		invoice *Invoice,
	) error

	FindInvoiceByID(
		id int,
	) (*Invoice, error)

	FindOwnerInvoices(
		ownerID string,
		ownerType string,
	) ([]Invoice, error)

	UpdateInvoice(
		invoice *Invoice,
	) error

	MarkInvoicePaid(
		id int,
	) error

	FindInvoiceByNumber(
		invoiceNumber string,
	) (*Invoice, error)

	FindInvoicesByOwner(
		ownerID string,
		ownerType string,
	) ([]Invoice, error)

	FindOutstandingInvoices(
		ownerID string,
		ownerType string,
	) ([]Invoice, error)
}
