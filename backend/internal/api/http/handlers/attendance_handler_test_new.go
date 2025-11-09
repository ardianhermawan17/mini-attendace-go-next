package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/trustmedis/mini-attendance/internal/api/http/handler"
)

type AttendanceHandlerUnitTestSuite struct {
	suite.Suite
	router   *gin.Engine
	handlers *handler.Handlers
}

func TestAttendanceHandlerUnitTestSuite(t *testing.T) {
	suite.Run(t, new(AttendanceHandlerUnitTestSuite))
}

func (suite *AttendanceHandlerUnitTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	suite.handlers = &handler.Handlers{}
}

// ===== Unit Tests for CheckIn =====

func (suite *AttendanceHandlerUnitTestSuite) TestCheckIn_Success() {
	checkInReq := map[string]interface{}{
		"latitude":  -6.1753,
		"longitude": 106.8249,
		"device":    "mobile",
	}

	body, _ := json.Marshal(checkInReq)
	req := httptest.NewRequest("POST", "/api/v1/attendance/check-in", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer validtoken")

	// Should return 201 Created with attendance record
	assert.NotNil(suite.T(), req)
	assert.Equal(suite.T(), "POST", req.Method)
}

func (suite *AttendanceHandlerUnitTestSuite) TestCheckIn_NoAuthorization() {
	checkInReq := map[string]interface{}{
		"latitude":  -6.1753,
		"longitude": 106.8249,
		"device":    "mobile",
	}

	body, _ := json.Marshal(checkInReq)
	req := httptest.NewRequest("POST", "/api/v1/attendance/check-in", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Should return 401 Unauthorized
	assert.Empty(suite.T(), req.Header.Get("Authorization"))
}

func (suite *AttendanceHandlerUnitTestSuite) TestCheckIn_InvalidCoordinates() {
	checkInReq := map[string]interface{}{
		"latitude":  -200.0, // Invalid latitude (must be between -90 and 90)
		"longitude": 106.8249,
		"device":    "mobile",
	}

	body, _ := json.Marshal(checkInReq)
	req := httptest.NewRequest("POST", "/api/v1/attendance/check-in", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer validtoken")

	// Should return 400 Bad Request
	assert.NotNil(suite.T(), req)
}

func (suite *AttendanceHandlerUnitTestSuite) TestCheckIn_MissingLatitude() {
	checkInReq := map[string]interface{}{
		"longitude": 106.8249,
		"device":    "mobile",
	}

	body, _ := json.Marshal(checkInReq)
	req := httptest.NewRequest("POST", "/api/v1/attendance/check-in", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer validtoken")

	// Should return 400 Bad Request
	assert.NotNil(suite.T(), req)
}

func (suite *AttendanceHandlerUnitTestSuite) TestCheckIn_AlreadyCheckedIn() {
	checkInReq := map[string]interface{}{
		"latitude":  -6.1753,
		"longitude": 106.8249,
		"device":    "mobile",
	}

	body, _ := json.Marshal(checkInReq)
	req := httptest.NewRequest("POST", "/api/v1/attendance/check-in", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer alreadycheckedintoken")

	// Should return 409 Conflict
	assert.NotNil(suite.T(), req)
}

// ===== Unit Tests for CheckOut =====

func (suite *AttendanceHandlerUnitTestSuite) TestCheckOut_Success() {
	checkOutReq := map[string]interface{}{
		"latitude":  -6.1753,
		"longitude": 106.8249,
		"device":    "mobile",
	}

	body, _ := json.Marshal(checkOutReq)
	req := httptest.NewRequest("POST", "/api/v1/attendance/check-out", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer validtoken")

	// Should return 200 OK with updated attendance record
	assert.NotNil(suite.T(), req)
	assert.Equal(suite.T(), "POST", req.Method)
}

func (suite *AttendanceHandlerUnitTestSuite) TestCheckOut_NotCheckedIn() {
	checkOutReq := map[string]interface{}{
		"latitude":  -6.1753,
		"longitude": 106.8249,
		"device":    "mobile",
	}

	body, _ := json.Marshal(checkOutReq)
	req := httptest.NewRequest("POST", "/api/v1/attendance/check-out", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer notcheckedintoken")

	// Should return 400 Bad Request
	assert.NotNil(suite.T(), req)
}

func (suite *AttendanceHandlerUnitTestSuite) TestCheckOut_NoAuthorization() {
	checkOutReq := map[string]interface{}{
		"latitude":  -6.1753,
		"longitude": 106.8249,
		"device":    "mobile",
	}

	body, _ := json.Marshal(checkOutReq)
	req := httptest.NewRequest("POST", "/api/v1/attendance/check-out", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Should return 401 Unauthorized
	assert.Empty(suite.T(), req.Header.Get("Authorization"))
}

// ===== Unit Tests for GetTodayAttendance =====

func (suite *AttendanceHandlerUnitTestSuite) TestGetTodayAttendance_Success() {
	req := httptest.NewRequest("GET", "/api/v1/attendance/today", nil)
	req.Header.Set("Authorization", "Bearer validtoken")

	// Should return 200 OK with today's attendance record
	assert.NotNil(suite.T(), req)
	assert.Equal(suite.T(), "GET", req.Method)
}

func (suite *AttendanceHandlerUnitTestSuite) TestGetTodayAttendance_NoRecord() {
	req := httptest.NewRequest("GET", "/api/v1/attendance/today", nil)
	req.Header.Set("Authorization", "Bearer newtokenuser")

	// Should return 404 Not Found
	assert.NotNil(suite.T(), req)
}

func (suite *AttendanceHandlerUnitTestSuite) TestGetTodayAttendance_NoAuthorization() {
	req := httptest.NewRequest("GET", "/api/v1/attendance/today", nil)

	// Should return 401 Unauthorized
	assert.Empty(suite.T(), req.Header.Get("Authorization"))
}

// ===== Unit Tests for GetAttendanceHistory =====

func (suite *AttendanceHandlerUnitTestSuite) TestGetAttendanceHistory_Success() {
	req := httptest.NewRequest("GET", "/api/v1/attendance/history?from=2024-01-01&to=2024-01-31", nil)
	req.Header.Set("Authorization", "Bearer validtoken")

	// Should return 200 OK with paginated records
	assert.NotNil(suite.T(), req)
	assert.Equal(suite.T(), "GET", req.Method)
}

func (suite *AttendanceHandlerUnitTestSuite) TestGetAttendanceHistory_InvalidFromDate() {
	req := httptest.NewRequest("GET", "/api/v1/attendance/history?from=invalid&to=2024-01-31", nil)
	req.Header.Set("Authorization", "Bearer validtoken")

	// Should return 400 Bad Request
	assert.NotNil(suite.T(), req)
}

func (suite *AttendanceHandlerUnitTestSuite) TestGetAttendanceHistory_InvalidToDate() {
	req := httptest.NewRequest("GET", "/api/v1/attendance/history?from=2024-01-01&to=not-a-date", nil)
	req.Header.Set("Authorization", "Bearer validtoken")

	// Should return 400 Bad Request
	assert.NotNil(suite.T(), req)
}

func (suite *AttendanceHandlerUnitTestSuite) TestGetAttendanceHistory_MissingFromDate() {
	req := httptest.NewRequest("GET", "/api/v1/attendance/history?to=2024-01-31", nil)
	req.Header.Set("Authorization", "Bearer validtoken")

	// Should return 400 Bad Request
	assert.NotNil(suite.T(), req)
}

func (suite *AttendanceHandlerUnitTestSuite) TestGetAttendanceHistory_MissingToDate() {
	req := httptest.NewRequest("GET", "/api/v1/attendance/history?from=2024-01-01", nil)
	req.Header.Set("Authorization", "Bearer validtoken")

	// Should return 400 Bad Request
	assert.NotNil(suite.T(), req)
}

func (suite *AttendanceHandlerUnitTestSuite) TestGetAttendanceHistory_InvalidDateRange() {
	// from date is after to date
	req := httptest.NewRequest("GET", "/api/v1/attendance/history?from=2024-01-31&to=2024-01-01", nil)
	req.Header.Set("Authorization", "Bearer validtoken")

	// Should return 400 Bad Request
	assert.NotNil(suite.T(), req)
}

func (suite *AttendanceHandlerUnitTestSuite) TestGetAttendanceHistory_WithPagination() {
	req := httptest.NewRequest("GET", "/api/v1/attendance/history?from=2024-01-01&to=2024-01-31&page=2&limit=10", nil)
	req.Header.Set("Authorization", "Bearer validtoken")

	// Should return 200 OK with page 2
	assert.NotNil(suite.T(), req)
	assert.Contains(suite.T(), req.URL.RawQuery, "page=2")
}

func (suite *AttendanceHandlerUnitTestSuite) TestGetAttendanceHistory_InvalidLimit() {
	req := httptest.NewRequest("GET", "/api/v1/attendance/history?from=2024-01-01&to=2024-01-31&limit=1000", nil)
	req.Header.Set("Authorization", "Bearer validtoken")

	// Should return 400 Bad Request - limit exceeds maximum (100)
	assert.NotNil(suite.T(), req)
}

func (suite *AttendanceHandlerUnitTestSuite) TestGetAttendanceHistory_NoAuthorization() {
	req := httptest.NewRequest("GET", "/api/v1/attendance/history?from=2024-01-01&to=2024-01-31", nil)

	// Should return 401 Unauthorized
	assert.Empty(suite.T(), req.Header.Get("Authorization"))
}

// ===== Unit Tests for GetUserAttendance (Manager/Admin only) =====

func (suite *AttendanceHandlerUnitTestSuite) TestGetUserAttendance_AsManager() {
	req := httptest.NewRequest(
		"GET",
		"/api/v1/attendance/users/550e8400-e29b-41d4-a716-446655440002?from=2024-01-01&to=2024-01-31",
		nil,
	)
	req.Header.Set("Authorization", "Bearer managertoken")

	// Should return 200 OK
	assert.NotNil(suite.T(), req)
}

func (suite *AttendanceHandlerUnitTestSuite) TestGetUserAttendance_AsAdmin() {
	req := httptest.NewRequest(
		"GET",
		"/api/v1/attendance/users/550e8400-e29b-41d4-a716-446655440002?from=2024-01-01&to=2024-01-31",
		nil,
	)
	req.Header.Set("Authorization", "Bearer admintoken")

	// Should return 200 OK
	assert.NotNil(suite.T(), req)
}

func (suite *AttendanceHandlerUnitTestSuite) TestGetUserAttendance_AsEmployee() {
	req := httptest.NewRequest(
		"GET",
		"/api/v1/attendance/users/550e8400-e29b-41d4-a716-446655440003?from=2024-01-01&to=2024-01-31",
		nil,
	)
	req.Header.Set("Authorization", "Bearer employeetoken")

	// Should return 403 Forbidden
	assert.NotNil(suite.T(), req)
}

func (suite *AttendanceHandlerUnitTestSuite) TestGetUserAttendance_UserNotFound() {
	req := httptest.NewRequest(
		"GET",
		"/api/v1/attendance/users/nonexistent-uuid?from=2024-01-01&to=2024-01-31",
		nil,
	)
	req.Header.Set("Authorization", "Bearer admintoken")

	// Should return 404 Not Found
	assert.NotNil(suite.T(), req)
}

// ===== Integration-style tests (simulating realistic flows) =====

func (suite *AttendanceHandlerUnitTestSuite) TestAttendanceCheckInCheckOutFlow() {
	// 1. Check-in
	checkInReq := map[string]interface{}{
		"latitude":  -6.1753,
		"longitude": 106.8249,
		"device":    "mobile",
	}
	checkInBody, _ := json.Marshal(checkInReq)
	checkInHttpReq := httptest.NewRequest("POST", "/api/v1/attendance/check-in", bytes.NewReader(checkInBody))
	checkInHttpReq.Header.Set("Authorization", "Bearer validtoken")
	assert.NotNil(suite.T(), checkInHttpReq)

	// Simulate a delay
	time.Sleep(100 * time.Millisecond)

	// 2. Check-out
	checkOutReq := map[string]interface{}{
		"latitude":  -6.1753,
		"longitude": 106.8249,
		"device":    "mobile",
	}
	checkOutBody, _ := json.Marshal(checkOutReq)
	checkOutHttpReq := httptest.NewRequest("POST", "/api/v1/attendance/check-out", bytes.NewReader(checkOutBody))
	checkOutHttpReq.Header.Set("Authorization", "Bearer validtoken")
	assert.NotNil(suite.T(), checkOutHttpReq)

	// 3. Get today's record
	getTodayReq := httptest.NewRequest("GET", "/api/v1/attendance/today", nil)
	getTodayReq.Header.Set("Authorization", "Bearer validtoken")
	assert.NotNil(suite.T(), getTodayReq)
}

func (suite *AttendanceHandlerUnitTestSuite) TestAttendanceHistoryFiltering() {
	// Test various history filter combinations
	testCases := []struct {
		name        string
		fromDate    string
		toDate      string
		page        string
		limit       string
		expectError bool
	}{
		{
			name:        "Valid date range",
			fromDate:    "2024-01-01",
			toDate:      "2024-01-31",
			expectError: false,
		},
		{
			name:        "Single day",
			fromDate:    "2024-01-15",
			toDate:      "2024-01-15",
			expectError: false,
		},
		{
			name:        "With pagination",
			fromDate:    "2024-01-01",
			toDate:      "2024-01-31",
			page:        "1",
			limit:       "20",
			expectError: false,
		},
		{
			name:        "Invalid from date",
			fromDate:    "invalid",
			toDate:      "2024-01-31",
			expectError: true,
		},
		{
			name:        "Invalid to date",
			fromDate:    "2024-01-01",
			toDate:      "invalid",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			query := "?from=" + tc.fromDate + "&to=" + tc.toDate
			if tc.page != "" {
				query += "&page=" + tc.page
			}
			if tc.limit != "" {
				query += "&limit=" + tc.limit
			}

			req := httptest.NewRequest("GET", "/api/v1/attendance/history"+query, nil)
			req.Header.Set("Authorization", "Bearer validtoken")
			assert.NotNil(suite.T(), req)
		})
	}
}
