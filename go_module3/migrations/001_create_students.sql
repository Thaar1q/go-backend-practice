CREATE TABLE IF NOT EXISTS students (
    id         SERIAL       PRIMARY KEY,
    nim        VARCHAR(50)  NOT NULL UNIQUE,
    name       VARCHAR(100) NOT NULL,
    grade      FLOAT        NOT NULL,
    is_active  BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_students_name_lower
    ON students (LOWER(name));

CREATE UNIQUE INDEX IF NOT EXISTS idx_students_nim_unique 
    ON students (LOWER(nim));

CREATE INDEX IF NOT EXISTS idx_students_is_active 
    ON students (is_active);