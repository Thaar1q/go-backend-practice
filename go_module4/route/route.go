package route

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"go_module4/app/service"
	"go_module4/helper"
	"go_module4/middleware"
)

func Register(app *fiber.App, pool *pgxpool.Pool, studentService *service.StudentService) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	api := app.Group("/api/v1")

	api.Get("/health", healthCheck(pool))

	students := api.Group("/students", middleware.RequireJSON)
	students.Get("/", studentService.ListStudents)
	students.Get("/:id", studentService.GetStudent)
	students.Post("/", studentService.CreateStudent)
	students.Put("/:id", studentService.ReplaceStudent)
	students.Patch("/:id", studentService.PatchStudent)
	students.Delete("/:id", studentService.DeleteStudent)

}

func healthCheck(pool *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
		defer cancel()
		if err := pool.Ping(ctx); err != nil {
			return helper.Fail(c, fiber.StatusServiceUnavailable, "database down")
		}
		return helper.Success(c, "server is running", fiber.Map{"timestamp": time.Now()})
	}
}
