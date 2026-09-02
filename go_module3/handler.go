package main

import (
	"errors"
	"strconv"
	"strings"

	"go_module3/app/model"
	"go_module3/app/repository"

	"github.com/gofiber/fiber/v2"
)

// 1. Handler Struct & Constructor
type StudentHandler struct {
	repo repository.StudentRepository
}

func NewStudentHandler(repo repository.StudentRepository) *StudentHandler {
	return &StudentHandler{repo: repo}
}

// Internal Helper
func paramID(c *fiber.Ctx) (int, bool) {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 1 {
		return 0, false
	}
	return id, true
}

func translateError(c *fiber.Ctx, err error, pesanUmum string) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return fail(c, fiber.StatusNotFound, "mahasiswa tidak ditemukan")
	case errors.Is(err, repository.ErrDuplicate):
		return fail(c, fiber.StatusConflict, "NIM sudah terdaftar")
	default:
		return fail(c, fiber.StatusInternalServerError, pesanUmum)
	}
}

// 2. GET ALL (list)
func (h *StudentHandler) ListStudents(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()
	q := parseListQuery(c)

	// Fetch from DB
	result, total, err := h.repo.FindAll(ctx, q)
	if err != nil {
		return fail(c, fiber.StatusInternalServerError, "gagal mengambil daftar mahasiswa")
	}

	totalPages := (total + q.Limit - 1) / q.Limit

	return okList(c, "daftar mahasiswa berhasil diambil", result, &model.Meta{
		Page: q.Page, Limit: q.Limit, Total: total, TotalPages: totalPages,
	})
}

// 3. GET ONE
func (h *StudentHandler) GetStudent(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()
	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	s, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return translateError(c, err, "gagal mencari mahasiswa")
	}

	return ok(c, "mahasiswa ditemukan", s)
}

// 4. POST (Create)
func (h *StudentHandler) CreateStudent(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()
	var req model.CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	errs := map[string]string{}
	req.NIM = strings.TrimSpace(req.NIM)
	req.Name = strings.TrimSpace(req.Name)

	if req.NIM == "" {
		errs["nim"] = "wajib diisi"
	}
	if req.Name == "" {
		errs["name"] = "wajib diisi"
	}
	if req.Grade < 0 || req.Grade > 4.0 {
		errs["grade"] = "harus di antara 0.0 - 4.0"
	}
	if len(errs) > 0 {
		return failValidation(c, errs)
	}

	baru := model.Student{
		NIM:      req.NIM,
		Name:     req.Name,
		Grade:    req.Grade,
		IsActive: true,
	}

	// Save to DB
	createdItem, err := h.repo.Create(ctx, baru)
	if err != nil {
		return translateError(c, err, "gagal menyimpan mahasiswa")
	}

	return created(c, "mahasiswa berhasil dibuat", createdItem, "/api/v1/students/"+strconv.Itoa(createdItem.ID))
}

// 5. PUT (Replace All)
func (h *StudentHandler) ReplaceStudent(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()
	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	var req model.ReplaceStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	errs := map[string]string{}
	req.NIM = strings.TrimSpace(req.NIM)
	req.Name = strings.TrimSpace(req.Name)

	if req.NIM == "" {
		errs["nim"] = "wajib diisi pada PUT"
	}
	if req.Name == "" {
		errs["name"] = "wajib diisi pada PUT"
	}
	if req.Grade == nil {
		errs["grade"] = "wajib diisi pada PUT"
	} else if *req.Grade < 0 || *req.Grade > 4.0 {
		errs["grade"] = "harus di antara 0.0 - 4.0"
	}
	if len(errs) > 0 {
		return failValidation(c, errs)
	}

	updated := model.Student{
		ID:       id,
		NIM:      req.NIM,
		Name:     req.Name,
		Grade:    *req.Grade,
		IsActive: req.IsActive,
	}

	// Update DB
	result, err := h.repo.Update(ctx, updated)
	if err != nil {
		return translateError(c, err, "gagal memperbarui mahasiswa")
	}

	return ok(c, "mahasiswa berhasil diganti seluruhnya", result)
}

// 6. PATCH (Update Partial)
func (h *StudentHandler) PatchStudent(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()
	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	var req model.PatchStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	if req.NIM == nil && req.Name == nil && req.Grade == nil && req.IsActive == nil {
		return fail(c, fiber.StatusBadRequest, "tidak ada field yang diubah")
	}

	// 1. Fetch current data to preserve unchanged fields
	s, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return translateError(c, err, "gagal mencari mahasiswa")
	}

	// 2. Validate and map new fields
	errs := map[string]string{}
	if req.NIM != nil {
		nimVal := strings.TrimSpace(*req.NIM)
		if nimVal == "" {
			errs["nim"] = "tidak boleh kosong"
		}
		s.NIM = nimVal
	}
	if req.Name != nil {
		nameVal := strings.TrimSpace(*req.Name)
		if nameVal == "" {
			errs["name"] = "tidak boleh kosong"
		}
		s.Name = nameVal
	}
	if req.Grade != nil {
		if *req.Grade < 0 || *req.Grade > 4.0 {
			errs["grade"] = "harus di antara 0.0 - 4.0"
		}
		s.Grade = *req.Grade
	}
	if req.IsActive != nil {
		s.IsActive = *req.IsActive
	}

	if len(errs) > 0 {
		return failValidation(c, errs)
	}

	// 3. Update DB
	result, err := h.repo.Update(ctx, s)
	if err != nil {
		return translateError(c, err, "gagal memperbarui mahasiswa")
	}

	return ok(c, "mahasiswa berhasil diperbarui sebagian", result)
}

// 7. DELETE
func (h *StudentHandler) DeleteStudent(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()
	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	err := h.repo.Delete(ctx, id)
	if err != nil {
		return translateError(c, err, "gagal menghapus mahasiswa")
	}

	return noContent(c)
}
