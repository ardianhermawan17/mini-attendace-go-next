package handler

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AbsenceReportResponse struct {
	ID           string     `json:"id"`
	UserID       string     `json:"user_id"`
	ReportDate   string     `json:"report_date"`
	Status       string     `json:"status"`
	CheckInTime  *time.Time `json:"check_in_time"`
	CheckOutTime *time.Time `json:"check_out_time"`
	WorkHours    *float64   `json:"work_hours"`
	Notes        *string    `json:"notes"`
}

// GetAbsenceReport godoc
// @Summary Get absence report
// @Description Get absence report with filtering options
// @Tags Reports
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param status query string false "Filter by status (hadir, terlambat, pulang_cepat, tidak_hadir)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Records per page" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /reports/absence [get]
func (h *Handlers) GetAbsenceReport(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userIDStr := userID.(string)
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, -3, 0).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))
	status := c.Query("status")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	// Parse pagination parameters
	page := 1
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	limit := 20
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
		limit = l
	}

	offset := (page - 1) * limit

	// Build query
	query := `SELECT id, user_id, report_date, status, check_in_time, check_out_time, work_hours, notes
	         FROM absence_reports 
	         WHERE user_id = $1 AND report_date BETWEEN $2 AND $3`
	args := []interface{}{userIDStr, startDate, endDate}

	if status != "" {
		query += " AND status = $4"
		args = append(args, status)
	}

	query += " ORDER BY report_date DESC LIMIT $" + fmt.Sprintf("%d", len(args)+1) + " OFFSET $" + fmt.Sprintf("%d", len(args)+2)
	args = append(args, limit, offset)

	rows, err := h.db.Query(ctx, query, args...)
	if err != nil {
		h.logger.Error("Failed to query absence report", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal_error"})
		return
	}
	defer rows.Close()

	var reports []AbsenceReportResponse
	for rows.Next() {
		var report AbsenceReportResponse
		if err := rows.Scan(&report.ID, &report.UserID, &report.ReportDate, &report.Status,
			&report.CheckInTime, &report.CheckOutTime, &report.WorkHours, &report.Notes); err != nil {
			h.logger.Error("Failed to scan report", "error", err)
			continue
		}
		reports = append(reports, report)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  reports,
		"count": len(reports),
	})
}

// ExportAbsenceReport godoc
// @Summary Export absence report as CSV
// @Description Export absence report data as CSV file
// @Tags Reports
// @Produce text/csv
// @Param Authorization header string true "Bearer {token}"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {file} file
// @Router /reports/absence/export [get]
func (h *Handlers) ExportAbsenceReport(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userIDStr := userID.(string)
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, -3, 0).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	rows, err := h.db.Query(ctx,
		`SELECT report_date, status, check_in_time, check_out_time, work_hours, notes
		 FROM absence_reports 
		 WHERE user_id = $1 AND report_date BETWEEN $2 AND $3
		 ORDER BY report_date DESC`,
		userIDStr, startDate, endDate,
	)

	if err != nil {
		h.logger.Error("Failed to query absence report for export", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal_error"})
		return
	}
	defer rows.Close()

	// Set response headers for CSV download
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=absence_report_%s.csv", time.Now().Format("2006-01-02")))

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Write header
	header := []string{"Date", "Status", "Check-in", "Check-out", "Work Hours", "Notes"}
	if err := writer.Write(header); err != nil {
		h.logger.Error("Failed to write CSV header", "error", err)
		return
	}

	// Write data
	for rows.Next() {
		var reportDate, status string
		var checkInTime, checkOutTime *time.Time
		var workHours *float64
		var notes *string

		if err := rows.Scan(&reportDate, &status, &checkInTime, &checkOutTime, &workHours, &notes); err != nil {
			h.logger.Error("Failed to scan report row", "error", err)
			continue
		}

		checkInStr := ""
		if checkInTime != nil {
			checkInStr = checkInTime.Format("15:04:05")
		}

		checkOutStr := ""
		if checkOutTime != nil {
			checkOutStr = checkOutTime.Format("15:04:05")
		}

		workHoursStr := ""
		if workHours != nil {
			workHoursStr = fmt.Sprintf("%.2f", *workHours)
		}

		notesStr := ""
		if notes != nil {
			notesStr = *notes
		}

		record := []string{reportDate, status, checkInStr, checkOutStr, workHoursStr, notesStr}
		if err := writer.Write(record); err != nil {
			h.logger.Error("Failed to write CSV record", "error", err)
			continue
		}
	}
}
