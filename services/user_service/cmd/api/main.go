package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pmii/pmii-backend/shared/jwt"
	"github.com/pmii/pmii-backend/shared/logger"
	"github.com/pmii/user-service/config"
	"github.com/pmii/user-service/internal/delivery/http/handler"
	"github.com/pmii/user-service/internal/delivery/http/router"
	"github.com/pmii/user-service/internal/infrastructure/database"
	"github.com/pmii/user-service/internal/repository/mysql"
	"github.com/pmii/user-service/internal/usecase/user"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// Initialize logger
	if err := logger.InitLogger(cfg.App.Env); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.Log.Sync()

	logger.Info("Starting User Service",
		zap.String("app", cfg.App.Name),
		zap.String("env", cfg.App.Env),
		zap.String("port", cfg.App.Port),
	)

	// Initialize database connection
	db, err := database.NewMySQLConnection(&cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()
	logger.Info("Database connected successfully")

	// Run database migrations
	logger.Info("Running database migrations...")
	if err := database.RunMigrations(&cfg.Database); err != nil {
		logger.Fatal("Failed to run migrations", zap.Error(err))
	}
	logger.Info("Migrations completed successfully")

	// Initialize JWT service
	jwtService := jwt.NewJWTService(cfg.JWT.Secret, cfg.JWT.ExpireHours)

	// Initialize repositories
	userRepo := mysql.NewMysqlUserRepository(db)

	// Initialize use cases
	userUsecase := user.NewUserUsecase(userRepo, cfg.JWT.Secret, cfg.JWT.ExpireHours)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userUsecase)

	// Setup router
	r := router.SetupRouter(userHandler, jwtService)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Info(fmt.Sprintf("Server running on port %s", cfg.App.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited properly")
}
