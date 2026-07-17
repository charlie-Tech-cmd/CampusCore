-- ============================================================================
-- CAMPUSCORE
-- Migration: 000005_create_financial_and_support.up.sql
-- Description:
-- Creates support ticket system.
-- ============================================================================

BEGIN;

-- ============================================================================
-- SUPPORT TICKETS
-- ============================================================================

CREATE TABLE support_tickets (
    id BIGSERIAL PRIMARY KEY,

    student_id VARCHAR(50) NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    category ticket_category NOT NULL,

    status ticket_status NOT NULL
        DEFAULT 'open',

    subject VARCHAR(150) NOT NULL,

    message TEXT NOT NULL,

    resolved_by VARCHAR(50)
        REFERENCES users(id)
        ON DELETE SET NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;