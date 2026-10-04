-- ============================================================================
-- CAMPUSCORE
-- Migration: 000009_create_billing_and_notifications.up.sql
-- Purpose:
--   Creates:
--     - Invoices
--     - Notifications
-- ============================================================================

BEGIN;

-- ============================================================================
-- INVOICES
-- ============================================================================

CREATE TABLE invoices (
    id BIGSERIAL PRIMARY KEY,

    invoice_number VARCHAR(100) NOT NULL UNIQUE,

    owner_id VARCHAR(50) NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    owner_type VARCHAR(30) NOT NULL,

    fee_type VARCHAR(50) NOT NULL,

    session VARCHAR(20) NOT NULL,

    amount NUMERIC(12,2) NOT NULL
        CHECK (amount > 0),

    status VARCHAR(20) NOT NULL
        CHECK (status IN ('pending', 'paid', 'expired', 'voided')),

    due_date TIMESTAMP NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_invoices_owner
    ON invoices(owner_id, owner_type);

CREATE INDEX idx_invoices_status
    ON invoices(status);

CREATE INDEX idx_invoices_due_date
    ON invoices(due_date);

-- ============================================================================
-- NOTIFICATIONS
-- ============================================================================

CREATE TABLE notifications (
    id BIGSERIAL PRIMARY KEY,

    user_id VARCHAR(50) NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    title VARCHAR(150) NOT NULL,

    message TEXT NOT NULL,

    type VARCHAR(30) NOT NULL,

    is_read BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notifications_user
    ON notifications(user_id);

CREATE INDEX idx_notifications_user_created
    ON notifications(user_id, created_at);

COMMIT;
