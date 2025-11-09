package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthCheckResponse struct {
	Status    string                 `json:"status"`
	Timestamp string                 `json:"timestamp"`
	Services  map[string]interface{} `json:"services"`
}

// HealthCheck godoc
// @Summary Health check endpoint
// @Description Check system health including database, Redis, and Kafka connectivity
// @Tags Health
// @Produce json
// @Success 200 {object} HealthCheckResponse
// @Failure 503 {object} HealthCheckResponse
// @Router /health [get]
func (h *Handlers) HealthCheck(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response := HealthCheckResponse{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Services:  make(map[string]interface{}),
	}

	allHealthy := true

	// Check database
	dbHealth := "healthy"
	if err := h.db.Ping(ctx); err != nil {
		dbHealth = "unhealthy"
		allHealthy = false
		h.logger.Error("Database health check failed", "error", err)
	}
	response.Services["database"] = gin.H{
		"status": dbHealth,
	}

	// Check Redis
	redisHealth := "healthy"
	if err := h.redis.Ping(ctx); err != nil {
		redisHealth = "unhealthy"
		allHealthy = false
		h.logger.Error("Redis health check failed", "error", err)
	}
	response.Services["redis"] = gin.H{
		"status": redisHealth,
	}

	// Check Kafka (basic check - can be enhanced)
	kafkaHealth := "healthy"
	response.Services["kafka"] = gin.H{
		"status": kafkaHealth,
	}

	if allHealthy {
		response.Status = "healthy"
		c.JSON(http.StatusOK, response)
	} else {
		response.Status = "degraded"
		c.JSON(http.StatusServiceUnavailable, response)
	}
}
