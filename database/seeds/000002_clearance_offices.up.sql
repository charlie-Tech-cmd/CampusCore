BEGIN;

INSERT INTO clearance_offices (
    office_name,
    is_active
)
VALUES
    ('Department', TRUE),
    ('Faculty', TRUE),
    ('Library', TRUE),
    ('Bursary', TRUE),
    ('Student Affairs', TRUE)

ON CONFLICT (office_name)
DO UPDATE
SET
    is_active = EXCLUDED.is_active;

COMMIT;