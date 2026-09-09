# AI Usage Disclosure

This document outlines the scope, tools, and practices regarding the use of Artificial Intelligence (AI) assistance in the development of this project[cite: 1].

## 1. Tools Used

- **Large Language Models (LLM):** Gemini (Google) / GitHub Copilot[cite: 1]
- **Development Environment:** Visual Studio Code, Windows PowerShell, Git Bash[cite: 1]
- **Purpose:** Conceptual guidance, syntax reference, debugging assistance, architecture design, and scaffolding structure[cite: 1].

## 2. Scope of AI Assistance

AI tools were used for the following areas:

- **Module 1 (Foundations):**
  - Syntax clarification and examples for Go pointers, slices, and map operations[cite: 1].
  - Troubleshooting Git upstream and history synchronization errors (`--allow-unrelated-histories`)[cite: 1].
  - Generating standard documentation templates (README, setup guides)[cite: 1].

- **Module 2 (REST API & HTTP Deep Dive):**
  - Struct modeling and DTO design for mutation endpoints (handling pointer fields for `PATCH` vs standard types for `PUT`)[cite: 1].
  - Debugging zero-value edge cases in JSON validation and index collision checks for unique constraints (`idx != i`)[cite: 1].
  - Query parameter parsing logic (pagination slicing, sort whitelisting, limit clamping, and filter evaluation)[cite: 1].
  - Shell-specific troubleshooting for CLI testing (PowerShell quotation escaping and `--%` stop-parsing syntax)[cite: 1].
  - Designing the API contract table and documentation structure[cite: 1].

- **Module 3 (Database & Repository Pattern):**
  - Structuring the `StudentRepository` interface and implementing the concrete data layer using `pgxpool`[cite: 1].
  - Implementing parameterized queries (`$1, $2`) for dynamic filtering, sorting whitelists, and pagination[cite: 1].
  - Adding `MinGrade` filter logic to `buildFilter` to match `model.ListQuery` schema[cite: 1].
  - Translating driver-level errors (`pgx.ErrNoRows`, PostgreSQL code `23505`) into domain sentinel errors (`ErrNotFound`, `ErrDuplicate`)[cite: 1].
  - Refactoring handler functions into struct methods to inject repository dependencies[cite: 1].
  - Refactoring handler context usage to enforce `context.WithTimeout` on all DB operations[cite: 1].
  - Centralizing HTTP error translation and translating Indonesian output messages to English[cite: 1].
  - Conceptual explanation of HTTP status codes (503 vs 500) for health checks[cite: 1].
  - Troubleshooting connection issues and error handling on health check endpoints[cite: 1].
  - Consistency check against module reference for missing implementations[cite: 1].
  - Formatting structural comments across all packages for code readability[cite: 1].

- **Module 4 (Clean Architecture):**
  - Clarified Clean Architecture layer definitions and verified folder responsibilities against the course rubric.
  - Advised on separating pure business rules (`user_rules.go`) from HTTP handlers and fixing function signatures (`ApplyPatch`).
  - Provided syntax templates for unit tests (`user_rules_test.go`) and lumberjack-based rotating loggers (`config/logger.go`).
  - Assisted in troubleshooting Windows PowerShell formatting issues, CLI escaping, and dynamic test automation scripts.
  - Formatted the layer boundary compliance audit table and ASCII dependency diagram for documentation.

## 3. Human Oversight & Verification

- **Code Review:** All generated code, structural patterns, and helper functions were reviewed, adapted, and tested manually[cite: 1].
- **Logic Ownership:** Core business logic, validation routines, route configurations, and HTTP status code behaviors were understood, verified, and debugged by the author[cite: 1].
- **Testing & Validation:** All API endpoints were manually tested via terminal HTTP requests (`curl.exe`), verifying request headers, status codes, and response envelopes against assignment specifications[cite: 1].
- **Integrity:** AI was utilized as a pair-programming partner and learning tool, adhering to academic and professional integrity standards[cite: 1].