package main

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

var students []Student
var nextID = 1

// Storage and Help

func findStudentIndex(id int) int {
	for i := range students {
		if students[i].ID == id {
			return i
		}
	}
	return -1
}

func searchMatch(s Student, search string) bool {
	search = strings.ToLower(search)
	return strings.Contains(strings.ToLower(s.Name), search)
}

func paramID(c *fiber.Ctx) (int, bool) {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 1 {
		return 0, false
	}
	return id, true
}

// List (w. Filter, Sort, Paginate)
func listStudents(c *fiber.Ctx) error {
	q := parseListQuery(c)

	// 1 - Filter
	result := []Student{}
	for _, s := range students {
		if q.IsActive != nil && s.IsActive != *q.IsActive {
			continue
		}
		if q.MinGrade != nil && s.Grade < *q.MinGrade {
			continue
		}
		if q.Search != "" && !searchMatch(s, q.Search) {
			continue
		}

		result = append(result, s)
	}

	// 2 - Sort
	sort.SliceStable(result, func(i, j int) bool {
		var smallerThan bool
		switch q.Sort {
		case "name":
			smallerThan = result[i].Name < result[j].Name
		case "nim":
			smallerThan = result[i].NIM < result[j].NIM
		case "grade":
			smallerThan = result[i].Grade < result[j].Grade
		case "created_at":
			smallerThan = result[i].CreatedAt.Before(result[j].CreatedAt)
		default:
			smallerThan = result[i].ID < result[j].ID
		}
		if q.Order == "desc" {
			return !smallerThan
		}
		return smallerThan
	})

	// 3 - Paginate
	total := len(result)
	totalPages := (total + q.Limit - 1) / q.Limit
	pageStart := (q.Page - 1) * q.Limit
	if pageStart > total {
		pageStart = total
	}
	pageEnd := pageStart + q.Limit
	if pageEnd > total {
		pageEnd = total
	}

	return okList(c, "daftar mahasiswa berhasil diambil", result[pageStart:pageEnd], &Meta{
		Page: q.Page, Limit: q.Limit, Total: total, TotalPages: totalPages,
	})
}

func getStudent(c *fiber.Ctx) error {
	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	i := findStudentIndex(id)
	if i == -1 {
		return fail(c, fiber.StatusNotFound, "mahasiswa tidak ditemukan")
	}

	return ok(c, "mahasiswa ditemukan", students[i])
}

func createStudent(c *fiber.Ctx) error {
	var req CreateStudentRequest
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
	for _, s := range students {
		if strings.EqualFold(s.NIM, req.NIM) {
			return fail(c, fiber.StatusConflict, "NIM sudah digunakan")
		}
	}
	if len(errs) > 0 {
		return failValidation(c, errs)
	}

	newItem := Student{
		ID:        nextID,
		NIM:       req.NIM,
		Name:      req.Name,
		Grade:     req.Grade,
		IsActive:  true,
		CreatedAt: time.Now(),
	}
	students = append(students, newItem)
	nextID++

	return created(c, "mahasiswa berhasil dibuat", newItem,
		"/api/v1/students/"+strconv.Itoa(newItem.ID))
}

func replaceStudent(c *fiber.Ctx) error {
	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	i := findStudentIndex(id)
	if i == -1 {
		return fail(c, fiber.StatusNotFound, "mahasiswa tidak ditemukan")
	}

	var req ReplaceStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	errs := map[string]string{}
	if strings.TrimSpace(req.NIM) == "" {
		errs["nim"] = "wajib diisi pada PUT"
	}
	if strings.TrimSpace(req.Name) == "" {
		errs["name"] = "wajib diisi pada PUT"
	}
	if req.Grade == nil {
		errs["grade"] = "wajib diisi pada PUT"
	} else if *req.Grade < 0 || *req.Grade > 4.0 {
		errs["grade"] = "harus di antara 0.0 - 4.0"
	}
	for idx, s := range students {
		if idx != i && strings.EqualFold(s.NIM, req.NIM) {
			return fail(c, fiber.StatusConflict, "NIM sudah digunakan")
		}
	}
	if len(errs) > 0 {
		return failValidation(c, errs)
	}

	students[i].NIM = req.NIM
	students[i].Name = req.Name
	students[i].Grade = *req.Grade
	students[i].IsActive = req.IsActive

	return ok(c, "mahasiswa berhasil diganti seluruhnya", students[i])
}

func patchStudent(c *fiber.Ctx) error {
	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	i := findStudentIndex(id)
	if i == -1 {
		return fail(c, fiber.StatusNotFound, "mahasiswa tidak ditemukan")
	}

	var req PatchStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	if req.NIM == nil && req.Name == nil && req.Grade == nil && req.IsActive == nil {
		return fail(c, fiber.StatusBadRequest, "tidak ada field yang diubah")
	}

	if req.NIM != nil {
		if strings.TrimSpace(*req.NIM) == "" {
			return failValidation(c, map[string]string{"nim": "tidak boleh kosong"})
		}
		for idx, s := range students {
			if idx != i && strings.EqualFold(s.NIM, *req.NIM) {
				return fail(c, fiber.StatusConflict, "NIM sudah digunakan")
			}
		}
		students[i].NIM = *req.NIM
	}
	if req.Name != nil {
		if strings.TrimSpace(*req.Name) == "" {
			return failValidation(c, map[string]string{"name": "tidak boleh kosong"})
		}
		students[i].Name = *req.Name
	}
	if req.Grade != nil {
		if *req.Grade < 0 || *req.Grade > 4.0 {
			return failValidation(c, map[string]string{"grade": "harus di antara 0.0 - 4.0"})
		}
		students[i].Grade = *req.Grade
	}
	if req.IsActive != nil {
		students[i].IsActive = *req.IsActive
	}

	return ok(c, "mahasiswa berhasil diperbarui sebagian", students[i])
}

func deleteStudent(c *fiber.Ctx) error {
	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	i := findStudentIndex(id)
	if i == -1 {
		return fail(c, fiber.StatusNotFound, "mahasiswa tidak ditemukan")
	}

	students = append(students[:i], students[i+1:]...)

	return noContent(c)
}
