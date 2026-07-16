BEGIN;

DELETE FROM system_configs
WHERE config_key IN (
    'system_initialized',

    'maintenance_mode',

    'registration_open',
    'clearance_open',
    'result_submission_open',

    'support_ticketing_enabled',
    'payment_gateway_enabled',

    'password_reset_enabled',
    'account_self_registration',

    'audit_logging_enabled'
);

COMMIT;