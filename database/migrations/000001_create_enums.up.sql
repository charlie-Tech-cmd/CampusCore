-- ============================================================================
-- CAMPUSCORE
-- Migration: 000001_create_enums.up.sql
-- Description: Create all PostgreSQL ENUM types
-- ============================================================================

-- ============================================================================
-- USER ROLES
-- ============================================================================

CREATE TYPE user_role AS ENUM (
    'student',
    'lecturer',
    'admin',
    'bursar',
    'librarian'
);

-- ============================================================================
-- ACADEMIC SEMESTERS
-- ============================================================================

CREATE TYPE academic_semester AS ENUM (
    'First',
    'Second',
    'Summer'
);

-- ============================================================================
-- RESULT APPROVAL WORKFLOW
-- ============================================================================

CREATE TYPE result_status AS ENUM (
    'submitted',
    'hod_approved',
    'dean_approved',
    'senate_approved',
    'finalized'
);

-- ============================================================================
-- PAYMENT STATUS
-- ============================================================================

CREATE TYPE payment_status AS ENUM (
    'pending',
    'successful',
    'failed',
    'refunded'
);

-- ============================================================================
-- FEE TYPES
-- ============================================================================

CREATE TYPE fee_type AS ENUM (
    'school_fees',
    'acceptance',
    'development',
    'hostel',
    'library',
    'graduation',
    'other'
);

-- ============================================================================
-- CLEARANCE STATUS
-- ============================================================================

CREATE TYPE clearance_status AS ENUM (
    'pending',
    'submitted',
    'cleared'
);

-- ============================================================================
-- SUPPORT TICKET STATUS
-- ============================================================================

CREATE TYPE ticket_status AS ENUM (
    'open',
    'resolved'
);

-- ============================================================================
-- SUPPORT TICKET CATEGORY
-- ============================================================================

CREATE TYPE ticket_category AS ENUM (
    'fees',
    'registration',
    'biodata',
    'general'
);