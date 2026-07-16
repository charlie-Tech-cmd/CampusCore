-- ============================================================================
-- CAMPUSCORE
-- Migration: 000005_create_financial_and_support.up.sql
-- Description:
-- Creates:
--   - Fee Structures
--   - Fee Payments
--   - Clearance Offices
--   - Student Clearances
--   - Clearance Attachments
--   - Support Tickets
-- ============================================================================

BEGIN;

-- ============================================================================
-- FEE STRUCTURES
-- ============================================================================

CREATE TABLE fee_structures (
    id BIGSERIAL PRIMARY KEY,

    department_id BIGINT NOT NULL
        REFERENCES departments(id)
        ON DELETE CASCADE,

    level INTEGER NOT NULL,

    fee_type VARCHAR(50) NOT NULL,

    amount_required NUMERIC(12,2) NOT NULL
        CHECK (amount_required >= 0),

    session VARCHAR(20) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (
        department_id,
        level,
        fee_type,
        session
    )
);

-- ============================================================================
-- FEE PAYMENTS
-- ============================================================================

CREATE TABLE fee_payments (
    id BIGSERIAL PRIMARY KEY,

    student_id VARCHAR(50) NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    gateway_reference VARCHAR(100) NOT NULL UNIQUE,

    amount_paid NUMERIC(12,2) NOT NULL
        CHECK (amount_paid >= 0),

    fee_type VARCHAR(50) NOT NULL,

    session VARCHAR(20) NOT NULL,

    status VARCHAR(20) NOT NULL,

    paid_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- CLEARANCE OFFICES
-- ============================================================================

CREATE TABLE clearance_offices (
    id BIGSERIAL PRIMARY KEY,

    office_name VARCHAR(100) NOT NULL UNIQUE,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- STUDENT CLEARANCES
-- ============================================================================

CREATE TABLE student_clearances (
    id BIGSERIAL PRIMARY KEY,

    student_id VARCHAR(50) NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    office_id BIGINT NOT NULL
        REFERENCES clearance_offices(id)
        ON DELETE CASCADE,

    status clearance_status NOT NULL
        DEFAULT 'pending',

    assigned_staff_id VARCHAR(50)
        REFERENCES users(id)
        ON DELETE SET NULL,

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(student_id, office_id)
);

-- ============================================================================
-- CLEARANCE ATTACHMENTS
-- ============================================================================

CREATE TABLE clearance_attachments (
    id BIGSERIAL PRIMARY KEY,

    clearance_id BIGINT NOT NULL
        REFERENCES student_clearances(id)
        ON DELETE CASCADE,

    file_path TEXT NOT NULL,

    file_hash VARCHAR(64) NOT NULL,

    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

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