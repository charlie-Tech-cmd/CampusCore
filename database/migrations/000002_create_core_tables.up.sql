-- ============================================================================
-- CAMPUSCORE
-- Migration: 000002_create_core_tables.up.sql
-- Description: Create core system tables
-- ============================================================================

-- ============================================================================
-- DEPARTMENTS
-- ============================================================================

CREATE TABLE departments (
    id              BIGSERIAL PRIMARY KEY,
    name            VARCHAR(100) NOT NULL UNIQUE,
    faculty         VARCHAR(100) NOT NULL,

    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- ACADEMIC SESSIONS
-- ============================================================================

CREATE TABLE academic_sessions (
    id              BIGSERIAL PRIMARY KEY,

    session_name    VARCHAR(20) NOT NULL UNIQUE,
    semester        academic_semester NOT NULL,

    registration_open BOOLEAN NOT NULL DEFAULT FALSE,
    result_submission_open BOOLEAN NOT NULL DEFAULT FALSE,

    starts_at       DATE,
    ends_at         DATE,

    is_active       BOOLEAN NOT NULL DEFAULT FALSE,

    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- USERS
-- ============================================================================

CREATE TABLE users (
    id                  VARCHAR(50) PRIMARY KEY,

    surname             VARCHAR(50) NOT NULL,
    first_name          VARCHAR(50) NOT NULL,
    middle_name         VARCHAR(50),

    email               VARCHAR(100) NOT NULL UNIQUE,
    phone               VARCHAR(20),

    password_hash       VARCHAR(255) NOT NULL,

    role                user_role NOT NULL,

    department_id       BIGINT
        REFERENCES departments(id)
        ON DELETE SET NULL,

    level               SMALLINT,

    is_active           BOOLEAN NOT NULL DEFAULT TRUE,

    last_login          TIMESTAMPTZ,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- SYSTEM CONFIGURATION
-- ============================================================================

CREATE TABLE system_configs (
    config_key      VARCHAR(100) PRIMARY KEY,

    config_value    TEXT NOT NULL,

    description     TEXT,

    updated_by      VARCHAR(50)
        REFERENCES users(id)
        ON DELETE SET NULL,

    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);