package route

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"go_module6/app/service"
	"go_module6/helper"
	"go_module6/middleware"
)

type Dependencies struct {
	Pool           *pgxpool.Pool
	JWT            *helper.JWTManager
	StudentService *service.StudentService
	AuthService    *service.AuthService
}

func Register(app *fiber.App, deps Dependencies) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	api := app.Group("/api/v1")

	// Public
	api.Get("/health", healthCheck(deps.Pool))

	// Auth routes (public / rate-limited)
	auth := api.Group("/auth", middleware.RequireJSON)
	auth.Post("/register", deps.AuthService.Register)
	auth.Post("/login", middleware.LoginRateLimiter(), deps.AuthService.Login)
	auth.Post("/refresh", deps.AuthService.Refresh)
	auth.Post("/logout", deps.AuthService.Logout)
	auth.Get("/me", middleware.RequireAuth(deps.JWT), deps.AuthService.Me)

	// Protected student routes
	students := api.Group("/students",
		middleware.RequireJSON,
		middleware.RequireAuth(deps.JWT),
	)
	students.Get("/", deps.StudentService.ListStudents)
	students.Get("/:id", deps.StudentService.GetStudent)
	students.Post("/", deps.StudentService.CreateStudent)
	students.Put("/:id", deps.StudentService.ReplaceStudent)
	students.Patch("/:id", deps.StudentService.PatchStudent)
	students.Delete("/:id", deps.StudentService.DeleteStudent)
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
