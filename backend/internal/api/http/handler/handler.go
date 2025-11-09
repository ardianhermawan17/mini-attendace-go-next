package handler

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trustmedis/mini-attendance/internal/infra/kafka"
	"github.com/trustmedis/mini-attendance/internal/infra/redis"
	"github.com/trustmedis/mini-attendance/internal/infra/observability"
)

type Handlers struct {
	db            *pgxpool.Pool
	redis         *redis.RedisClient
	kafkaProducer *kafka.Producer
	logger        observability.Logger
}

func NewHandlers(
	db *pgxpool.Pool,
	redis *redis.RedisClient,
	kafkaProducer *kafka.Producer,
	logger observability.Logger,
) *Handlers {
	return &Handlers{
		db:            db,
		redis:         redis,
		kafkaProducer: kafkaProducer,
		logger:        logger,
	}
}
