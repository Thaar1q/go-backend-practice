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

##### Module 4
* Clean Architecture layer separation (Entities, Use Cases, Interface Adapters, Frameworks & Drivers)
* Pure business rules and validation decoupled from HTTP frameworks
* Independent unit testing of business rules without mock servers or databases
* Structured JSON logging with file rotation (log/slog and lumberjack.v2)
* Global security middleware, panic recovery, and scoped request validation
* Centralized dependency injection, composition root assembly, and graceful shutdown

##### Module 5
* Authentication & Security architecture using JWT access tokens and database-backed refresh tokens
* Secure password hashing with bcrypt (cost 12) and constant-time comparison mitigations
* Single-use refresh token rotation and revocation lifecycle on logout
* Route authorization middleware (`RequireAuth`) injecting authenticated user claims into request context
* Defense implementations against brute-force (IP-based rate limiting) and timing-based user enumeration
* Pure unit-tested validation rules for credentials and password strength

##### Module 6
* Role-Based Access Control (RBAC) with database tables (`roles`, `permissions`, `role_permissions`)
* Data ownership tracking (`owner_id`) so users can manage their own records
* In-memory permission cache (`PermissionSet`) for instant permission checks
* Two-layer access control: route middleware checks permissions, service layer checks resource ownership
* Protection against common access control issues (IDOR, admin self-demotion, and self-deletion)
* Unit tests (`student_authz_rules_test.go`) and automated multi-role test script (`_Test/data_test_module6.ps1`)

##### Module 7
* Input validation using `validator/v10` tags (required fields, NIM format, password rules, and safe PATCH updates with `omitnil`)
* Cursor-based pagination using base64 tokens instead of page numbers (faster, prevents skipping or duplicate items when new data is added)
* Database index on `(created_at, id)` to keep pagination queries fast
* Content negotiation supporting both JSON and CSV export based on the `Accept` header (returns 406 if format not supported)
* Unified error format (`code`, `message`, `fields`) so errors are consistent and clear across the entire API
* Automated PowerShell test scripts (`_Test/data_test_module7_D.ps1`) to verify validation, pagination, CSV exports, and errors

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
3. Execute the migration scripts:
   psql -U postgres -d go_backend_practice -f migrations/001_create_students.sql
   psql -U postgres -d go_backend_practice -f migrations/002_auth.sql
   psql -U postgres -d go_backend_practice -f migrations/003_rbac.sql
   psql -U postgres -d go_backend_practice -f migrations/004_student_permissions.sql
   psql -U postgres -d go_backend_practice -f migrations/005_cursor_index.sql

##### Database Schema
```
CREATE TABLE IF NOT EXISTS roles (
    name        VARCHAR(20) PRIMARY KEY,
    description VARCHAR(150) NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS permissions (
    name        VARCHAR(50) PRIMARY KEY,
    description VARCHAR(150) NOT NULL
);

CREATE TABLE IF NOT EXISTS role_permissions (
    role_name       VARCHAR(20) NOT NULL REFERENCES roles(name) ON DELETE CASCADE,
    permission_name VARCHAR(50) NOT NULL REFERENCES permissions(name) ON DELETE CASCADE,
    PRIMARY KEY (role_name, permission_name)
);

CREATE TABLE IF NOT EXISTS students (
    id         SERIAL PRIMARY KEY,
    nim        VARCHAR(50) NOT NULL UNIQUE,
    name       VARCHAR(100) NOT NULL,
    grade      FLOAT NOT NULL,
    is_active  BOOLEAN NOT NULL DEFAULT TRUE,
    role       VARCHAR(20) NOT NULL DEFAULT 'student' REFERENCES roles(name) ON UPDATE CASCADE,
    password   TEXT NOT NULL DEFAULT '',
    owner_id   INTEGER REFERENCES students(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_students_name_lower ON students (LOWER(name));
CREATE UNIQUE INDEX IF NOT EXISTS idx_students_nim_unique ON students (LOWER(nim));
CREATE INDEX IF NOT EXISTS idx_students_is_active ON students (is_active);
CREATE INDEX IF NOT EXISTS idx_students_role ON students (role);
CREATE INDEX IF NOT EXISTS idx_students_owner_id ON students (owner_id);
CREATE INDEX IF NOT EXISTS idx_students_created_at_id_desc ON students (created_at DESC, id DESC);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id         BIGSERIAL PRIMARY KEY,
    user_id    INTEGER NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS refresh_tokens_user_id_idx ON refresh_tokens (user_id);
```

##### Environment Configuration
Create a .env file in the root of the Module 3 project by copying .env.example:
```
cp .env.example .env
```

Fill in the environment variables according to your local setup:
```
DB_USER=postgres
DB_PASSWORD=your_password
DB_HOST=localhost
DB_PORT=5432
DB_NAME=go_backend_practice
DB_SSLMODE=disable
DB_MAX_CONNS=10
LOG_LEVEL=info
JWT_SECRET=super_secret_key_at_least_32_chars
JWT_ISSUER=go_backend_practice
```

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

   * Module 4:
     ```
     cd go_module4
     go run .
     ```

   * Module 5:
     ```
     cd go_module5
     go run .
     ```

   * Module 6:
     ```
     cd go_module6
     go run .
     ```

   * Module 7:
     ```
     cd go_module7
     go run .
     ```
4. Run Unit Tests (currently for module 4 and above):
   ```
   cd go_module[x]
   go test -v ./app/service/...
   ```
---

#### API Contract
Base URL: http://localhost:3000/api/v1

| Method | Endpoint | Query / Path / Header Params | Request Body (JSON) | Possible Statuses | Response Body Example |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/auth/register` | *(None)* | `{"nim":"123","name":"User","grade":3.5,"password":"password123"}` | 201 Created<br>400 Bad Request<br>409 Conflict<br>422 Unprocessable Entity | `{"success":true,"message":"registration successful","data":{"id":1,"nim":"123","name":"User","grade":3.5,"role":"student","is_active":true,"owner_id":1,"created_at":"..."}}` |
| `POST` | `/api/v1/auth/login` | *(None)* | `{"nim":"123","password":"password123"}` | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>429 Too Many Requests | `{"success":true,"message":"login successful","data":{"access_token":"...","refresh_token":"...","token_type":"Bearer","expires_in":900}}` |
| `POST` | `/api/v1/auth/refresh` | *(None)* | `{"refresh_token":"..."}` | 200 OK<br>400 Bad Request<br>401 Unauthorized | `{"success":true,"message":"token successfully refreshed","data":{"access_token":"...","refresh_token":"...","token_type":"Bearer","expires_in":900}}` |
| `POST` | `/api/v1/auth/logout` | *(None)* | `{"refresh_token":"..."}` | 200 OK<br>400 Bad Request | `{"success":true,"message":"logout successful","data":null}` |
| `GET` | `/api/v1/auth/me` | `Authorization: Bearer <token>` | *(None)* | 200 OK<br>401 Unauthorized | `{"success":true,"message":"profile retrieved successfully","data":{"student":{"id":1,"nim":"...","role":"student",...},"permissions":["..."]}}` |
| `GET` | `/api/v1/students` | `Authorization: Bearer <token>`<br>*(Guard: `student:list` — Admin, Staff)*<br>**Headers:** `Accept: application/json` or `text/csv`<br>**Cursor Params (M7):**<br>`cursor` (string, Base64URL)<br>`limit` (int, default 10, max 100)<br>`search` (string)<br>`is_active` (bool)<br>**Offset Params (M2-M6):**<br>`page`, `limit`, `sort`, `order`, `min_grade` | *(None)* | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>403 Forbidden<br>406 Not Acceptable | **JSON (Keyset):**<br>`{"success":true,"message":"student list retrieved successfully","data":[...],"meta":{"limit":10,"next_cursor":"...","has_more":true}}`<br><br>**CSV:**<br>`id,nim,name,grade,role,is_active,created_at`<br>`1,123456789,John,3.50,student,true,2026-09-30T10:00:00Z` |
| `GET` | `/api/v1/students/:id` | `Authorization: Bearer <token>`<br>*(Guard: Ownership or `student:read:any`)*<br>`:id` (int, path) | *(None)* | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>403 Forbidden<br>404 Not Found | `{"success":true,"message":"student found","data":{"id":1,"nim":"000000001","name":"mahasiswaAA","grade":3.75,"role":"student","is_active":true,"owner_id":1,"created_at":"..."}}` |
| `POST` | `/api/v1/students` | `Authorization: Bearer <token>`<br>*(Guard: `student:create` — Admin, Staff)* | `{"nim":"000000001","name":"mahasiswaAA","grade":3.75}` | 201 Created<br>400 Bad Request<br>401 Unauthorized<br>403 Forbidden<br>409 Conflict<br>415 Unsupported Media Type<br>422 Unprocessable Entity | `{"success":true,"message":"student successfully created","data":{...}}`<br>Header: `Location: /api/v1/students/1` |
| `PUT` | `/api/v1/students/:id` | `Authorization: Bearer <token>`<br>*(Guard: Ownership or `student:update:any`)*<br>`:id` (int, path) | `{"nim":"000000001","name":"mahasiswaAA_baru","grade":3.80,"is_active":false}` | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>403 Forbidden<br>404 Not Found<br>409 Conflict<br>415 Unsupported Media Type<br>422 Unprocessable Entity | `{"success":true,"message":"student successfully replaced","data":{...}}` |
| `PATCH` | `/api/v1/students/:id` | `Authorization: Bearer <token>`<br>*(Guard: Ownership or `student:update:any`)*<br>`:id` (int, path) | `{"grade":3.90}` | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>403 Forbidden<br>404 Not Found<br>409 Conflict<br>415 Unsupported Media Type<br>422 Unprocessable Entity | `{"success":true,"message":"student successfully partially updated","data":{...}}` |
| `PATCH` | `/api/v1/students/:id/role` | `Authorization: Bearer <token>`<br>*(Guard: `role:assign` — Admin only, Non-Self)*<br>`:id` (int, path) | `{"role":"staff"}` | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>403 Forbidden<br>404 Not Found<br>422 Unprocessable Entity | `{"success":true,"message":"student role successfully updated","data":{"id":2,"nim":"...","role":"staff",...}}` |
| `DELETE` | `/api/v1/students/:id` | `Authorization: Bearer <token>`<br>*(Guard: `student:delete` — Admin only, Non-Self)*<br>`:id` (int, path) | *(None)* | 204 No Content<br>400 Bad Request<br>401 Unauthorized<br>403 Forbidden<br>404 Not Found | *(Empty Body)* |

##### Standardized Error Envelope
```json
{
  "success": false,
  "code": "VALIDATION_ERROR",
  "message": "validasi gagal",
  "fields": {
    "nim": "must be 9 to 18 numeric digits"
  },
  "request_id": "c1a2..."
}
```
Machine-readable codes: `VALIDATION_ERROR` (422), `BAD_REQUEST` (400), `UNAUTHORIZED` (401), `FORBIDDEN` (403), `NOT_FOUND` (404), `CONFLICT` (409), `UNSUPPORTED_MEDIA_TYPE` (415), `NOT_ACCEPTABLE` (406), `TOO_MANY_REQUESTS` (429), `SERVICE_UNAVAILABLE` (503), `INTERNAL_ERROR` (500).

---

#### References Used
* https://youtu.be/8uiZC0l4Ajw
* https://docs.gofiber.io/
* https://pkg.go.dev/github.com/jackc/pgx/v5
* https://www.postgresql.org/docs/
* https://developer.mozilla.org/en-US/docs/Web/HTTP
* https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html
* https://pkg.go.dev/gopkg.in/natefinch/lumberjack.v2
* https://pkg.go.dev/github.com/golang-jwt/jwt/v5
* https://pkg.go.dev/golang.org/x/crypto/bcrypt 
* https://owasp.org/www-project-top-ten/ 
* https://owasp.org/Top10/A01_2021-Broken_Access_Control/
* https://csrc.nist.gov/projects/role-based-access-control 
* https://github.com/go-playground/validator
* https://use-the-index-luke.com/no-offset
* https://www.rfc-editor.org/rfc/rfc9110#section-12.5.1 