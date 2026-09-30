CREATE INDEX IF NOT EXISTS idx_students_created_at_id_desc
    ON students (created_at DESC, id DESC);