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
   psql -U postgres -d go_backend_practice -f migrations/002_auth.sql

##### Database Schema
```
CREATE TABLE IF NOT EXISTS students (
    id         SERIAL PRIMARY KEY,
    nim        VARCHAR(50) NOT NULL UNIQUE,
    name       VARCHAR(100) NOT NULL,
    grade      FLOAT NOT NULL,
    is_active  BOOLEAN NOT NULL DEFAULT TRUE,
    role       VARCHAR(20) NOT NULL DEFAULT 'student',
    password   TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_students_name_lower ON students (LOWER(name));
CREATE UNIQUE INDEX IF NOT EXISTS idx_students_nim_unique ON students (LOWER(nim));
CREATE INDEX IF NOT EXISTS idx_students_is_active ON students (is_active);

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
4. Run Unit Tests (currently for module 4 and above):
   ```
   cd go_module[x]
   go test -v ./app/service/...
   ```
---

#### API Contract
Base URL: http://localhost:3000/api/v1/students

| Method | Endpoint | Query / Path / Header Params | Request Body (JSON) | Possible Statuses | Response Body Example |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/auth/register` | *(None)* | `{"nim":"123","name":"User","grade":3.5,"password":"password123"}` | 201 Created<br>400 Bad Request<br>409 Conflict<br>422 Unprocessable Entity | `{"success":true,"message":"registration successful","data":{"id":1,"nim":"123","name":"User","grade":3.5,"role":"student","is_active":true,"created_at":"..."}}`[cite: 2] |
| `POST` | `/api/v1/auth/login` | *(None)* | `{"nim":"123","password":"password123"}` | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>429 Too Many Requests | `{"success":true,"message":"login successful","data":{"access_token":"...","refresh_token":"...","token_type":"Bearer","expires_in":900}}`[cite: 2] |
| `POST` | `/api/v1/auth/refresh` | *(None)* | `{"refresh_token":"..."}` | 200 OK<br>400 Bad Request<br>401 Unauthorized | `{"success":true,"message":"token successfully refreshed","data":{"access_token":"...","refresh_token":"...","token_type":"Bearer","expires_in":900}}`[cite: 2] |
| `POST` | `/api/v1/auth/logout` | *(None)* | `{"refresh_token":"..."}` | 200 OK<br>400 Bad Request | `{"success":true,"message":"logout successful","data":null}`[cite: 2] |
| `GET` | `/api/v1/auth/me` | `Authorization: Bearer <token>` | *(None)* | 200 OK<br>401 Unauthorized | `{"success":true,"message":"user profile retrieved","data":{"user_id":1,"username":"...","role":"student"}}`[cite: 2] |
| `GET` | `/api/v1/students` | `Authorization: Bearer <token>`<br>page (int, default 1)<br>limit (int, default 10, max 100)<br>search (string)<br>sort (id, nim, name, grade)<br>order (asc, desc)<br>is_active (bool)<br>min_grade (float) | *(None)* | 200 OK<br>401 Unauthorized | `{"success":true,"message":"student list successfully retrieved","data":[...],"meta":{"page":1,"limit":10,"total":1,"total_pages":1}}`[cite: 1, 2] |
| `GET` | `/api/v1/students/:id` | `Authorization: Bearer <token>`<br>`:id` (int, path) | *(None)* | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>404 Not Found | `{"success":true,"message":"student found","data":{"id":1,"nim":"000000001","name":"mahasiswaAA","grade":3.75,"is_active":true,"created_at":"..."}}`[cite: 1, 2] |
| `POST` | `/api/v1/students` | `Authorization: Bearer <token>` | `{"nim":"000000001","name":"mahasiswaAA","grade":3.75}` | 201 Created<br>400 Bad Request<br>401 Unauthorized<br>409 Conflict<br>415 Unsupported Media Type<br>422 Unprocessable Entity | `{"success":true,"message":"student successfully created","data":{...}}`<br>Header: `Location: /api/v1/students/1`[cite: 1, 2] |
| `PUT` | `/api/v1/students/:id` | `Authorization: Bearer <token>`<br>`:id` (int, path) | `{"nim":"000000001","name":"mahasiswaAA_baru","grade":3.80,"is_active":false}` | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>404 Not Found<br>409 Conflict<br>415 Unsupported Media Type<br>422 Unprocessable Entity | `{"success":true,"message":"student successfully replaced","data":{...}}`[cite: 1, 2] |
| `PATCH` | `/api/v1/students/:id` | `Authorization: Bearer <token>`<br>`:id` (int, path) | `{"grade":3.90}` | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>404 Not Found<br>409 Conflict<br>415 Unsupported Media Type<br>422 Unprocessable Entity | `{"success":true,"message":"student successfully partially updated","data":{...}}`[cite: 1, 2] |
| `DELETE` | `/api/v1/students/:id` | `Authorization: Bearer <token>`<br>`:id` (int, path) | *(None)* | 204 No Content<br>400 Bad Request<br>401 Unauthorized<br>404 Not Found | *(Empty Body)*[cite: 1, 2] |
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