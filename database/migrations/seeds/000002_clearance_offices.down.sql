BEGIN;

DELETE FROM clearance_offices
WHERE office_name IN (
    'Department',
    'Faculty',
    'Library',
    'Bursary',
    'Student Affairs'
);

COMMIT;