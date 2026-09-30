-- Student permissions, role mappings, and student ownership schema
INSERT INTO permissions (name, description) VALUES
    ('student:list',       'Melihat daftar seluruh student'),
    ('student:read:any',   'Melihat data student mana pun'),
    ('student:update:any', 'Mengubah data student mana pun'),
    ('student:create',     'Menambahkan student baru'),
    ('student:delete',     'Menghapus student'),
    ('role:assign',        'Mengubah role milik student lain')
ON CONFLICT (name) DO NOTHING;

INSERT INTO role_permissions (role_name, permission_name) VALUES
    ('admin', 'student:list'),
    ('admin', 'student:read:any'),
    ('admin', 'student:create'),
    ('admin', 'student:update:any'),
    ('admin', 'student:delete'),
    ('admin', 'role:assign'),
    ('staff', 'student:list'),
    ('staff', 'student:read:any'),
    ('staff', 'student:create')
ON CONFLICT DO NOTHING;

-- Fix existing student roles and bind foreign key
UPDATE students SET role = 'student' WHERE role NOT IN (SELECT name FROM roles);

ALTER TABLE students DROP CONSTRAINT IF EXISTS fkey_students_role;
ALTER TABLE students
    ADD CONSTRAINT fkey_students_role
    FOREIGN KEY (role) REFERENCES roles(name) ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_students_role ON students (role);

-- Add owner_id and handle old rows safely
ALTER TABLE students ADD COLUMN IF NOT EXISTS owner_id INTEGER;

-- Backfill old rows
UPDATE students SET owner_id = id WHERE owner_id IS NULL;

-- Bind foreign key to students(id)
ALTER TABLE students DROP CONSTRAINT IF EXISTS fkey_students_owner;
ALTER TABLE students
    ADD CONSTRAINT fkey_students_owner
    FOREIGN KEY (owner_id) REFERENCES students(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_students_owner_id ON students (owner_id);
