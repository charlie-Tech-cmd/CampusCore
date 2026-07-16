BEGIN;

DELETE FROM system_configs
WHERE config_key IN (

    -- Institution
    'institution_name',
    'institution_short_name',
    'institution_email',
    'institution_phone',
    'institution_address',
    'institution_website',

    -- Academic Calendar
    'current_session',
    'current_semester',
    'registration_start_date',
    'registration_end_date',

    -- Payment
    'payment_gateway',
    'payment_public_key',
    'payment_secret_key',

    -- SMTP
    'smtp_host',
    'smtp_port',
    'smtp_username',
    'smtp_password',
    'smtp_sender_email',
    'smtp_sender_name',

    -- Defaults
    'default_timezone',
    'default_currency'
);

COMMIT;