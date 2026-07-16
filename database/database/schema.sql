-- ============================================================================
-- CAMPUSCORE DATABASE SCHEMA
-- PostgreSQL 16+
-- Part 1 - Core Foundation
-- ============================================================================

-- ============================================================================
-- EXTENSIONS
-- ============================================================================

CREATE EXTENSION IF NOT EXISTS citext;

-- ============================================================================
-- ENUMS
-- ============================================================================

CREATE TYPE user_role AS ENUM (
    'student',
    'lecturer',
    'admin',
    'bursar',
    'librarian'
);

CREATE TYPE result_status AS ENUM (
    'submitted',
    'hod_approved',
    'dean_approved',
    'senate_approved',
    'finalized'
);

CREATE TYPE ticket_category AS ENUM (
    'fees',
    'registration',
    'biodata',
    'general'
);

CREATE TYPE ticket_status AS ENUM (
    'open',
    'resolved'
);

CREATE TYPE clearance_status AS ENUM (
    'pending',
    'submitted',
    'cleared'
);

CREATE TYPE payment_status AS ENUM (
    'pending',
    'successful',
    'failed'
);

CREATE TYPE fee_type AS ENUM (
    'school_fees',
    'acceptance_fee',
    'hostel_fee',
    'medical_fee',
    'late_registration'
);

-- ============================================================================
-- ACADEMIC SESSIONS
-- ============================================================================

CREATE TABLE academic_sessions (
    id SERIAL PRIMARY KEY,

    session_name VARCHAR(20) NOT NULL,
    semester VARCHAR(20) NOT NULL,

    registration_open BOOLEAN NOT NULL DEFAULT TRUE,
    result_submission_open BOOLEAN NOT NULL DEFAULT FALSE,
    clearance_open BOOLEAN NOT NULL DEFAULT FALSE,

    is_active BOOLEAN NOT NULL DEFAULT FALSE,

    starts_at DATE,
    ends_at DATE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(session_name, semester)
);

-- ============================================================================
-- DEPARTMENTS
-- ============================================================================

CREATE TABLE departments (
    id SERIAL PRIMARY KEY,

    name VARCHAR(100) NOT NULL UNIQUE,
    faculty VARCHAR(100) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    deleted_at TIMESTAMPTZ
);

-- ============================================================================
-- USERS
-- ============================================================================

CREATE TABLE users (
    id VARCHAR(50) PRIMARY KEY,

    surname VARCHAR(50) NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    middle_name VARCHAR(50),

    email CITEXT NOT NULL UNIQUE,

    phone VARCHAR(20),

    password_hash VARCHAR(255) NOT NULL,

    role user_role NOT NULL,

    department_id INT
        REFERENCES departments(id),

    level INT NOT NULL DEFAULT 100
        CHECK (level IN (100,200,300,400,500,600)),

    last_login TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    deleted_at TIMESTAMPTZ
);

-- ============================================================================
-- COURSES
-- ============================================================================

CREATE TABLE courses (
    code VARCHAR(10) PRIMARY KEY,

    title VARCHAR(150) NOT NULL,

    credit_units INT NOT NULL
        CHECK (credit_units > 0),

    level INT NOT NULL
        CHECK (level IN (100,200,300,400,500,600)),

    department_id INT NOT NULL
        REFERENCES departments(id),

    max_capacity INT NOT NULL
        CHECK (max_capacity > 0),

    current_enrolled INT NOT NULL DEFAULT 0
        CHECK (current_enrolled >= 0),

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    deleted_at TIMESTAMPTZ,

    CHECK (current_enrolled <= max_capacity)
);

-- ============================================================================
-- COURSE PREREQUISITES
-- ============================================================================

CREATE TABLE course_prerequisites (
    course_code VARCHAR(10)
        REFERENCES courses(code)
        ON DELETE CASCADE,

    prerequisite_code VARCHAR(10)
        REFERENCES courses(code)
        ON DELETE CASCADE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (
        course_code,
        prerequisite_code
    ),

    CHECK (course_code <> prerequisite_code)
);

-- ============================================================================
-- STUDENT COURSE REGISTRATION
-- ============================================================================

CREATE TABLE student_courses (
    id BIGSERIAL PRIMARY KEY,

    student_id VARCHAR(50) NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    course_code VARCHAR(10) NOT NULL
        REFERENCES courses(code)
        ON DELETE CASCADE,

    session_id INT NOT NULL
        REFERENCES academic_sessions(id),

    status VARCHAR(20) NOT NULL DEFAULT 'approved',

    registered_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (
        student_id,
        course_code,
        session_id
    )
);

-- ============================================================================
-- RESULTS
-- ============================================================================

CREATE TABLE results (
    id BIGSERIAL PRIMARY KEY,

    student_id VARCHAR(50) NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    course_code VARCHAR(10) NOT NULL
        REFERENCES courses(code)
        ON DELETE CASCADE,

    session_id INT NOT NULL
        REFERENCES academic_sessions(id),

    score NUMERIC(5,2)
        CHECK (
            score IS NULL
            OR (score >= 0 AND score <= 100)
        ),

    grade CHAR(2),

    gp NUMERIC(3,2)
        CHECK (
            gp IS NULL
            OR (gp >= 0 AND gp <= 5)
        ),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (
        student_id,
        course_code,
        session_id
    )
);

-- ============================================================================
-- RESULT APPROVAL WORKFLOW
-- ============================================================================

CREATE TABLE approvals (
    id BIGSERIAL PRIMARY KEY,

    course_code VARCHAR(10) NOT NULL
        REFERENCES courses(code),

    session_id INT NOT NULL
        REFERENCES academic_sessions(id),

    current_state result_status NOT NULL
        DEFAULT 'submitted',

    action_by VARCHAR(50)
        REFERENCES users(id),

    remarks TEXT,

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (
        course_code,
        session_id
    )
);

-- ============================================================================
-- APPROVAL HISTORY
-- ============================================================================

CREATE TABLE approval_history (
    id BIGSERIAL PRIMARY KEY,

    approval_id BIGINT NOT NULL
        REFERENCES approvals(id)
        ON DELETE CASCADE,

    previous_state result_status,

    new_state result_status NOT NULL,

    actor_id VARCHAR(50)
        REFERENCES users(id),

    remarks TEXT,

    acted_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- FEE STRUCTURES
-- ============================================================================

CREATE TABLE fee_structures (
    id BIGSERIAL PRIMARY KEY,

    department_id INT NOT NULL
        REFERENCES departments(id),

    level INT NOT NULL
        CHECK (level IN (100,200,300,400,500,600)),

    fee_type fee_type NOT NULL,

    amount_required NUMERIC(12,2) NOT NULL
        CHECK (amount_required >= 0),

    session_id INT NOT NULL
        REFERENCES academic_sessions(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (
        department_id,
        level,
        fee_type,
        session_id
    )
);

-- ============================================================================
-- FEE PAYMENTS
-- ============================================================================

CREATE TABLE fee_payments (
    id BIGSERIAL PRIMARY KEY,

    student_id VARCHAR(50) NOT NULL
        REFERENCES users(id),

    gateway_reference VARCHAR(100) NOT NULL UNIQUE,

    amount_paid NUMERIC(12,2) NOT NULL
        CHECK (amount_paid > 0),

    fee_type fee_type NOT NULL,

    session_id INT NOT NULL
        REFERENCES academic_sessions(id),

    status payment_status NOT NULL,

    gateway_name VARCHAR(50),

    gateway_transaction_id VARCHAR(150),

    payment_channel VARCHAR(50),

    currency CHAR(3) NOT NULL DEFAULT 'NGN',

    paid_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- PAYMENT WEBHOOK EVENTS
-- ============================================================================

CREATE TABLE payment_webhooks (
    id BIGSERIAL PRIMARY KEY,

    gateway_reference VARCHAR(100) NOT NULL,

    payload JSONB NOT NULL,

    received_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    processed BOOLEAN NOT NULL DEFAULT FALSE
);

-- ============================================================================
-- CLEARANCE OFFICES
-- ============================================================================

CREATE TABLE clearance_offices (
    id SERIAL PRIMARY KEY,

    office_name VARCHAR(100) NOT NULL UNIQUE,

    office_code VARCHAR(20) NOT NULL UNIQUE,

    office_order INT NOT NULL
        CHECK (office_order > 0),

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- STUDENT CLEARANCE
-- ============================================================================

CREATE TABLE student_clearances (
    id BIGSERIAL PRIMARY KEY,

    student_id VARCHAR(50) NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    office_id INT NOT NULL
        REFERENCES clearance_offices(id)
        ON DELETE CASCADE,

    status clearance_status NOT NULL
        DEFAULT 'pending',

    assigned_staff_id VARCHAR(50)
        REFERENCES users(id),

    remarks TEXT,

    signed_off_at TIMESTAMPTZ,

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

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

    original_filename VARCHAR(255) NOT NULL,

    stored_filename VARCHAR(255) NOT NULL,

    file_path TEXT NOT NULL,

    file_hash VARCHAR(64) NOT NULL,

    mime_type VARCHAR(100),

    uploaded_by VARCHAR(50)
        REFERENCES users(id),

    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- CLEARANCE HISTORY
-- ============================================================================

CREATE TABLE clearance_history (
    id BIGSERIAL PRIMARY KEY,

    clearance_id BIGINT NOT NULL
        REFERENCES student_clearances(id)
        ON DELETE CASCADE,

    previous_status clearance_status,

    new_status clearance_status NOT NULL,

    actor_id VARCHAR(50)
        REFERENCES users(id),

    remarks TEXT,

    acted_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- SUPPORT TICKETS
-- ============================================================================

CREATE TABLE support_tickets (
    id BIGSERIAL PRIMARY KEY,

    student_id VARCHAR(50) NOT NULL
        REFERENCES users(id),

    category ticket_category NOT NULL,

    status ticket_status NOT NULL
        DEFAULT 'open',

    priority VARCHAR(20) NOT NULL
        DEFAULT 'medium'
        CHECK (priority IN ('low','medium','high','urgent')),

    subject VARCHAR(150) NOT NULL,

    message TEXT NOT NULL,

    assigned_to VARCHAR(50)
        REFERENCES users(id),

    resolved_by VARCHAR(50)
        REFERENCES users(id),

    resolved_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- SUPPORT TICKET COMMENTS
-- ============================================================================

CREATE TABLE support_ticket_comments (
    id BIGSERIAL PRIMARY KEY,

    ticket_id BIGINT NOT NULL
        REFERENCES support_tickets(id)
        ON DELETE CASCADE,

    author_id VARCHAR(50)
        REFERENCES users(id),

    comment TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- AUDIT TRAIL
-- ============================================================================

CREATE TABLE audit_trails (
    id BIGSERIAL PRIMARY KEY,

    actor_id VARCHAR(50)
        REFERENCES users(id),

    action_type VARCHAR(100) NOT NULL,

    target_table VARCHAR(100),

    target_id VARCHAR(100),

    details JSONB,

    ip_address INET,

    user_agent TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- SYSTEM CONFIGURATION
-- ============================================================================

CREATE TABLE system_configs (
    config_key VARCHAR(100) PRIMARY KEY,

    config_value TEXT NOT NULL,

    description TEXT,

    updated_by VARCHAR(50)
        REFERENCES users(id),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);