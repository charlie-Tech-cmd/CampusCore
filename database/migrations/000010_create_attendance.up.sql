-- ============================================================================
-- CAMPUSCORE
-- Migration: 000010_create_attendance.up.sql
-- Description: Student attendance management
-- ============================================================================

CREATE TABLE attendance (
    id BIGSERIAL PRIMARY KEY,

    student_id VARCHAR(50) NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    course_code VARCHAR(10) NOT NULL
        REFERENCES courses(code)
        ON DELETE CASCADE,

    lecturer_id VARCHAR(50)
        REFERENCES users(id)
        ON DELETE SET NULL,

    session VARCHAR(20) NOT NULL,

    semester VARCHAR(20) NOT NULL,

    class_date DATE NOT NULL,

    status VARCHAR(20) NOT NULL
        CHECK (status IN ('present', 'absent', 'excused')),

    marked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (
        student_id,
        course_code,
        session,
        semester,
        class_date
    )
);

CREATE INDEX idx_attendance_student
    ON attendance(student_id);

CREATE INDEX idx_attendance_course
    ON attendance(course_code);

CREATE INDEX idx_attendance_lecturer
    ON attendance(lecturer_id);

CREATE INDEX idx_attendance_date
    ON attendance(class_date);

CREATE INDEX idx_attendance_status
    ON attendance(status);