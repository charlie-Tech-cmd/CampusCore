-- ============================================================================
-- CAMPUSCORE
-- Migration: 000003_create_academic_tables.up.sql
-- Description: Academic management tables
-- ============================================================================

-- ============================================================================
-- COURSES
-- ============================================================================

CREATE TABLE courses (
    code                    VARCHAR(10) PRIMARY KEY,

    title                   VARCHAR(150) NOT NULL,

    credit_units            SMALLINT NOT NULL
        CHECK (credit_units > 0),

    level                   SMALLINT NOT NULL
        CHECK (level >= 100),

    department_id           BIGINT NOT NULL
        REFERENCES departments(id)
        ON DELETE RESTRICT,

    max_capacity            INTEGER NOT NULL
        CHECK (max_capacity > 0),

    current_enrolled        INTEGER NOT NULL DEFAULT 0
        CHECK (current_enrolled >= 0),

    is_active               BOOLEAN NOT NULL DEFAULT TRUE,

    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- COURSE PREREQUISITES
-- ============================================================================

CREATE TABLE course_prerequisites (
    course_code             VARCHAR(10)
        REFERENCES courses(code)
        ON DELETE CASCADE,

    prerequisite_code       VARCHAR(10)
        REFERENCES courses(code)
        ON DELETE CASCADE,

    PRIMARY KEY (
        course_code,
        prerequisite_code
    ),

    CHECK (course_code <> prerequisite_code)
);

-- ============================================================================
-- STUDENT COURSE REGISTRATION
-- ============================================================================

CREATE TABLE student_courses (
    id                      BIGSERIAL PRIMARY KEY,

    student_id              VARCHAR(50) NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    course_code             VARCHAR(10) NOT NULL
        REFERENCES courses(code)
        ON DELETE CASCADE,

    session_id              BIGINT NOT NULL
        REFERENCES academic_sessions(id)
        ON DELETE RESTRICT,

    status                  VARCHAR(20) NOT NULL
        DEFAULT 'approved'
        CHECK (
            status IN (
                'approved',
                'dropped'
            )
        ),

    registered_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (
        student_id,
        course_code,
        session_id
    )
);

-- ============================================================================
-- RESULTS
-- ============================================================================

CREATE TABLE results (
    id                      BIGSERIAL PRIMARY KEY,

    student_id              VARCHAR(50) NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    course_code             VARCHAR(10) NOT NULL
        REFERENCES courses(code)
        ON DELETE CASCADE,

    session_id              BIGINT NOT NULL
        REFERENCES academic_sessions(id)
        ON DELETE RESTRICT,

    score                   NUMERIC(5,2)
        CHECK (
            score >= 0
            AND score <= 100
        ),

    grade                   CHAR(2),

    gp                      NUMERIC(3,2)
        CHECK (
            gp >= 0
            AND gp <= 5
        ),

    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (
        student_id,
        course_code,
        session_id
    )
);

-- ============================================================================
-- RESULT APPROVAL WORKFLOW
-- ============================================================================

CREATE TABLE approvals (
    id                      BIGSERIAL PRIMARY KEY,

    course_code             VARCHAR(10) NOT NULL
        REFERENCES courses(code)
        ON DELETE CASCADE,

    session_id              BIGINT NOT NULL
        REFERENCES academic_sessions(id)
        ON DELETE RESTRICT,

    current_state           result_status NOT NULL
        DEFAULT 'submitted',

    action_by               VARCHAR(50)
        REFERENCES users(id),

    remarks                 TEXT,

    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (
        course_code,
        session_id
    )
);

-- ============================================================================
-- APPROVAL HISTORY
-- ============================================================================

CREATE TABLE approval_history (
    id                      BIGSERIAL PRIMARY KEY,

    approval_id             BIGINT NOT NULL
        REFERENCES approvals(id)
        ON DELETE CASCADE,

    previous_state          result_status,

    new_state               result_status NOT NULL,

    acted_by                VARCHAR(50)
        REFERENCES users(id),

    remarks                 TEXT,

    acted_at                TIMESTAMPTZ NOT NULL DEFAULT NOW()
);