# Go Backend Programming

Ongoing project about learning Go Language.

## Currently Available Modules

### Module 1

* Hello World
* Variable declaration, map declaration and operations
* Pointer functions
* Struct functions

### Module 2

* REST API implementation with Fiber v2
* In-memory CRUD operations (`GET`, `POST`, `PUT`, `PATCH`, `DELETE`)
* HTTP semantics, status codes, and custom response envelopes (`WebResponse`, `Meta`)
* DTO validation and mutation patterns (full replacement vs. partial update)
* Query string processing (pagination, search, sorting, filtering, and limit caps)

---

## How to Initialize

### Prerequisites

* [Go](https://go.dev/dl/) (version 1.20 or newer recommended)
* Git

### Installation

1. Clone the repository:
```bash
git clone https://github.com/Thaar1q/go-backend-practice.git
cd go-backend-practice

```


2. Download and verify dependencies:
```bash
go mod tidy

```


3. Run the tasks / modules:
* **Module 1**:
```bash
go run ./cmd/task1/main.go
go run ./cmd/task2/main.go
go run ./cmd/task3/main.go
go run ./cmd/task4/main.go

```


* **Module 2**:
```bash
cd api-students
go run .

```





---

## API Contract (Module 2)

Base URL: `http://localhost:3000/api/v1/students`

| Method | Endpoint | Query / Path Parameters | Request Body (JSON) | Possible Statuses | Response Body Example |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `GET` | `/api/v1/students` | `page` (int, default 1)<br>`limit` (int, default 10, max 100)<br>`search` (string, case-insensitive)<br>`sort` (id, nim, name, grade)<br>`order` (asc, desc)<br>`is_active` (bool) | *(None)* | `200 OK` | `{"success": true, "message": "daftar mahasiswa berhasil diambil", "data": [...], "meta": {"page": 1, "limit": 10, "total": 1, "total_pages": 1}}` |
| `GET` | `/api/v1/students/:id` | `:id` (int, path) | *(None)* | `200 OK`<br>`400 Bad Request`<br>`404 Not Found` | `{"success": true, "message": "mahasiswa ditemukan", "data": {"id": 1, "nim": "000000001", "name": "mahasiswaAA", "grade": 3.75, "is_active": true, "created_at": "..."}}` |
| `POST` | `/api/v1/students` | *(None)* | `{"nim": "000000001", "name": "mahasiswaAA", "grade": 3.75}` | `201 Created`<br>`400 Bad Request`<br>`409 Conflict`<br>`415 Unsupported Media Type`<br>`422 Unprocessable Entity` | `{"success": true, "message": "mahasiswa berhasil dibuat", "data": {...}}`<br>*Header: Location: /api/v1/students/1* |
| `PUT` | `/api/v1/students/:id` | `:id` (int, path) | `{"nim": "000000001", "name": "mahasiswaAA_baru", "grade": 3.80, "is_active": false}` *(All fields mandatory)* | `200 OK`<br>`400 Bad Request`<br>`404 Not Found`<br>`409 Conflict`<br>`415 Unsupported Media Type`<br>`422 Unprocessable Entity` | `{"success": true, "message": "mahasiswa berhasil diganti seluruhnya", "data": {...}}` |
| `PATCH` | `/api/v1/students/:id` | `:id` (int, path) | `{"grade": 3.90}` *(Only modified fields required)* | `200 OK`<br>`400 Bad Request`<br>`404 Not Found`<br>`409 Conflict`<br>`415 Unsupported Media Type`<br>`422 Unprocessable Entity` | `{"success": true, "message": "mahasiswa berhasil diperbarui sebagian", "data": {...}}` |
| `DELETE`| `/api/v1/students/:id` | `:id` (int, path) | *(None)* | `204 No Content`<br>`400 Bad Request`<br>`404 Not Found` | *(Empty Body)* |
---

## References Used

* [https://www.youtube.com/watch?v=8uiZC0l4Ajw](https://www.youtube.com/watch?v=8uiZC0l4Ajw)
* [https://www.youtube.com/watch?v=d_L64KT3SFM](https://www.youtube.com/watch?v=d_L64KT3SFM)
* [https://docs.gofiber.io/](https://docs.gofiber.io/)
* [https://developer.mozilla.org/en-US/docs/Web/HTTP](https://developer.mozilla.org/en-US/docs/Web/HTTP)