-- ============================================================================
-- CAMPUSCORE
-- Migration: 000001_create_enums.down.sql
-- Description: Drop all PostgreSQL ENUM types
-- ============================================================================

DROP TYPE IF EXISTS ticket_category;

DROP TYPE IF EXISTS ticket_status;

DROP TYPE IF EXISTS clearance_status;

DROP TYPE IF EXISTS fee_type;

DROP TYPE IF EXISTS payment_status;

DROP TYPE IF EXISTS result_status;

DROP TYPE IF EXISTS academic_semester;

DROP TYPE IF EXISTS user_role;