# Go Backend Programming

Ongoing project about learning Go Language.

#### Currently Available Modules

##### Module 1
* Hello World
* Variable declaration, map declaration and operations
* Pointer functions
* Struct functions

##### Module 2
* REST API implementation with Fiber v2
* In-memory CRUD operations (GET, POST, PUT, PATCH, DELETE)
* HTTP semantics, status codes, and custom response envelopes (WebResponse, Meta)
* DTO validation and mutation patterns (full replacement vs. partial update)
* Query string processing (pagination, search, sorting, filtering, and limit caps)

##### Module 3
* Persistent storage integration using PostgreSQL and pgx/v5 (pgxpool)
* Database schema migration and index optimizations
* Repository Pattern implementation using Go interfaces
* Parameterized SQL queries preventing SQL injection
* Environment variable configuration (.env and .env.example)
* Sentinel error translation from database errors to HTTP status codes

---

#### How to Initialize

##### Prerequisites
* Go (version 1.22 or newer recommended)
* Git
* PostgreSQL (version 14 or newer)
* curl / Postman

##### Database Setup
1. Ensure the PostgreSQL service is active and accessible.
2. Create the target database (e.g., go_backend_practice):
   CREATE DATABASE go_backend_practice;
3. Execute the migration script located at migrations/001_create_students.sql:
   psql -U postgres -d go_backend_practice -f migrations/001_create_students.sql

##### Database Schema
```
CREATE TABLE IF NOT EXISTS students (
    id         SERIAL PRIMARY KEY,
    nim        VARCHAR(50) NOT NULL UNIQUE,
    name       VARCHAR(100) NOT NULL,
    grade      FLOAT NOT NULL,
    is_active  BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_students_name_lower ON students (LOWER(name));
CREATE UNIQUE INDEX IF NOT EXISTS idx_students_nim_unique ON students (LOWER(nim));
CREATE INDEX IF NOT EXISTS idx_students_is_active ON students (is_active);
```

##### Environment Configuration
Create a .env file in the root of the Module 3 project by copying .env.example:
cp .env.example .env

Fill in the environment variables according to your local setup:
DB_USER=postgres
DB_PASSWORD=your_password
DB_HOST=localhost
DB_PORT=5432
DB_NAME=go_backend_practice
DB_SSLMODE=disable
DB_MAX_CONNS=10

##### Installation & Execution
1. Clone the repository:
   ```
   git clone https://github.com/Thaar1q/go-backend-practice.git
   cd go-backend-practice
   ```

2. Download and verify dependencies:
   ```
   go mod tidy
   ```

3. Run the tasks / modules:
   * Module 1:
     ```
     go run ./cmd/task1/main.go
     go run ./cmd/task2/main.go
     go run ./cmd/task3/main.go
     go run ./cmd/task4/main.go
     ```

   * Module 2:
     ```
     cd go_module2
     go run .
     ```

   * Module 3:
     ```
     cd go_module3
     go run .
     ```

---

#### API Contract
Base URL: http://localhost:3000/api/v1/students

| Method | Endpoint | Query / Path Parameters | Request Body (JSON) | Possible Statuses | Response Body Example |
| :--- | :--- | :--- | :--- | :--- | :--- |
| GET | /api/v1/students | page (int, default 1)<br>limit (int, default 10, max 100)<br>search (string, case-insensitive)<br>sort (id, nim, name, grade)<br>order (asc, desc)<br>is_active (bool) | (None) | 200 OK | {"success": true, "message": "daftar mahasiswa berhasil diambil", "data": [...], "meta": {"page": 1, "limit": 10, "total": 1, "total_pages": 1}} |
| GET | /api/v1/students/:id | :id (int, path) | (None) | 200 OK<br>400 Bad Request<br>404 Not Found | {"success": true, "message": "mahasiswa ditemukan", "data": {"id": 1, "nim": "000000001", "name": "mahasiswaAA", "grade": 3.75, "is_active": true, "created_at": "..."}} |
| POST | /api/v1/students | (None) | {"nim": "000000001", "name": "mahasiswaAA", "grade": 3.75} | 201 Created<br>400 Bad Request<br>409 Conflict<br>415 Unsupported Media Type<br>422 Unprocessable Entity | {"success": true, "message": "mahasiswa berhasil dibuat", "data": {...}}<br>Header: Location: /api/v1/students/1 |
| PUT | /api/v1/students/:id | :id (int, path) | {"nim": "000000001", "name": "mahasiswaAA_baru", "grade": 3.80, "is_active": false} (All fields mandatory) | 200 OK<br>400 Bad Request<br>404 Not Found<br>409 Conflict<br>415 Unsupported Media Type<br>422 Unprocessable Entity | {"success": true, "message": "mahasiswa berhasil diganti seluruhnya", "data": {...}} |
| PATCH | /api/v1/students/:id | :id (int, path) | {"grade": 3.90} (Only modified fields required) | 200 OK<br>400 Bad Request<br>404 Not Found<br>409 Conflict<br>415 Unsupported Media Type<br>422 Unprocessable Entity | {"success": true, "message": "mahasiswa berhasil diperbarui sebagian", "data": {...}} |
| DELETE | /api/v1/students/:id | :id (int, path) | (None) | 204 No Content<br>400 Bad Request<br>404 Not Found | (Empty Body) |

---

#### References Used
* https://www.youtube.com/watch?v=8uiZC0l4Ajw
* https://docs.gofiber.io/
* https://pkg.go.dev/github.com/jackc/pgx/v5
* https://www.postgresql.org/docs/
* https://developer.mozilla.org/en-US/docs/Web/HTTP

```