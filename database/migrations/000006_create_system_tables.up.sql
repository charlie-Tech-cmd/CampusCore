-- ============================================================================
-- CAMPUSCORE
-- Migration: 000006_create_system_tables.up.sql
-- Description:
-- Creates:
--   - Audit Trails
--   - System Configurations
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

-- ============================================================================
-- SYSTEM CONFIGURATIONS
-- ============================================================================

CREATE TABLE system_configs (
    config_key VARCHAR(100) PRIMARY KEY,

    config_value TEXT NOT NULL,

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;