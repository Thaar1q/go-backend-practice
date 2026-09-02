package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"go_module3/app/repository"
	"go_module3/config"
	"go_module3/database"
)

var methodsWithBody = map[string]bool{
	fiber.MethodPost:  true,
	fiber.MethodPut:   true,
	fiber.MethodPatch: true,
}

func requireJSON(c *fiber.Ctx) error {
	if methodsWithBody[c.Method()] {
		ct := c.Get("Content-Type")
		if !strings.HasPrefix(ct, fiber.MIMEApplicationJSON) {
			return fail(c, fiber.StatusUnsupportedMediaType,
				"Content-Type must be application/json")
		}
	}
	return c.Next()
}

func main() {
	// 1. Config
	config.LoadEnv()

	// 2. Load Database
	pool, err := database.NewPool(context.Background())
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()

	// 3. pool -> repository -> handler
	StudentRepository := repository.NewStudentRepository(pool)
	StudentHandler := NewStudentHandler(StudentRepository)

	app := fiber.New(fiber.Config{
		AppName: "Praktikum Backend Lanjut - Pertemuan 3",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			status := fiber.StatusInternalServerError
			pesan := "500 Internal Server Error"
			if e, ok := err.(*fiber.Error); ok {
				status = e.Code
				pesan = e.Message
			}
			return fail(c, status, pesan)
		},
	})

	// Global middleware
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${locals:requestid} ${method} ${path} ${status} ${latency}\n",
	}))
	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	api := app.Group("/api/v1")

	api.Get("/health", func(c *fiber.Ctx) error {
		return ok(c, "server berjalan", fiber.Map{"timestamp": time.Now()})
	})

	u := api.Group("/students", requireJSON)
	u.Get("/", StudentHandler.ListStudents)
	u.Get("/:id", StudentHandler.GetStudent)
	u.Post("/", StudentHandler.CreateStudent)
	u.Put("/:id", StudentHandler.ReplaceStudent)
	u.Patch("/:id", StudentHandler.PatchStudent)
	u.Delete("/:id", StudentHandler.DeleteStudent)

	app.Use(func(c *fiber.Ctx) error {
		return fail(c, fiber.StatusNotFound, "endpoint tidak ditemukan")
	})

	fmt.Println("Server berjalan di http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}
