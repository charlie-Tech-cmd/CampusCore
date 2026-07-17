-- ============================================================================
-- CAMPUSCORE
-- Migration: 000006_create_system_tables.up.sql
-- Description:
-- -- Creates:
--   - Audit Trails
-- ============================================================================

BEGIN;

-- ============================================================================
-- AUDIT TRAILS
-- ============================================================================

CREATE TABLE audit_trails (
    id BIGSERIAL PRIMARY KEY,

    actor_id VARCHAR(50) NOT NULL,

    action_type VARCHAR(100) NOT NULL,

    target_entity VARCHAR(100),

    details TEXT,

    ip_address INET,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;