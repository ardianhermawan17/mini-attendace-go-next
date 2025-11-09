package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/trustmedis/mini-attendance/internal/config"
)

type Producer struct {
	writer *kafka.Writer
}

type Event struct {
	EventID       string                 `json:"event_id"`
	EventType     string                 `json:"event_type"`
	AggregateID   string                 `json:"aggregate_id"`
	AggregateType string                 `json:"aggregate_type"`
	Timestamp     int64                  `json:"timestamp"`
	Data          map[string]interface{} `json:"data"`
	Metadata      map[string]string      `json:"metadata"`
}

func NewProducer(cfg config.KafkaConfig) (*Producer, error) {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Brokers...),
		Balancer: &kafka.LeastBytes{},
	}

	return &Producer{writer: writer}, nil
}

func (p *Producer) PublishCheckInEvent(ctx context.Context, topic string, userID string, attendanceID string, data map[string]interface{}) error {
	event := Event{
		EventID:       uuid.New().String(),
		EventType:     "attendance.check_in",
		AggregateID:   attendanceID,
		AggregateType: "attendance",
		Timestamp:     int64(0), // Will be set by Kafka
		Data:          data,
		Metadata: map[string]string{
			"user_id": userID,
		},
	}

	return p.publishEvent(ctx, topic, event)
}

func (p *Producer) PublishCheckOutEvent(ctx context.Context, topic string, userID string, attendanceID string, data map[string]interface{}) error {
	event := Event{
		EventID:       uuid.New().String(),
		EventType:     "attendance.check_out",
		AggregateID:   attendanceID,
		AggregateType: "attendance",
		Timestamp:     int64(0),
		Data:          data,
		Metadata: map[string]string{
			"user_id": userID,
		},
	}

	return p.publishEvent(ctx, topic, event)
}

func (p *Producer) publishEvent(ctx context.Context, topic string, event Event) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	message := kafka.Message{
		Topic: topic,
		Key:   []byte(event.AggregateID),
		Value: eventBytes,
	}

	err = p.writer.WriteMessages(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
