package repository

import "context"

// UserRepository defines the contract for user data management
type UserRepository interface {
	// Add your existing user methods here if needed
}

// GovernanceRepository defines the contract for institutional rules/approvals
type GovernanceRepository interface {
	// Add your existing governance methods here if needed
}

// FinancialRepository defines the contract for processing transactions and clearance
type FinancialRepository interface {
	SavePayment(ctx context.Context, id, studentID string, amount float64, status, reference string) error
	FetchStatus(ctx context.Context, id string) (string, error)
}
