package config

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"go_module5/helper"
	"go_module5/middleware"
	"go_module5/route"
)

func NewApp(logger *slog.Logger, deps route.Dependencies) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      GetEnv("APP_NAME", "Go Backend Practice"),
		ErrorHandler: newErrorHandler(logger),
		BodyLimit:    1 * 1024 * 1024, // 1 MB
	})

	allowedOrigins := GetEnv("CORS_ALLOWED_ORIGINS", "*")
	middleware.Register(app, logger, allowedOrigins)

	route.Register(app, deps)

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
