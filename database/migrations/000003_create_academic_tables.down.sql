-- ============================================================================
-- CAMPUSCORE
-- Migration: 000003_create_academic_tables.down.sql
-- Description: Drop academic tables
-- ============================================================================

DROP TABLE IF EXISTS approval_history;

DROP TABLE IF EXISTS approvals;

DROP TABLE IF EXISTS results;

DROP TABLE IF EXISTS student_courses;

DROP TABLE IF EXISTS course_prerequisites;

DROP TABLE IF EXISTS courses;