package service

import (
	"errors"
	"strconv"
	"strings"

	"go_module6/app/model"
	"go_module6/app/repository"
	"go_module6/helper"

	"github.com/gofiber/fiber/v2"
)

// 1. Handler Struct & Constructor
type StudentService struct {
	repo        repository.StudentRepository
	permissions *helper.PermissionSet
}

func NewStudentService(
	repo repository.StudentRepository,
	permissions *helper.PermissionSet,
) *StudentService {
	return &StudentService{
		repo:        repo,
		permissions: permissions,
	}
}

func translateError(c *fiber.Ctx, err error, pesanUmum string) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return helper.Fail(c, fiber.StatusNotFound, "student not found")
	case errors.Is(err, repository.ErrDuplicate):
		return helper.Fail(c, fiber.StatusConflict, "NIM already registered")
	default:
		return helper.Fail(c, fiber.StatusInternalServerError, pesanUmum)
	}
}

// 2. GET ALL (list)
func (h *StudentService) ListStudents(c *fiber.Ctx) error {
	ctx, cancel := helper.ReqCtx(c)
	defer cancel()
	q := helper.ParseListQuery(c)

	// Fetch from DB
	result, total, err := h.repo.FindAll(ctx, q)
	if err != nil {
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to fetch student list")
	}

	totalPages := CountTotalPages(total, q.Limit)

	return helper.SuccessList(c, "student list successfully retrieved", result, &model.Meta{
		Page: q.Page, Limit: q.Limit, Total: total, TotalPages: totalPages,
	})
}

// 3. GET ONE
func (h *StudentService) GetStudent(c *fiber.Ctx) error {
	ctx, cancel := helper.ReqCtx(c)
	defer cancel()

	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Fail(c, fiber.StatusUnauthorized, "unauthenticated")
	}

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.Fail(c, fiber.StatusBadRequest, "id must be a positive number")
	}

	s, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return translateError(c, err, "failed to find student")
	}

	if !CanAccessStudent(current, s.OwnerID, h.permissions, "student:read:any") {
		return helper.Fail(c, fiber.StatusForbidden, "not allowed to access others' data")
	}

	return helper.Success(c, "student found", s)
}

// 4. POST (Create)
func (h *StudentService) CreateStudent(c *fiber.Ctx) error {
	ctx, cancel := helper.ReqCtx(c)
	defer cancel()

	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Fail(c, fiber.StatusUnauthorized, "unauthenticated")
	}

	var req model.CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "body must be valid JSON")
	}

	req.NIM = strings.TrimSpace(req.NIM)
	req.Name = strings.TrimSpace(req.Name)

	if errs := ValidateCreate(req); len(errs) > 0 {
		return helper.FailValidation(c, errs)
	}

	new := model.Student{
		NIM:      req.NIM,
		Name:     req.Name,
		Grade:    req.Grade,
		Role:     "student",
		IsActive: true,
		OwnerID:  current.UserID,
	}

	// Save to DB
	createdItem, err := h.repo.Create(ctx, new)
	if err != nil {
		return translateError(c, err, "failed to save student")
	}

	return helper.Created(c, "student successfully created", createdItem, "/api/v1/students/"+strconv.Itoa(createdItem.ID))
}

// 5. PUT (Replace All)
func (h *StudentService) ReplaceStudent(c *fiber.Ctx) error {
	ctx, cancel := helper.ReqCtx(c)
	defer cancel()

	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Fail(c, fiber.StatusUnauthorized, "unauthenticated")
	}

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.Fail(c, fiber.StatusBadRequest, "id must be a positive number")
	}

	s, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return translateError(c, err, "failed to find student")
	}

	if !CanAccessStudent(current, s.OwnerID, h.permissions, "student:update:any") {
		return helper.Fail(c, fiber.StatusForbidden, "not allowed to update others' data")
	}

	var req model.ReplaceStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "body must be valid JSON")
	}

	req.NIM = strings.TrimSpace(req.NIM)
	req.Name = strings.TrimSpace(req.Name)

	if errs := ValidateReplace(req); len(errs) > 0 {
		return helper.FailValidation(c, errs)
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
		return translateError(c, err, "failed to update student")
	}

	return helper.Success(c, "student successfully replaced", result)
}

// 6. PATCH (Update Partial)
func (h *StudentService) PatchStudent(c *fiber.Ctx) error {
	ctx, cancel := helper.ReqCtx(c)
	defer cancel()

	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Fail(c, fiber.StatusUnauthorized, "unauthenticated")
	}

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.Fail(c, fiber.StatusBadRequest, "id must be a positive number")
	}

	s, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return translateError(c, err, "failed to find student")
	}

	if !CanAccessStudent(current, s.OwnerID, h.permissions, "student:update:any") {
		return helper.Fail(c, fiber.StatusForbidden, "not allowed to update others' data")
	}

	var req model.PatchStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "body must be valid JSON")
	}

	if IsEmptyPatch(req) {
		return helper.Fail(c, fiber.StatusBadRequest, "no fields to update")
	}

	// 1. Validate and map new fields
	updated, errs := ApplyPatch(s, req)
	if len(errs) > 0 {
		return helper.FailValidation(c, errs)
	}

	// 2. Update DB
	result, err := h.repo.Update(ctx, updated)
	if err != nil {
		return translateError(c, err, "failed to update student")
	}

	return helper.Success(c, "student successfully partially updated", result)
}

// 7. PATCH ROLE
func (h *StudentService) AssignRole(c *fiber.Ctx) error {
	ctx, cancel := helper.ReqCtx(c)
	defer cancel()

	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Fail(c, fiber.StatusUnauthorized, "unauthenticated")
	}

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.Fail(c, fiber.StatusBadRequest, "id must be a positive number")
	}

	var req model.AssignRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "body must be valid JSON")
	}

	if errs := ValidateAssignRole(current, id, req, h.permissions); len(errs) > 0 {
		return helper.FailValidation(c, errs)
	}

	result, err := h.repo.UpdateRole(ctx, id, strings.TrimSpace(req.Role))
	if err != nil {
		return translateError(c, err, "failed to update student role")
	}

	return helper.Success(c, "student role successfully updated", result)
}

// 8. DELETE
func (h *StudentService) DeleteStudent(c *fiber.Ctx) error {
	ctx, cancel := helper.ReqCtx(c)
	defer cancel()

	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Fail(c, fiber.StatusUnauthorized, "unauthenticated")
	}

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.Fail(c, fiber.StatusBadRequest, "id must be a positive number")
	}

	if current.UserID == id {
		return helper.Fail(c, fiber.StatusForbidden, "cannot delete own account")
	}

	err := h.repo.Delete(ctx, id)
	if err != nil {
		return translateError(c, err, "failed to delete student")
	}

	return helper.NoContent(c)
}
