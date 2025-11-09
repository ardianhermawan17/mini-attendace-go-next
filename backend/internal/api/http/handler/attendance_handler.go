package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CheckInRequest struct {
	Source   string            `json:"source" binding:"required"`
	Metadata map[string]string `json:"metadata"`
}

type CheckOutRequest struct {
	Source   string            `json:"source" binding:"required"`
	Metadata map[string]string `json:"metadata"`
}

type AttendanceResponse struct {
	ID             string     `json:"id"`
	UserID         string     `json:"user_id"`
	AttendanceDate string     `json:"attendance_date"`
	CheckInTime    *time.Time `json:"check_in_time"`
	CheckOutTime   *time.Time `json:"check_out_time"`
	Status         string     `json:"status"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// CheckIn godoc
// @Summary Check-in for attendance
// @Description Record user check-in with distributed lock to prevent double check-in
// @Tags Attendance
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param request body CheckInRequest true "Check-in data"
// @Success 201 {object} AttendanceResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /attendance/check-in [post]
func (h *Handlers) CheckIn(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	var req CheckInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userIDStr := userID.(string)
	today := time.Now().Format("2006-01-02")

	// Acquire distributed lock
	lockKey := fmt.Sprintf("lock:attendance:checkin:%s:%s", userIDStr, today)
	lockToken := uuid.New().String()

	locked, err := h.redis.SetNX(ctx, lockKey, lockToken, 5*time.Second)
	if err != nil || !locked {
		h.logger.Warn("Failed to acquire lock for check-in", "user_id", userIDStr, "error", err)
		c.JSON(http.StatusConflict, ErrorResponse{
			Error:   "conflict",
			Message: "Please retry, another check-in is in progress",
		})
		return
	}

	defer func() {
		// Release lock
		h.redis.Del(ctx, lockKey)
	}()

	// Check if already checked in today
	var existingID *string
	err = h.db.QueryRow(ctx,
		`SELECT id FROM attendance_records 
		 WHERE user_id = $1 AND attendance_date = $2 AND check_in_at IS NOT NULL AND check_out_at IS NULL`,
		userIDStr, today,
	).Scan(&existingID)

	if err == nil && existingID != nil {
		h.logger.Warn("User already checked in", "user_id", userIDStr)
		c.JSON(http.StatusConflict, ErrorResponse{
			Error:   "already_checked_in",
			Message: "You are already checked in. Please check out first.",
		})
		return
	}

	// Begin transaction
	tx, err := h.db.Begin(ctx)
	if err != nil {
		h.logger.Error("Failed to begin transaction", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal_error"})
		return
	}
	defer tx.Rollback(ctx)

	// Create or update attendance record
	attendanceID := uuid.New().String()
	checkInTime := time.Now()

	_, err = tx.Exec(ctx,
		`INSERT INTO attendance_records (id, user_id, attendance_date, check_in_at, check_in_source, check_in_meta)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 ON CONFLICT (user_id, attendance_date) DO UPDATE SET
		 check_in_at = $4, check_in_source = $5, check_in_meta = $6, updated_at = CURRENT_TIMESTAMP`,
		attendanceID, userIDStr, today, checkInTime, req.Source, req.Metadata,
	)

	if err != nil {
		h.logger.Error("Failed to insert attendance record", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal_error"})
		return
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		h.logger.Error("Failed to commit transaction", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal_error"})
		return
	}

	// Publish event to Kafka
	eventData := map[string]interface{}{
		"user_id":          userIDStr,
		"attendance_date":  today,
		"check_in_time":    checkInTime,
		"check_in_source":  req.Source,
		"check_in_meta":    req.Metadata,
	}

	if err := h.kafkaProducer.PublishCheckInEvent(ctx, "attendance.check-in", userIDStr, attendanceID, eventData); err != nil {
		h.logger.Error("Failed to publish check-in event", "error", err)
		// Don't fail the request, event publishing is async
	}

	h.logger.Info("User checked in successfully", "user_id", userIDStr, "attendance_id", attendanceID)

	c.JSON(http.StatusCreated, AttendanceResponse{
		ID:             attendanceID,
		UserID:         userIDStr,
		AttendanceDate: today,
		CheckInTime:    &checkInTime,
		Status:         "checked_in",
	})
}

// CheckOut godoc
// @Summary Check-out from attendance
// @Description Record user check-out
// @Tags Attendance
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param request body CheckOutRequest true "Check-out data"
// @Success 200 {object} AttendanceResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /attendance/check-out [post]
func (h *Handlers) CheckOut(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	var req CheckOutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userIDStr := userID.(string)
	today := time.Now().Format("2006-01-02")

	// Begin transaction
	tx, err := h.db.Begin(ctx)
	if err != nil {
		h.logger.Error("Failed to begin transaction", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal_error"})
		return
	}
	defer tx.Rollback(ctx)

	// Get attendance record
	var attendanceID string
	var checkInTime time.Time
	err = tx.QueryRow(ctx,
		`SELECT id, check_in_at FROM attendance_records 
		 WHERE user_id = $1 AND attendance_date = $2 AND check_in_at IS NOT NULL AND check_out_at IS NULL
		 FOR UPDATE`,
		userIDStr, today,
	).Scan(&attendanceID, &checkInTime)

	if err != nil {
		h.logger.Warn("No active check-in found", "user_id", userIDStr)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "not_found",
			Message: "No active check-in found. Please check in first.",
		})
		return
	}

	// Update attendance record
	checkOutTime := time.Now()
	_, err = tx.Exec(ctx,
		`UPDATE attendance_records 
		 SET check_out_at = $1, check_out_source = $2, check_out_meta = $3, updated_at = CURRENT_TIMESTAMP
		 WHERE id = $4`,
		checkOutTime, req.Source, req.Metadata, attendanceID,
	)

	if err != nil {
		h.logger.Error("Failed to update attendance record", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal_error"})
		return
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		h.logger.Error("Failed to commit transaction", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal_error"})
		return
	}

	// Publish event to Kafka
	eventData := map[string]interface{}{
		"user_id":          userIDStr,
		"attendance_date":  today,
		"check_in_time":    checkInTime,
		"check_out_time":   checkOutTime,
		"check_out_source": req.Source,
		"check_out_meta":   req.Metadata,
	}

	if err := h.kafkaProducer.PublishCheckOutEvent(ctx, "attendance.check-out", userIDStr, attendanceID, eventData); err != nil {
		h.logger.Error("Failed to publish check-out event", "error", err)
	}

	h.logger.Info("User checked out successfully", "user_id", userIDStr, "attendance_id", attendanceID)

	c.JSON(http.StatusOK, AttendanceResponse{
		ID:             attendanceID,
		UserID:         userIDStr,
		AttendanceDate: today,
		CheckInTime:    &checkInTime,
		CheckOutTime:   &checkOutTime,
		Status:         "checked_out",
	})
}

// GetTodayAttendance godoc
// @Summary Get today's attendance
// @Description Get current user's attendance record for today
// @Tags Attendance
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Success 200 {object} AttendanceResponse
// @Failure 404 {object} ErrorResponse
// @Router /attendance/today [get]
func (h *Handlers) GetTodayAttendance(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userIDStr := userID.(string)
	today := time.Now().Format("2006-01-02")

	var attendance AttendanceResponse
	err := h.db.QueryRow(ctx,
		`SELECT id, user_id, attendance_date, check_in_at, check_out_at
		 FROM attendance_records 
		 WHERE user_id = $1 AND attendance_date = $2`,
		userIDStr, today,
	).Scan(&attendance.ID, &attendance.UserID, &attendance.AttendanceDate, &attendance.CheckInTime, &attendance.CheckOutTime)

	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "not_found"})
		return
	}

	if attendance.CheckOutTime != nil {
		attendance.Status = "checked_out"
	} else if attendance.CheckInTime != nil {
		attendance.Status = "checked_in"
	}

	c.JSON(http.StatusOK, attendance)
}

// GetAttendanceHistory godoc
// @Summary Get attendance history
// @Description Get user's attendance history with pagination
// @Tags Attendance
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Records per page" default(10)
// @Success 200 {object} map[string]interface{}
// @Router /attendance/history [get]
func (h *Handlers) GetAttendanceHistory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userIDStr := userID.(string)
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	rows, err := h.db.Query(ctx,
		`SELECT id, user_id, attendance_date, check_in_at, check_out_at
		 FROM attendance_records 
		 WHERE user_id = $1 AND attendance_date BETWEEN $2 AND $3
		 ORDER BY attendance_date DESC
		 LIMIT $4 OFFSET $5`,
		userIDStr, startDate, endDate, limit, (page + "0"),
	)

	if err != nil {
		h.logger.Error("Failed to query attendance history", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal_error"})
		return
	}
	defer rows.Close()

	var records []AttendanceResponse
	for rows.Next() {
		var record AttendanceResponse
		if err := rows.Scan(&record.ID, &record.UserID, &record.AttendanceDate, &record.CheckInTime, &record.CheckOutTime); err != nil {
			h.logger.Error("Failed to scan attendance record", "error", err)
			continue
		}

		if record.CheckOutTime != nil {
			record.Status = "checked_out"
		} else if record.CheckInTime != nil {
			record.Status = "checked_in"
		}

		records = append(records, record)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  records,
		"count": len(records),
	})
}
