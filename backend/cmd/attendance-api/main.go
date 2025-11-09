package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/trustmedis/mini-attendance/internal/config"
	"github.com/trustmedis/mini-attendance/internal/infra/db"
	"github.com/trustmedis/mini-attendance/internal/infra/kafka"
	"github.com/trustmedis/mini-attendance/internal/infra/redis"
	"github.com/trustmedis/mini-attendance/internal/infra/observability"
	"github.com/trustmedis/mini-attendance/internal/api/http/handler"
	"github.com/trustmedis/mini-attendance/internal/api/http/middleware"
)

// @title Mini Attendance API
// @version 1.0
// @description Attendance management system with check-in/check-out functionality
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

func main() {
	// Load environment variables
	_ = godotenv.Load()

	// Initialize logger
	logger := observability.InitLogger()
	defer logger.Sync()

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	dbConn, err := db.NewPostgresConnection(cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", "error", err)
	}
	defer dbConn.Close(context.Background())

	// Run migrations
	if err := db.RunMigrations(dbConn, cfg.Database.MigrationsPath); err != nil {
		logger.Fatal("Failed to run migrations", "error", err)
	}

	// Initialize Redis
	redisClient := redis.NewRedisClient(cfg.Redis)
	defer redisClient.Close()

	// Initialize Kafka producer
	kafkaProducer, err := kafka.NewProducer(cfg.Kafka)
	if err != nil {
		logger.Fatal("Failed to create Kafka producer", "error", err)
	}
	defer kafkaProducer.Close()

	// Setup Gin router
	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Apply middlewares
	router.Use(middleware.LoggingMiddleware(logger))
	router.Use(middleware.ErrorHandlingMiddleware())
	router.Use(middleware.CORSMiddleware())

	// Initialize handlers
	h := handler.NewHandlers(dbConn, redisClient, kafkaProducer, logger)

	// Register routes
	registerRoutes(router, h, logger)

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.App.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Starting server", "port", cfg.App.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server error", "error", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
	}

	logger.Info("Server exited")
}

func registerRoutes(router *gin.Engine, h *handler.Handlers, logger observability.Logger) {
	// Health check
	router.GET("/health", h.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")

	// Auth routes
	auth := v1.Group("/auth")
	{
		auth.POST("/login", h.Login)
		auth.POST("/register", h.Register)
		auth.POST("/refresh", h.RefreshToken)
	}

	// Attendance routes (protected)
	attendance := v1.Group("/attendance")
	attendance.Use(middleware.AuthMiddleware())
	{
		attendance.POST("/check-in", h.CheckIn)
		attendance.POST("/check-out", h.CheckOut)
		attendance.GET("/today", h.GetTodayAttendance)
		attendance.GET("/history", h.GetAttendanceHistory)
	}

	// Report routes (protected)
	reports := v1.Group("/reports")
	reports.Use(middleware.AuthMiddleware())
	{
		reports.GET("/absence", h.GetAbsenceReport)
		reports.GET("/absence/export", h.ExportAbsenceReport)
	}

	// Swagger documentation
	router.GET("/swagger/*any", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Swagger UI available at /swagger/index.html"})
	})
}
