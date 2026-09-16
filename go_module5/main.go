package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go_module5/app/repository"
	"go_module5/app/service"
	"go_module5/config"
	"go_module5/database"
)

func main() {
	// 1. Config
	config.LoadEnv()
	logger := config.NewLogger()

	// 2. Load Database
	pool, err := database.NewPool(context.Background())
	if err != nil {
		logger.Error("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.Close()

	// 3. pool -> repository -> handler
	studentRepository := repository.NewStudentRepository(pool)
	studentService := service.NewStudentService(studentRepository)

	// 4. Fiber App Initialization
	app := config.NewApp(logger, pool, studentService)
	port := config.GetEnv("APP_PORT", "3000")

	go func() {
		if err := app.Listen(":" + port); err != nil {
			logger.Error("server stopped", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	logger.Info("server running", slog.String("port", port))

	// 5. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutdown signal received, exiting")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Error("failed to gracefully shutdown server",
			slog.String("error", err.Error()))
	}
	logger.Info("server shutdown complete")
}
