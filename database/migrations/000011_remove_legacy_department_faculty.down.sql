-- ============================================================================
-- CAMPUSCORE
-- Migration: 000011_remove_legacy_department_faculty.down.sql
-- Description: Restore legacy departments.faculty column
-- ============================================================================

ALTER TABLE departments
ADD COLUMN faculty VARCHAR(100);

UPDATE departments d
SET faculty = f.name
FROM faculties f
WHERE d.faculty_id = f.id;

ALTER TABLE departments
ALTER COLUMN faculty SET NOT NULL;
