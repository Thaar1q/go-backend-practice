package config

import (
	"go_module5/app/service"
	"go_module5/helper"
	"go_module5/middleware"
	"go_module5/route"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewApp(
	logger *slog.Logger, pool *pgxpool.Pool, studentService *service.StudentService,
) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      GetEnv("APP_NAME", "Go Backend Practice"),
		ErrorHandler: newErrorHandler(logger),
	})
	middleware.Register(app, logger)
	route.Register(app, pool, studentService)

	app.Use(func(c *fiber.Ctx) error {
		return helper.Fail(c, fiber.StatusNotFound, "endpoint not found")
	})
	return app
}

func newErrorHandler(logger *slog.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		status := fiber.StatusInternalServerError
		message := "500 Internal Server Error"
		if e, ok := err.(*fiber.Error); ok {
			status = e.Code
			message = e.Message
		}
		logger.Error("unhandled_error",
			slog.String("path", c.Path()),
			slog.Int("status", status),
			slog.String("error", err.Error()),
		)
		return helper.Fail(c, status, message)
	}
}
