-- Table structures and default roles/permissions
CREATE TABLE IF NOT EXISTS roles (
    name        VARCHAR(20)  PRIMARY KEY,
    description VARCHAR(150) NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

INSERT INTO roles (name, description) VALUES
    ('admin', 'Akses penuh terhadap seluruh data dan pengaturan'),
    ('staff', 'Boleh melihat data seluruh user/student, tetapi tidak boleh mengubah'),
    ('student', 'Hanya boleh mengelola datanya sendiri')
ON CONFLICT (name) DO NOTHING;

CREATE TABLE IF NOT EXISTS permissions (
    name        VARCHAR(50)  PRIMARY KEY,
    description VARCHAR(150) NOT NULL
);

CREATE TABLE IF NOT EXISTS role_permissions (
    role_name       VARCHAR(20) NOT NULL REFERENCES roles(name) ON DELETE CASCADE,
    permission_name VARCHAR(50) NOT NULL REFERENCES permissions(name) ON DELETE CASCADE,
    PRIMARY KEY (role_name, permission_name)
);
