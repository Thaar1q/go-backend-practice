# AI Usage Disclosure

This document outlines the scope, tools, and practices regarding the use of Artificial Intelligence (AI) assistance in the development of this project.

## 1. Tools Used

* **Large Language Models (LLM):** Gemini (Google) / GitHub Copilot
* **Development Environment:** Visual Studio Code, Windows PowerShell, Git Bash
* **Purpose:** Conceptual guidance, syntax reference, debugging assistance, architecture design, and scaffolding structure.

## 2. Scope of AI Assistance

AI tools were used for the following areas:

* **Module 1 (Foundations):**
* Syntax clarification and examples for Go pointers, slices, and map operations.
* Troubleshooting Git upstream and history synchronization errors (`--allow-unrelated-histories`).
* Generating standard documentation templates (README, setup guides).


* **Module 2 (REST API & HTTP Deep Dive):**
* Struct modeling and DTO design for mutation endpoints (handling pointer fields for `PATCH` vs standard types for `PUT`).
* Debugging zero-value edge cases in JSON validation and index collision checks for unique constraints (`idx != i`).
* Query parameter parsing logic (pagination slicing, sort whitelisting, limit clamping, and filter evaluation).
* Shell-specific troubleshooting for CLI testing (PowerShell quotation escaping and `--%` stop-parsing syntax).
* Designing the API contract table and documentation structure.


## 3. Human Oversight & Verification

* **Code Review:** All generated code, structural patterns, and helper functions were reviewed, adapted, and tested manually.
* **Logic Ownership:** Core business logic, validation routines, route configurations, and HTTP status code behaviors were understood, verified, and debugged by the author.
* **Testing & Validation:** All API endpoints were manually tested via terminal HTTP requests (`curl.exe`), verifying request headers, status codes, and response envelopes against assignment specifications.
* **Integrity:** AI was utilized as a pair-programming partner and learning tool, adhering to academic and professional integrity standards.