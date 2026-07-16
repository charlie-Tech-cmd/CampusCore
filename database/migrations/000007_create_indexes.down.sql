BEGIN;

DROP INDEX IF EXISTS idx_audit_created;
DROP INDEX IF EXISTS idx_audit_action;
DROP INDEX IF EXISTS idx_audit_actor;

DROP INDEX IF EXISTS idx_ticket_resolved_by;
DROP INDEX IF EXISTS idx_ticket_category;
DROP INDEX IF EXISTS idx_ticket_status;
DROP INDEX IF EXISTS idx_ticket_student;

DROP INDEX IF EXISTS idx_clearance_attachment;
DROP INDEX IF EXISTS idx_clearance_staff;
DROP INDEX IF EXISTS idx_clearance_status;
DROP INDEX IF EXISTS idx_clearance_office;
DROP INDEX IF EXISTS idx_clearance_student;

DROP INDEX IF EXISTS idx_fee_payment_session;
DROP INDEX IF EXISTS idx_fee_payment_status;
DROP INDEX IF EXISTS idx_fee_payment_student;
DROP INDEX IF EXISTS idx_fee_payment_reference;
DROP INDEX IF EXISTS idx_fee_structure_lookup;
DROP INDEX IF EXISTS idx_fee_structure_department;

DROP INDEX IF EXISTS idx_approvals_action_by;
DROP INDEX IF EXISTS idx_approvals_state;
DROP INDEX IF EXISTS idx_approvals_course;

DROP INDEX IF EXISTS idx_results_semester;
DROP INDEX IF EXISTS idx_results_session;
DROP INDEX IF EXISTS idx_results_course;
DROP INDEX IF EXISTS idx_results_student;

DROP INDEX IF EXISTS idx_student_courses_semester;
DROP INDEX IF EXISTS idx_student_courses_session;
DROP INDEX IF EXISTS idx_student_courses_course;
DROP INDEX IF EXISTS idx_student_courses_student;

DROP INDEX IF EXISTS idx_courses_level;
DROP INDEX IF EXISTS idx_courses_department;

DROP INDEX IF EXISTS idx_users_last_login;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_department;
DROP INDEX IF EXISTS idx_users_role;

COMMIT;