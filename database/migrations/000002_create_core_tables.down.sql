-- ============================================================================
-- CAMPUSCORE
-- Migration: 000002_create_core_tables.down.sql
-- Description: Drop core system tables
-- ============================================================================

DROP TABLE IF EXISTS system_configs;

DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS academic_sessions;

DROP TABLE IF EXISTS departments;