# AI Usage Disclosure

This document outlines the scope, tools, and practices regarding the use of Artificial Intelligence (AI) assistance in the development of this project.

## 1. Tools Used

- **Large Language Models (LLM):** Gemini (Google) / GitHub Copilot
- **Development Environment:** Visual Studio Code, Windows PowerShell, Git Bash
- **Purpose:** Conceptual guidance, syntax reference, debugging assistance, architecture design, and scaffolding structure.

## 2. Scope of AI Assistance

AI tools were used for the following areas:

- **Module 1 (Foundations):**
  - Syntax clarification and examples for Go pointers, slices, and map operations.
  - Troubleshooting Git upstream and history synchronization errors (`--allow-unrelated-histories`).
  - Generating standard documentation templates (README, setup guides).

- **Module 2 (REST API & HTTP Deep Dive):**
  - Struct modeling and DTO design for mutation endpoints (handling pointer fields for `PATCH` vs standard types for `PUT`).
  - Debugging zero-value edge cases in JSON validation and index collision checks for unique constraints (`idx != i`).
  - Query parameter parsing logic (pagination slicing, sort whitelisting, limit clamping, and filter evaluation).
  - Shell-specific troubleshooting for CLI testing (PowerShell quotation escaping and `--%` stop-parsing syntax).
  - Designing the API contract table and documentation structure.

- **Module 3 (Database & Repository Pattern):**
  - Structuring the `StudentRepository` interface and implementing the concrete data layer using `pgxpool`.
  - Implementing parameterized queries (`$1, $2`) for dynamic filtering, sorting whitelists, and pagination.
  - Adding `MinGrade` filter logic to `buildFilter` to match `model.ListQuery` schema.
  - Translating driver-level errors (`pgx.ErrNoRows`, PostgreSQL code `23505`) into domain sentinel errors (`ErrNotFound`, `ErrDuplicate`).
  - Refactoring handler functions into struct methods to inject repository dependencies.
  - Refactoring handler context usage to enforce `context.WithTimeout` on all DB operations.
  - Centralizing HTTP error translation and translating Indonesian output messages to English.
  - Conceptual explanation of HTTP status codes (503 vs 500) for health checks.
  - Troubleshooting connection issues and error handling on health check endpoints.
  - Consistency check against module reference for missing implementations.
  - Formatting structural comments across all packages for code readability.

- **Module 4 (Clean Architecture):**
  - Clarified Clean Architecture layer definitions and verified folder responsibilities against the course rubric.
  - Advised on separating pure business rules (`user_rules.go`) from HTTP handlers and fixing function signatures (`ApplyPatch`).
  - Provided syntax templates for unit tests (`user_rules_test.go`) and lumberjack-based rotating loggers (`config/logger.go`).
  - Assisted in troubleshooting Windows PowerShell formatting issues, CLI escaping, and dynamic test automation scripts.
  - Formatted the layer boundary compliance audit table and ASCII dependency diagram for documentation.

- **Module 5 (Authentication & Security):**
  - Formatted SQL migration schemas for password hashing persistence and indexed token storage (`refresh_tokens`).
  - Advised on cryptographic hashing parameters using `bcrypt` (cost 12) and implementing constant-time dummy verification routines to defeat user enumeration and timing attacks.
  - Outlined the token lifecycle structure (HMAC-SHA256 JWT access tokens paired with SHA-256 hashed single-use refresh tokens in PostgreSQL).
  - Designed the `RequireAuth` bearer middleware pattern and context-binding routines (`c.Locals`).
  - Debugged DTO deserialization binding issues (`LoginRequest` tag mapping for NIM vs. Username) causing 401 Unauthorized failures during automated testing.
  - Suggested conventional commit messages and outlined the compliance test matrix for rate limiting (429) and route authorization.

## 3. Human Oversight & Verification

- **Code Review:** All generated code, structural patterns, and helper functions were reviewed, adapted, and tested manually.
- **Logic Ownership:** Core business logic, validation routines, route configurations, and HTTP status code behaviors were understood, verified, and debugged by the author.
- **Testing & Validation:** All API endpoints were manually tested via terminal HTTP requests (`curl.exe`), verifying request headers, status codes, and response envelopes against assignment specifications.
- **Integrity:** AI was utilized as a pair-programming partner and learning tool, adhering to academic and professional integrity standards.