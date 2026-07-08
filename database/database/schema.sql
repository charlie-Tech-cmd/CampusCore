-- ============================================================================
-- CAMPUSCORE DATABASE SCHEMA
-- PostgreSQL 16+
-- ============================================================================

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

-- ============================================================================
-- DEPARTMENTS
-- ============================================================================

CREATE TABLE departments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    faculty VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- USERS
-- ============================================================================

CREATE TABLE users (
    id VARCHAR(50) PRIMARY KEY,
    surname VARCHAR(50) NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    middle_name VARCHAR(50),
    email VARCHAR(100) UNIQUE NOT NULL,
    phone VARCHAR(20),
    password_hash VARCHAR(255) NOT NULL,
    role user_role NOT NULL,
    department_id INT REFERENCES departments(id),
    level INT DEFAULT 100,
    last_login TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- COURSES
-- ============================================================================

CREATE TABLE courses (
    code VARCHAR(10) PRIMARY KEY,
    title VARCHAR(150) NOT NULL,
    credit_units INT NOT NULL,
    level INT NOT NULL,
    department_id INT REFERENCES departments(id),
    max_capacity INT NOT NULL,
    current_enrolled INT DEFAULT 0
);

CREATE TABLE course_prerequisites (
    course_code VARCHAR(10) REFERENCES courses(code) ON DELETE CASCADE,
    prerequisite_code VARCHAR(10) REFERENCES courses(code) ON DELETE CASCADE,
    PRIMARY KEY(course_code, prerequisite_code)
);

CREATE TABLE student_courses (
    id SERIAL PRIMARY KEY,
    student_id VARCHAR(50) REFERENCES users(id) ON DELETE CASCADE,
    course_code VARCHAR(10) REFERENCES courses(code) ON DELETE CASCADE,
    session VARCHAR(20) NOT NULL,
    semester VARCHAR(20) NOT NULL,
    status VARCHAR(20) DEFAULT 'approved',
    registered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(student_id, course_code, session, semester)
);

-- ============================================================================
-- RESULTS
-- ============================================================================

CREATE TABLE results (
    id SERIAL PRIMARY KEY,
    student_id VARCHAR(50) REFERENCES users(id) ON DELETE CASCADE,
    course_code VARCHAR(10) REFERENCES courses(code) ON DELETE CASCADE,
    session VARCHAR(20) NOT NULL,
    semester VARCHAR(20) NOT NULL,
    score NUMERIC(5,2),
    grade CHAR(2),
    gp NUMERIC(3,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE approvals (
    id SERIAL PRIMARY KEY,
    course_code VARCHAR(10) REFERENCES courses(code),
    session VARCHAR(20),
    semester VARCHAR(20),
    current_state result_status DEFAULT 'submitted',
    action_by VARCHAR(50) REFERENCES users(id),
    remarks TEXT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(course_code, session, semester)
);

-- ============================================================================
-- FEES
-- ============================================================================

CREATE TABLE fee_structures (
    id SERIAL PRIMARY KEY,
    department_id INT REFERENCES departments(id),
    level INT NOT NULL,
    fee_type VARCHAR(50) NOT NULL,
    amount_required NUMERIC(12,2) NOT NULL,
    session VARCHAR(20) NOT NULL,
    UNIQUE(department_id, level, fee_type, session)
);

CREATE TABLE fee_payments (
    id SERIAL PRIMARY KEY,
    student_id VARCHAR(50) REFERENCES users(id),
    gateway_reference VARCHAR(100) UNIQUE NOT NULL,
    amount_paid NUMERIC(12,2) NOT NULL,
    fee_type VARCHAR(50) NOT NULL,
    session VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL,
    paid_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- CLEARANCE
-- ============================================================================

CREATE TABLE clearance_offices (
    id SERIAL PRIMARY KEY,
    office_name VARCHAR(100) UNIQUE NOT NULL,
    is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE student_clearances (
    id SERIAL PRIMARY KEY,
    student_id VARCHAR(50) REFERENCES users(id) ON DELETE CASCADE,
    office_id INT REFERENCES clearance_offices(id) ON DELETE CASCADE,
    status clearance_status DEFAULT 'pending',
    assigned_staff_id VARCHAR(50) REFERENCES users(id),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(student_id, office_id)
);

CREATE TABLE clearance_attachments (
    id SERIAL PRIMARY KEY,
    clearance_id INT REFERENCES student_clearances(id) ON DELETE CASCADE,
    file_path TEXT NOT NULL,
    file_hash VARCHAR(64) NOT NULL,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- SUPPORT
-- ============================================================================

CREATE TABLE support_tickets (
    id SERIAL PRIMARY KEY,
    student_id VARCHAR(50) REFERENCES users(id),
    category ticket_category NOT NULL,
    status ticket_status DEFAULT 'open',
    subject VARCHAR(150) NOT NULL,
    message TEXT NOT NULL,
    resolved_by VARCHAR(50) REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- AUDIT LOGS
-- ============================================================================

CREATE TABLE audit_trails (
    id BIGSERIAL PRIMARY KEY,
    actor_id VARCHAR(50) NOT NULL,
    action_type VARCHAR(100) NOT NULL,
    target_entity VARCHAR(100),
    details TEXT,
    ip_address VARCHAR(45),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- SYSTEM CONFIG
-- ============================================================================

CREATE TABLE system_configs (
    config_key VARCHAR(100) PRIMARY KEY,
    config_value TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- INDEXES
-- ============================================================================

CREATE INDEX idx_users_role
ON users(role);

CREATE INDEX idx_users_department
ON users(department_id);

CREATE INDEX idx_courses_department
ON courses(department_id);

CREATE INDEX idx_student_courses
ON student_courses(student_id);

CREATE INDEX idx_results_student
ON results(student_id);

CREATE INDEX idx_results_course
ON results(course_code);

CREATE INDEX idx_fee_payment_reference
ON fee_payments(gateway_reference);

CREATE INDEX idx_fee_payment_student
ON fee_payments(student_id);

CREATE INDEX idx_clearance_student
ON student_clearances(student_id);

CREATE INDEX idx_ticket_student
ON support_tickets(student_id);

CREATE INDEX idx_ticket_status
ON support_tickets(status);

CREATE INDEX idx_audit_actor
ON audit_trails(actor_id);

-- ============================================================================
-- DEFAULT CLEARANCE OFFICES
-- ============================================================================

INSERT INTO clearance_offices (office_name)
VALUES
('Department'),
('Faculty'),
('Library'),
('Bursary'),
('Student Affairs')
ON CONFLICT DO NOTHING;

-- ============================================================================
-- DEFAULT SYSTEM CONFIG
-- ============================================================================

INSERT INTO system_configs(config_key, config_value)
VALUES
('registration_open','true'),
('current_session','2025/2026'),
('current_semester','First')
ON CONFLICT (config_key) DO NOTHING;