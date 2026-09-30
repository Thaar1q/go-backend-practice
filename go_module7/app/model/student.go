package model

import "time"

// 1. Base Entity
type Student struct {
	ID        int       `json:"id"`
	NIM       string    `json:"nim"`
	Name      string    `json:"name"`
	Grade     float64   `json:"grade"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	OwnerID   int       `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
}

// 2. Request Payloads
// POST
type CreateStudentRequest struct {
	NIM      string  `json:"nim" validate:"required,min=3,max=30,number"`
	Name     string  `json:"name" validate:"required,min=3,max=30"`
	Grade    float64 `json:"grade" validate:"min=0,max=4"`
	Password string  `json:"password" validate:"required,min=8,max=72,nospace"`
	IsActive bool    `json:"is_active"`
}

// PUT
type ReplaceStudentRequest struct {
	NIM      string   `json:"nim" validate:"required,min=3,max=30,number"`
	Name     string   `json:"name" validate:"required,min=3,max=30"`
	Grade    *float64 `json:"grade" validate:"required,min=0,max=4"`
	IsActive bool     `json:"is_active"`
}

// PATCH
type PatchStudentRequest struct {
	NIM      *string  `json:"nim,omitempty" validate:"omitnil,min=3,max=30,number"`
	Name     *string  `json:"name,omitempty" validate:"omitnil,min=3,max=30"`
	Grade    *float64 `json:"grade,omitempty" validate:"omitnil,min=0,max=4"`
	IsActive *bool    `json:"is_active,omitempty"`
}

// ROLES
type AssignRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=admin student"`
}

// 3. Response Structure
// RESPONSE TEMPLATE
type WebResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

type ErrorResponse struct {
	Success   bool              `json:"success"`
	Code      string            `json:"code"`
	Message   string            `json:"message"`
	Fields    map[string]string `json:"fields,omitempty"`
	RequestID string            `json:"request_id,omitempty"`
}

// 4. Query & Metadata
type Meta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type ListQuery struct {
	Page     int
	Limit    int
	Search   string
	Sort     string
	Order    string
	IsActive *bool
	MinGrade *float64
}

func (q ListQuery) Offset() int {
	return (q.Page - 1) * q.Limit
}

// 5. Cursor
type Cursor struct {
	CreatedAt time.Time
	ID        int
}

type CursorQuery struct {
	Limit    int
	Search   string
	IsActive *bool
	After    *Cursor
}

type CursorMeta struct {
	Limit      int    `json:"limit"`
	NextCursor string `json:"next_cursor,omitempty"`
	HasMore    bool   `json:"has_more"`
}
