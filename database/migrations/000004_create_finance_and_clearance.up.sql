-- ============================================================================
-- CAMPUSCORE
-- Migration: 000004_create_finance_and_clearance.up.sql
-- Purpose:
--   Creates:
--     - Fee Structures
--     - Fee Payments
--     - Clearance Offices
--     - Student Clearances
--     - Clearance Attachments
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

    level INTEGER NOT NULL
        CHECK (level >= 100),

    fee_type VARCHAR(50) NOT NULL,

    amount_required NUMERIC(12,2) NOT NULL
        CHECK (amount_required >= 0),

    session VARCHAR(20) NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

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

    paid_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- CLEARANCE OFFICES
-- ============================================================================

CREATE TABLE clearance_offices (
    id BIGSERIAL PRIMARY KEY,

    office_name VARCHAR(100) NOT NULL UNIQUE,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
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

    status clearance_status
        NOT NULL
        DEFAULT 'pending',

    assigned_staff_id VARCHAR(50)
        REFERENCES users(id)
        ON DELETE SET NULL,

    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (
        student_id,
        office_id
    )
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

    uploaded_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- INDEXES
-- ============================================================================

CREATE INDEX idx_fee_structures_lookup
ON fee_structures (
    department_id,
    level,
    session
);

CREATE INDEX idx_fee_payments_student
ON fee_payments (
    student_id
);

CREATE INDEX idx_fee_payments_reference
ON fee_payments (
    gateway_reference
);

CREATE INDEX idx_fee_payments_status
ON fee_payments (
    status
);

CREATE INDEX idx_student_clearances_student
ON student_clearances (
    student_id
);

CREATE INDEX idx_student_clearances_status
ON student_clearances (
    status
);

CREATE INDEX idx_clearance_attachments_clearance
ON clearance_attachments (
    clearance_id
);

COMMIT;