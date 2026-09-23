package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"go_module6/app/model"
)

var (
	ErrNotFound  = errors.New("data not found")
	ErrDuplicate = errors.New("data already exists")
)

// 1. Interface & Structs
type StudentRepository interface {
	FindAll(ctx context.Context, q model.ListQuery) ([]model.Student, int, error)
	FindByID(ctx context.Context, id int) (model.Student, error)
	FindByUsername(ctx context.Context, username string) (model.Student, error)
	Create(ctx context.Context, s model.Student) (model.Student, error)
	Update(ctx context.Context, s model.Student) (model.Student, error)
	UpdateRole(ctx context.Context, id int, role string) (model.Student, error)
	Delete(ctx context.Context, id int) error
}

// Whitelist
var kolomUrut = map[string]string{
	"id":         "id",
	"nim":        "nim",
	"name":       "name",
	"grade":      "grade",
	"role":       "role",
	"is_active":  "is_active",
	"created_at": "created_at",
}

type studentPostgresRepository struct {
	pool *pgxpool.Pool
}

// 2. Constructor
func NewStudentRepository(pool *pgxpool.Pool) StudentRepository {
	return &studentPostgresRepository{pool: pool}
}

// 3. Query Builder Helper
func buildFilter(q model.ListQuery) (string, []any) {
	where := " WHERE 1 = 1"
	args := []any{}

	if q.Search != "" {
		where += fmt.Sprintf(" AND name ILIKE $%d", len(args)+1)
		args = append(args, "%"+q.Search+"%")
	}

	if q.IsActive != nil {
		where += fmt.Sprintf(" AND is_active = $%d", len(args)+1)
		args = append(args, *q.IsActive)
	}

	if q.MinGrade != nil {
		where += fmt.Sprintf(" AND grade >= $%d", len(args)+1)
		args = append(args, *q.MinGrade)
	}

	return where, args
}

// 4. Find All (List)
func (r *studentPostgresRepository) FindAll(ctx context.Context, q model.ListQuery) ([]model.Student, int, error) {
	where, args := buildFilter(q)

	var total int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM students"+where, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count students: %w", err)
	}

	arah := "ASC"
	if q.Order == "desc" {
		arah = "DESC"
	}

	sqlText := fmt.Sprintf(
		`SELECT id, nim, name, grade, password, role, is_active, owner_id, created_at
         FROM students%s
         ORDER BY %s %s
         LIMIT $%d OFFSET $%d`,
		where, kolomUrut[q.Sort], arah, len(args)+1, len(args)+2,
	)

	offset := (q.Page - 1) * q.Limit
	args = append(args, q.Limit, offset)

	rows, err := r.pool.Query(ctx, sqlText, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("fetch student list: %w", err)
	}
	defer rows.Close()

	hasil := []model.Student{}
	for rows.Next() {
		var s model.Student
		if err := rows.Scan(&s.ID, &s.NIM, &s.Name, &s.Grade, &s.Password, &s.Role, &s.IsActive, &s.OwnerID, &s.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("read student row: %w", err)
		}
		hasil = append(hasil, s)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("read query result: %w", err)
	}

	return hasil, total, nil
}

// 5. Find By ID (Get One)
func (r *studentPostgresRepository) FindByID(ctx context.Context, id int) (model.Student, error) {
	var s model.Student

	err := r.pool.QueryRow(ctx,
		`SELECT id, nim, name, grade, password, role, is_active, owner_id, created_at
         FROM students WHERE id = $1`, id,
	).Scan(&s.ID, &s.NIM, &s.Name, &s.Grade, &s.Password, &s.Role, &s.IsActive, &s.OwnerID, &s.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Student{}, ErrNotFound
		}
		return model.Student{}, fmt.Errorf("fetch student: %w", err)
	}

	return s, nil
}

// 6. Find by Username
func (r *studentPostgresRepository) FindByUsername(ctx context.Context, username string) (model.Student, error) {
	var s model.Student

	err := r.pool.QueryRow(ctx,
		`SELECT id, nim, name, grade, password, role, is_active, owner_id, created_at
		 FROM students
		 WHERE LOWER(name) = LOWER($1) OR LOWER(nim) = LOWER($1)`, username,
	).Scan(&s.ID, &s.NIM, &s.Name, &s.Grade, &s.Password, &s.Role, &s.IsActive, &s.OwnerID, &s.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Student{}, ErrNotFound
		}
		return model.Student{}, fmt.Errorf("mengambil student: %w", err)
	}

	return s, nil
}

// 7. Create (Insert)
func (r *studentPostgresRepository) Create(ctx context.Context, s model.Student) (model.Student, error) {
	if s.Role == "" {
		s.Role = "student"
	}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO students (nim, name, grade, password, role, is_active, owner_id)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, created_at`,
		s.NIM, s.Name, s.Grade, s.Password, s.Role, s.IsActive, s.OwnerID,
	).Scan(&s.ID, &s.CreatedAt)

	if err != nil {
		if isUniqueViolation(err) {
			return model.Student{}, ErrDuplicate
		}
		return model.Student{}, fmt.Errorf("save student: %w", err)
	}

	return s, nil
}

// 8. Update (Replace/Patch)
func (r *studentPostgresRepository) Update(ctx context.Context, s model.Student) (model.Student, error) {
	err := r.pool.QueryRow(ctx,
		`UPDATE students SET nim = $1, name = $2, grade = $3, is_active = $4
         WHERE id = $5
		 RETURNING id, nim, name, grade, password, role, is_active, owner_id, created_at`,
		s.NIM, s.Name, s.Grade, s.IsActive, s.ID,
	).Scan(&s.ID, &s.NIM, &s.Name, &s.Grade, &s.Password, &s.Role, &s.IsActive, &s.OwnerID, &s.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Student{}, ErrNotFound
		}
		if isUniqueViolation(err) {
			return model.Student{}, ErrDuplicate
		}
		return model.Student{}, fmt.Errorf("update student: %w", err)
	}

	return s, nil
}

// 9. Update Role
func (r *studentPostgresRepository) UpdateRole(ctx context.Context, id int, role string) (model.Student, error) {
	var s model.Student
	err := r.pool.QueryRow(ctx,
		`UPDATE students SET role = $1
		 WHERE id = $2
		 RETURNING id, nim, name, grade, password, role, is_active, owner_id, created_at`,
		role, id,
	).Scan(&s.ID, &s.NIM, &s.Name, &s.Grade, &s.Password, &s.Role, &s.IsActive, &s.OwnerID, &s.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Student{}, ErrNotFound
		}
		return model.Student{}, fmt.Errorf("update student role: %w", err)
	}

	return s, nil
}

// 10. Delete
func (r *studentPostgresRepository) Delete(ctx context.Context, id int) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM students WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete student: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// 11. Error Check Helper
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
