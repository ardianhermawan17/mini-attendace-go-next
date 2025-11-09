package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"
	"github.com/trustmedis/mini-attendance/internal/config"
	"github.com/trustmedis/mini-attendance/internal/infra/observability"
)

type Consumer struct {
	reader *kafka.Reader
}

type CheckInEventData struct {
	UserID         string            `json:"user_id"`
	AttendanceDate string            `json:"attendance_date"`
	CheckInTime    string            `json:"check_in_time"`
	CheckInSource  string            `json:"check_in_source"`
	CheckInMeta    map[string]string `json:"check_in_meta"`
}

type CheckOutEventData struct {
	UserID         string            `json:"user_id"`
	AttendanceDate string            `json:"attendance_date"`
	CheckInTime    string            `json:"check_in_time"`
	CheckOutTime   string            `json:"check_out_time"`
	CheckOutSource string            `json:"check_out_source"`
	CheckOutMeta   map[string]string `json:"check_out_meta"`
}

func NewConsumer(cfg config.KafkaConfig) (*Consumer, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.Brokers,
		Topic:          cfg.Topics.CheckIn,
		GroupID:        "attendance-report-service",
		StartOffset:    kafka.LastOffset,
		CommitInterval: time.Second,
	})

	return &Consumer{reader: reader}, nil
}

func (c *Consumer) ConsumeCheckInEvents(ctx context.Context, db *pgxpool.Pool, logger observability.Logger) error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"localhost:9092"}, // Should come from config
		Topic:          "attendance.check-in",
		GroupID:        "attendance-report-service",
		StartOffset:    kafka.LastOffset,
		CommitInterval: time.Second,
	})
	defer reader.Close()

	logger.Info("Starting check-in event consumer")

	for {
		message, err := reader.ReadMessage(ctx)
		if err != nil {
			logger.Error("Error reading message", "error", err)
			continue
		}

		// Parse event
		var event Event
		if err := json.Unmarshal(message.Value, &event); err != nil {
			logger.Error("Failed to unmarshal event", "error", err)
			continue
		}

		// Check if event already processed (idempotency)
		var processed bool
		err = db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM processed_events WHERE event_id = $1)", event.EventID).Scan(&processed)
		if err != nil {
			logger.Error("Failed to check processed events", "error", err)
			continue
		}

		if processed {
			logger.Debug("Event already processed", "event_id", event.EventID)
			continue
		}

		// Process check-in event
		if err := c.processCheckInEvent(ctx, db, event, logger); err != nil {
			logger.Error("Failed to process check-in event", "event_id", event.EventID, "error", err)
			continue
		}

		// Mark event as processed
		_, err = db.Exec(ctx,
			"INSERT INTO processed_events (id, event_id, event_type) VALUES ($1, $2, $3)",
			uuid.New(), event.EventID, event.EventType,
		)
		if err != nil {
			logger.Error("Failed to mark event as processed", "error", err)
		}

		logger.Info("Check-in event processed", "event_id", event.EventID)
	}
}

func (c *Consumer) ConsumeCheckOutEvents(ctx context.Context, db *pgxpool.Pool, logger observability.Logger) error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"localhost:9092"}, // Should come from config
		Topic:          "attendance.check-out",
		GroupID:        "attendance-report-service",
		StartOffset:    kafka.LastOffset,
		CommitInterval: time.Second,
	})
	defer reader.Close()

	logger.Info("Starting check-out event consumer")

	for {
		message, err := reader.ReadMessage(ctx)
		if err != nil {
			logger.Error("Error reading message", "error", err)
			continue
		}

		// Parse event
		var event Event
		if err := json.Unmarshal(message.Value, &event); err != nil {
			logger.Error("Failed to unmarshal event", "error", err)
			continue
		}

		// Check if event already processed (idempotency)
		var processed bool
		err = db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM processed_events WHERE event_id = $1)", event.EventID).Scan(&processed)
		if err != nil {
			logger.Error("Failed to check processed events", "error", err)
			continue
		}

		if processed {
			logger.Debug("Event already processed", "event_id", event.EventID)
			continue
		}

		// Process check-out event
		if err := c.processCheckOutEvent(ctx, db, event, logger); err != nil {
			logger.Error("Failed to process check-out event", "event_id", event.EventID, "error", err)
			continue
		}

		// Mark event as processed
		_, err = db.Exec(ctx,
			"INSERT INTO processed_events (id, event_id, event_type) VALUES ($1, $2, $3)",
			uuid.New(), event.EventID, event.EventType,
		)
		if err != nil {
			logger.Error("Failed to mark event as processed", "error", err)
		}

		logger.Info("Check-out event processed", "event_id", event.EventID)
	}
}

func (c *Consumer) processCheckInEvent(ctx context.Context, db *pgxpool.Pool, event Event, logger observability.Logger) error {
	// Extract data from event
	userID := event.Metadata["user_id"]
	attendanceDate := event.Data["attendance_date"].(string)

	// Determine status based on check-in time
	status := "hadir" // Default status

	// Upsert absence report
	_, err := db.Exec(ctx,
		`INSERT INTO absence_reports (id, user_id, report_date, status, check_in_time)
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (user_id, report_date) DO UPDATE SET
		 status = $4, check_in_time = $5, updated_at = CURRENT_TIMESTAMP`,
		uuid.New(), userID, attendanceDate, status, time.Now(),
	)

	return err
}

func (c *Consumer) processCheckOutEvent(ctx context.Context, db *pgxpool.Pool, event Event, logger observability.Logger) error {
	// Extract data from event
	userID := event.Metadata["user_id"]
	attendanceDate := event.Data["attendance_date"].(string)
	checkInTime := event.Data["check_in_time"]
	checkOutTime := event.Data["check_out_time"]

	// Calculate work hours
	var workHours float64
	if checkInTime != nil && checkOutTime != nil {
		// Parse times and calculate difference
		checkInStr := fmt.Sprintf("%v", checkInTime)
		checkOutStr := fmt.Sprintf("%v", checkOutTime)

		checkInT, _ := time.Parse("2006-01-02 15:04:05", checkInStr)
		checkOutT, _ := time.Parse("2006-01-02 15:04:05", checkOutStr)

		workHours = checkOutT.Sub(checkInT).Hours()
	}

	// Determine status based on work hours
	status := "hadir"
	if workHours < 8 {
		status = "pulang_cepat"
	}

	// Upsert absence report
	_, err := db.Exec(ctx,
		`INSERT INTO absence_reports (id, user_id, report_date, status, check_in_time, check_out_time, work_hours)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 ON CONFLICT (user_id, report_date) DO UPDATE SET
		 status = $4, check_out_time = $6, work_hours = $7, updated_at = CURRENT_TIMESTAMP`,
		uuid.New(), userID, attendanceDate, status, checkInTime, checkOutTime, workHours,
	)

	return err
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
