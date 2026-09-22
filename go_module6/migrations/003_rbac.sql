-- New Tables and Values
CREATE TABLE IF NOT EXISTS roles (
    name        VARCHAR(20)  PRIMARY KEY,
    description VARCHAR(150) NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
 
INSERT INTO roles (name, description) VALUES
    ('admin', 'Akses penuh terhadap seluruh data dan pengaturan'),
    ('staff', 'Boleh melihat data seluruh student, tetapi tidak boleh mengubah'),
    ('student',  'Hanya boleh mengelola datanya sendiri')
ON CONFLICT (name) DO NOTHING;
 
CREATE TABLE IF NOT EXISTS permissions (
    name        VARCHAR(50)  PRIMARY KEY,
    description VARCHAR(150) NOT NULL
);

INSERT INTO permissions (name, description) VALUES
    ('student:list',       'Melihat daftar seluruh student'),
    ('student:read:any',   'Melihat data student mana pun'),
    ('student:update:any', 'Mengubah data student mana pun'),
    ('student:delete',     'Menghapus student'),
    ('role:assign',      'Mengubah role milik student lain')
ON CONFLICT (name) DO NOTHING;
 
CREATE TABLE IF NOT EXISTS role_permissions (
    role_name       VARCHAR(20) NOT NULL REFERENCES roles(name) ON DELETE CASCADE,
    permission_name VARCHAR(50) NOT NULL REFERENCES permissions(name) ON DELETE CASCADE,
    PRIMARY KEY (role_name, permission_name)
);
 
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

-- Fix Existing Roles & Bind Foreign Key
UPDATE students SET role = 'student' WHERE role NOT IN (SELECT name FROM roles);
 
ALTER TABLE students DROP CONSTRAINT IF EXISTS fkey_students_role;
ALTER TABLE students
    ADD CONSTRAINT fkey_students_role
    FOREIGN KEY (role) REFERENCES roles(name) ON UPDATE CASCADE;
 
CREATE INDEX IF NOT EXISTS idx_students_role ON students (role);

-- Add owner_id & handle old rows safely
ALTER TABLE students ADD COLUMN IF NOT EXISTS owner_id INTEGER;

-- Backfill old rows
UPDATE students SET owner_id = id WHERE owner_id IS NULL;

-- Bind foreign key to students(id)
ALTER TABLE students DROP CONSTRAINT IF EXISTS fkey_students_owner;
ALTER TABLE students
    ADD CONSTRAINT fkey_students_owner
    FOREIGN KEY (owner_id) REFERENCES students(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_students_owner_id ON students (owner_id);
