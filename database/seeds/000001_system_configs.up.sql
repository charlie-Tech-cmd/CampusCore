BEGIN;

INSERT INTO system_configs (config_key, config_value)
VALUES
    -- =========================================================================
    -- SYSTEM INITIALIZATION
    -- =========================================================================
    ('system_initialized', 'false'),

    -- =========================================================================
    -- APPLICATION STATE
    -- =========================================================================
    ('maintenance_mode', 'false'),

    -- =========================================================================
    -- FEATURE FLAGS
    -- =========================================================================
    ('registration_open', 'false'),
    ('clearance_open', 'false'),
    ('result_submission_open', 'false'),
    ('support_ticketing_enabled', 'true'),
    ('payment_gateway_enabled', 'true'),

    -- =========================================================================
    -- SECURITY
    -- =========================================================================
    ('password_reset_enabled', 'true'),
    ('account_self_registration', 'false'),

    -- =========================================================================
    -- AUDITING
    -- =========================================================================
    ('audit_logging_enabled', 'true')

ON CONFLICT (config_key)
DO NOTHING;

COMMIT;