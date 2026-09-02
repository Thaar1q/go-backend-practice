CREATE TABLE IF NOT EXISTS students (
    ID        SERIAL       PRIMARY KEY,
    NIM       VARCHAR(50)  NOT NULL UNIQUE,
    Name      VARCHAR(100) NOT NULL,
    Grade     FLOAT        NOT NULL,
    IsActive  BOOLEAN      NOT NULL DEFAULT TRUE,
    CreatedAt TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
 
CREATE INDEX IF NOT EXISTS idx_students_name_lower
    ON students (LOWER(Name));

CREATE UNIQUE INDEX IF NOT EXISTS idx_students_nim_unique 
    ON students (LOWER(NIM));
 
 CREATE INDEX IF NOT EXISTS idx_students_is_active 
    ON students (IsActive);