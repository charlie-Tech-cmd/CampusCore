BEGIN;

INSERT INTO system_configs
(config_key, config_value)
VALUES

-- ======================================================
-- Institution
-- ======================================================

('institution_name',''),
('institution_short_name',''),
('institution_email',''),
('institution_phone',''),
('institution_address',''),
('institution_website',''),

-- ======================================================
-- Academic Calendar
-- ======================================================

('current_session',''),
('current_semester',''),
('registration_start_date',''),
('registration_end_date',''),

-- ======================================================
-- Localization
-- ======================================================

('default_timezone','Africa/Lagos'),
('default_currency','NGN'),
('date_format','DD/MM/YYYY'),

-- ======================================================
-- Branding
-- ======================================================

('logo_url',''),
('portal_theme','default'),

ON CONFLICT (config_key)
DO NOTHING;

COMMIT;