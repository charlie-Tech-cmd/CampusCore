-- ============================================================================
-- CAMPUSCORE
-- Migration: 000011_remove_legacy_department_faculty.up.sql
-- Description: Remove legacy departments.faculty column
-- ============================================================================

ALTER TABLE departments
DROP COLUMN faculty;
