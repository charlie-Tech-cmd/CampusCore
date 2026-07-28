BEGIN;

DELETE FROM clearance_offices
WHERE office_name IN (
    'Department',
    'Faculty',
    'Library',
    'Bursary',
    'Security',
    'Student Affairs'
    
);

COMMIT;