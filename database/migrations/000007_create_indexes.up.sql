-- ============================================================================
-- CAMPUSCORE
-- Migration: 000007_create_indexes.up.sql
-- Description:
-- Performance indexes for production workloads.
-- ============================================================================

BEGIN;

-- ============================================================================
-- USERS
-- ============================================================================

CREATE INDEX idx_users_role
ON users(role);

CREATE INDEX idx_users_department
ON users(department_id);

CREATE INDEX idx_users_email
ON users(email);

CREATE INDEX idx_users_last_login
ON users(last_login);

-- ============================================================================
-- COURSES
-- ============================================================================

CREATE INDEX idx_courses_department
ON courses(department_id);

CREATE INDEX idx_courses_level
ON courses(level);

-- ============================================================================
-- STUDENT COURSE REGISTRATION
-- ============================================================================

CREATE INDEX idx_student_courses_student
ON student_courses(student_id);

CREATE INDEX idx_student_courses_course
ON student_courses(course_code);

CREATE INDEX idx_student_courses_session
ON student_courses(session);

CREATE INDEX idx_student_courses_semester
ON student_courses(semester);

-- ============================================================================
-- RESULTS
-- ============================================================================

CREATE INDEX idx_results_student
ON results(student_id);

CREATE INDEX idx_results_course
ON results(course_code);

CREATE INDEX idx_results_session
ON results(session);

CREATE INDEX idx_results_semester
ON results(semester);

-- ============================================================================
-- APPROVALS
-- ============================================================================

CREATE INDEX idx_approvals_course
ON approvals(course_code);

CREATE INDEX idx_approvals_state
ON approvals(current_state);

CREATE INDEX idx_approvals_action_by
ON approvals(action_by);

-- ============================================================================
-- FEES
-- ============================================================================

CREATE INDEX idx_fee_structure_department
ON fee_structures(department_id);

CREATE INDEX idx_fee_structure_lookup
ON fee_structures(department_id, level, fee_type, session);

CREATE INDEX idx_fee_payment_reference
ON fee_payments(gateway_reference);

CREATE INDEX idx_fee_payment_student
ON fee_payments(student_id);

CREATE INDEX idx_fee_payment_status
ON fee_payments(status);

CREATE INDEX idx_fee_payment_session
ON fee_payments(session);

-- ============================================================================
-- CLEARANCE
-- ============================================================================

CREATE INDEX idx_clearance_student
ON student_clearances(student_id);

CREATE INDEX idx_clearance_office
ON student_clearances(office_id);

CREATE INDEX idx_clearance_status
ON student_clearances(status);

CREATE INDEX idx_clearance_staff
ON student_clearances(assigned_staff_id);

CREATE INDEX idx_clearance_attachment
ON clearance_attachments(clearance_id);

-- ============================================================================
-- SUPPORT TICKETS
-- ============================================================================

CREATE INDEX idx_ticket_student
ON support_tickets(student_id);

CREATE INDEX idx_ticket_status
ON support_tickets(status);

CREATE INDEX idx_ticket_category
ON support_tickets(category);

CREATE INDEX idx_ticket_resolved_by
ON support_tickets(resolved_by);

-- ============================================================================
-- AUDIT
-- ============================================================================

CREATE INDEX idx_audit_actor
ON audit_trails(actor_id);

CREATE INDEX idx_audit_action
ON audit_trails(action_type);

CREATE INDEX idx_audit_created
ON audit_trails(created_at);

COMMIT;