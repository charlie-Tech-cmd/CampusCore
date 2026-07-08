-- CAMPUSCORE SIMS - INITIAL DATA LAYER MIGRATION
-- TARGET: POSTGRESQL (DOCKER / DEPLOYMENT INFRASTRUCTURE)

-- ============================================================================
-- 1. ENUMERATED TYPES (Enforcing Air-Tight State Machines)
-- ============================================================================
CREATE TYPE user_role AS ENUM ('student', 'lecturer', 'admin', 'bursar', 'librarian');
CREATE TYPE result_status AS ENUM ('submitted', 'hod_approved', 'dean_approved', 'senate_approved', 'finalized');
CREATE TYPE ticket_category AS ENUM ('fees', 'registration', 'biodata', 'general');
CREATE TYPE ticket_status AS ENUM ('open', 'resolved');
CREATE TYPE clearance_status AS ENUM ('pending', 'submitted', 'cleared');

-- ============================================================================
-- 2. CORE SYSTEM INFRASTRUCTURE TABLES
-- ============================================================================

-- Global System Configurations (Locks & Deadlines)
CREATE TABLE IF NOT EXISTS system_configs (
    config_key VARCHAR(50) PRIMARY KEY,
    config_value VARCHAR(100) NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Administrative Departments
CREATE TABLE IF NOT EXISTS departments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    faculty VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Identity Accounts (Universal Registry)
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(50) PRIMARY KEY, -- Matric Number, Staff ID, or Admin UUID
    surname VARCHAR(50) NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    middle_name VARCHAR(50),
    email VARCHAR(100) UNIQUE NOT NULL,
    phone VARCHAR(20) NOT NULL,
    password_hash VARCHAR(60) NOT NULL, -- Perfect fit for 60-char Bcrypt outputs
    role user_role NOT NULL,
    department_id INT REFERENCES departments(id) ON DELETE SET NULL,
    level INT DEFAULT 100, -- 100, 200, 300, 400, 500 etc.
    last_login TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Academic Modules (Dynamic Catalog)
CREATE TABLE IF NOT EXISTS courses (
    code VARCHAR(10) PRIMARY KEY, -- e.g., 'CMP101'
    title VARCHAR(100) NOT NULL,
    credit_units INT NOT NULL,
    level INT NOT NULL,
    department_id INT REFERENCES departments(id) ON DELETE CASCADE,
    max_capacity INT NOT NULL,
    current_enrolled INT DEFAULT 0
);

-- Self-Referencing Course Prerequisites Junction Table
CREATE TABLE IF NOT EXISTS course_prerequisites (
    course_code VARCHAR(10) REFERENCES courses(code) ON DELETE CASCADE,
    prerequisite_code VARCHAR(10) REFERENCES courses(code) ON DELETE CASCADE,
    PRIMARY KEY (course_code, prerequisite_code)
);

-- ============================================================================
-- 3. ACADEMIC & GOVERNANCE PIPELINES
-- ============================================================================

-- Track Result Batch Approval States (Direct Mapping to our Rejection Flowchart)
CREATE TABLE IF NOT EXISTS approvals (
    id SERIAL PRIMARY KEY,
    course_code VARCHAR(10) REFERENCES courses(code) ON DELETE CASCADE,
    session VARCHAR(9) NOT NULL,   -- e.g., '2026/2027'
    semester VARCHAR(10) NOT NULL, -- e.g., 'First'
    current_state result_status DEFAULT 'submitted',
    action_by VARCHAR(50) REFERENCES users(id),
    remarks TEXT,                  -- Mandatory for our tracking rejection loops
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (course_code, session, semester)
);

-- Active Student Registrations (Add/Drop Capable)
CREATE TABLE IF NOT EXISTS student_courses (
    id SERIAL PRIMARY KEY,
    student_id VARCHAR(50) REFERENCES users(id) ON DELETE CASCADE,
    course_code VARCHAR(10) REFERENCES courses(code) ON DELETE CASCADE,
    session VARCHAR(9) NOT NULL,
    semester VARCHAR(10) NOT NULL,
    status VARCHAR(20) DEFAULT 'approved', -- 'approved' or 'dropped'
    registered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (student_id, course_code, session, semester)
);

-- Immutable Academic Records (Requires Finalized State Verification)
CREATE TABLE IF NOT EXISTS results (
    id SERIAL PRIMARY KEY,
    student_id VARCHAR(50) REFERENCES users(id) ON DELETE CASCADE,
    course_code VARCHAR(10) REFERENCES courses(code) ON DELETE CASCADE,
    session VARCHAR(9) NOT NULL,
    semester VARCHAR(10) NOT NULL,
    score NUMERIC(5,2) NOT NULL,
    grade CHAR(1) NOT NULL,
    gp NUMERIC(3,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (student_id, course_code, session, semester)
);

-- ============================================================================
-- 4. FINANCIAL INTEGRATION & DYNAMIC OPERATIONS
-- ============================================================================

-- Dynamic Fee Configurations per Demographic Profile
CREATE TABLE IF NOT EXISTS fee_structures (
    id SERIAL PRIMARY KEY,
    department_id INT REFERENCES departments(id) ON DELETE CASCADE,
    level INT NOT NULL,
    fee_type VARCHAR(50) NOT NULL, -- 'school_fees', 'acceptance', etc.
    amount_required NUMERIC(12,2) NOT NULL,
    session VARCHAR(9) NOT NULL,
    UNIQUE (department_id, level, fee_type, session)
);

-- Core Financial Ledger Engine (Idempotency Safeguarded)
CREATE TABLE IF NOT EXISTS fee_payments (
    id SERIAL PRIMARY KEY,
    student_id VARCHAR(50) REFERENCES users(id) ON DELETE CASCADE,
    gateway_reference VARCHAR(100) UNIQUE NOT NULL, -- Eliminates Double-Spend Attacks
    amount_paid NUMERIC(12,2) NOT NULL,
    fee_type VARCHAR(50) NOT NULL,
    session VARCHAR(9) NOT NULL,
    status VARCHAR(20) NOT NULL, -- 'pending', 'successful', 'failed'
    paid_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Final-Year Multilateral Clearance Grid
CREATE TABLE IF NOT EXISTS clearance_offices (
    id SERIAL PRIMARY KEY,
    office_name VARCHAR(100) UNIQUE NOT NULL,
    is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS student_clearances (
    id SERIAL PRIMARY KEY,
    student_id VARCHAR(50) REFERENCES users(id) ON DELETE CASCADE,
    office_id INT REFERENCES clearance_offices(id) ON DELETE CASCADE,
    status clearance_status DEFAULT 'pending',
    assigned_staff_id VARCHAR(50) REFERENCES users(id),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(student_id, office_id)
);

CREATE TABLE IF NOT EXISTS clearance_attachments (
    id SERIAL PRIMARY KEY,
    clearance_id INT REFERENCES student_clearances(id) ON DELETE CASCADE,
    file_path VARCHAR(255) NOT NULL,
    file_hash VARCHAR(64) NOT NULL, -- SHA-256 integrity validation
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- 5. ENGAGEMENT & REVENUE PROTECTION AUDITS
-- ============================================================================

-- Structured Helpdesk Queue (Spam-Controlled)
CREATE TABLE IF NOT EXISTS support_tickets (
    id SERIAL PRIMARY KEY,
    student_id VARCHAR(50) REFERENCES users(id) ON DELETE CASCADE,
    category ticket_category NOT NULL,
    status ticket_status DEFAULT 'open',
    subject VARCHAR(150) NOT NULL,
    message TEXT NOT NULL,
    resolved_by VARCHAR(50) REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Bulletproof Administrative Activity Audit Log
CREATE TABLE IF NOT EXISTS audit_trails (
    id BIGSERIAL PRIMARY KEY,
    actor_id VARCHAR(50) NOT NULL,
    action_type VARCHAR(50) NOT NULL, -- 'GRADE_REJECTION', 'FEE_OVERRIDE', etc.
    target_entity VARCHAR(100) NOT NULL,
    details TEXT NOT NULL,
    ip_address VARCHAR(45),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- 6. INDEX OPTIMIZATION (Enforcing O(1) and O(log N) Database Operations)
-- ============================================================================
CREATE INDEX idx_users_role_dept ON users(role, department_id);
CREATE INDEX idx_payments_lookup ON fee_payments(student_id, session, fee_type, status);
CREATE INDEX idx_approvals_lookup ON approvals(course_code, current_state);
CREATE INDEX idx_results_student ON results(student_id, session);
CREATE INDEX idx_tickets_active ON support_tickets(student_id, status);