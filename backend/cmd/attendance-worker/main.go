package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/trustmedis/mini-attendance/internal/config"
	"github.com/trustmedis/mini-attendance/internal/infra/db"
	"github.com/trustmedis/mini-attendance/internal/infra/kafka"
	"github.com/trustmedis/mini-attendance/internal/infra/observability"
)

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

	// Initialize Kafka consumer
	consumer, err := kafka.NewConsumer(cfg.Kafka)
	if err != nil {
		logger.Fatal("Failed to create Kafka consumer", "error", err)
	}
	defer consumer.Close()

	logger.Info("Starting attendance worker service", "version", cfg.App.Version)

	// Start consuming events
	go func() {
		if err := consumer.ConsumeCheckInEvents(context.Background(), dbConn, logger); err != nil {
			logger.Error("Check-in consumer error", "error", err)
		}
	}()

	go func() {
		if err := consumer.ConsumeCheckOutEvents(context.Background(), dbConn, logger); err != nil {
			logger.Error("Check-out consumer error", "error", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down worker service...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Graceful shutdown
	if err := consumer.Close(); err != nil {
		logger.Error("Error closing consumer", "error", err)
	}

	logger.Info("Worker service exited")
}
