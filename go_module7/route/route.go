package route

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"go_module7/app/service"
	"go_module7/helper"
	"go_module7/middleware"
)

type Dependencies struct {
	Pool           *pgxpool.Pool
	JWT            *helper.JWTManager
	Permissions    *helper.PermissionSet
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

	// Auth routes
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

	perms := deps.Permissions

	// Route-level permission guards
	students.Get("/", middleware.RequirePermission(perms, "student:list"), deps.StudentService.ListStudents)
	students.Post("/", middleware.RequirePermission(perms, "student:create"), deps.StudentService.CreateStudent)
	students.Delete("/:id", middleware.RequirePermission(perms, "student:delete"), deps.StudentService.DeleteStudent)
	students.Patch("/:id/role", middleware.RequirePermission(perms, "role:assign"), deps.StudentService.AssignRole)

	// Ownership-checked
	students.Get("/:id", deps.StudentService.GetStudent)
	students.Put("/:id", deps.StudentService.ReplaceStudent)
	students.Patch("/:id", deps.StudentService.PatchStudent)
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
