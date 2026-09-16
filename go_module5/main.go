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
	"go_module5/helper"
	"go_module5/route"
)

const minSecretLength = 32

func main() {
	// 1. Config
	config.LoadEnv()
	logger := config.NewLogger()

	jwtSecret := config.GetEnv("JWT_SECRET", "")
	if len(jwtSecret) < minSecretLength {
		logger.Error("JWT_SECRET is not set or too short",
			slog.Int("min_chars", minSecretLength))
		os.Exit(1)
	}

	// 2. Database
	pool, err := database.NewPool(context.Background())
	if err != nil {
		logger.Error("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.Close()

	// 3. Helpers & Managers
	jwtManager := helper.NewJWTManager(
		jwtSecret,
		config.GetEnv("JWT_ISSUER", "praktikum-backend"),
		time.Duration(config.GetEnvInt("JWT_ACCESS_TTL_MINUTES", 15))*time.Minute,
	)

	// 4. Repositories & Services
	studentRepo := repository.NewStudentRepository(pool)
	tokenRepo := repository.NewTokenRepository(pool)

	studentService := service.NewStudentService(studentRepo)
	authService := service.NewAuthService(
		studentRepo,
		tokenRepo,
		jwtManager,
		time.Duration(config.GetEnvInt("JWT_REFRESH_TTL_DAYS", 7))*24*time.Hour,
	)

	// 5. App & Routes
	app := config.NewApp(logger, route.Dependencies{
		Pool:           pool,
		JWT:            jwtManager,
		StudentService: studentService,
		AuthService:    authService,
	})

	port := config.GetEnv("APP_PORT", "3000")

	go func() {
		if err := app.Listen(":" + port); err != nil {
			logger.Error("server stopped", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	logger.Info("server running", slog.String("port", port))

	// 6. Graceful Shutdown
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
