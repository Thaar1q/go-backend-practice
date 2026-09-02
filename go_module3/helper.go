package main

import (
	"context"
	"strconv"
	"strings"
	"time"

	"go_module3/app/model"

	"github.com/gofiber/fiber/v2"
)

// 1. Response Wrappers
func ok(c *fiber.Ctx, message string, data any) error {
	return c.Status(fiber.StatusOK).JSON(model.WebResponse{
		Success: true, Message: message, Data: data,
	})
}

func okList(c *fiber.Ctx, message string, data any, meta *model.Meta) error {
	return c.Status(fiber.StatusOK).JSON(model.WebResponse{
		Success: true, Message: message, Data: data, Meta: meta,
	})
}

func created(c *fiber.Ctx, message string, data any, location string) error {
	c.Set("Location", location)
	return c.Status(fiber.StatusCreated).JSON(model.WebResponse{
		Success: true, Message: message, Data: data,
	})
}

func noContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

func fail(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(model.WebResponse{Success: false, Message: message})
}

func failValidation(c *fiber.Ctx, errs map[string]string) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(model.WebResponse{
		Success: false, Message: "validation failed", Errors: errs,
	})
}

// 2. Query Parser
var allowedSort = map[string]bool{
	"id": true, "name": true, "nim": true, "grade": true, "is_active": true, "created_at": true,
}

func parseListQuery(c *fiber.Ctx) model.ListQuery {
	q := model.ListQuery{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 10),
		Search: strings.TrimSpace(c.Query("search")),
		Sort:   c.Query("sort", "id"),
		Order:  strings.ToLower(c.Query("order", "asc")),
	}

	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 {
		q.Limit = 10
	}
	if q.Limit > 100 { // Going with 100 for now
		q.Limit = 100
	}
	if !allowedSort[q.Sort] {
		q.Sort = "id"
	}
	if q.Order != "desc" {
		q.Order = "asc"
	}

	if raw := c.Query("is_active"); raw != "" {
		if v, err := strconv.ParseBool(raw); err == nil {
			q.IsActive = &v
		}
	}
	if raw := c.Query("min_grade"); raw != "" {
		if v, err := strconv.ParseFloat(raw, 64); err == nil {
			q.MinGrade = &v
		}
	}

	return q
}

// 3. Context Helper
func reqCtx(c *fiber.Ctx) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.UserContext(), 5*time.Second)
}
