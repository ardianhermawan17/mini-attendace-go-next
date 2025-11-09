package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// E2ETestSuite contains all end-to-end integration tests
type E2ETestSuite struct {
	suite.Suite
	baseURL    string
	httpClient *http.Client
	tokens     map[string]string // Store tokens for different users
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}

func (suite *E2ETestSuite) SetupSuite() {
	// Initialize test environment
	suite.baseURL = "http://localhost:8080/api/v1"
	suite.httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}
	suite.tokens = make(map[string]string)

	// Wait for server to be ready
	suite.waitForServer()

	// Setup test data (register/login test users)
	suite.setupTestUsers()
}

func (suite *E2ETestSuite) TearDownSuite() {
	// Cleanup test data if needed
}

// ===== Helper Methods =====

func (suite *E2ETestSuite) waitForServer() {
	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		resp, err := suite.httpClient.Get(suite.baseURL + "/health")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(time.Second)
	}
	suite.T().Fatal("Server not ready after 30 seconds")
}

func (suite *E2ETestSuite) setupTestUsers() {
	// Register test users if they don't exist
	testUsers := []map[string]string{
		{
			"email":     "admin@trustmedis.com",
			"password":  "admin123",
			"full_name": "Admin User",
			"role":      "admin",
		},
		{
			"email":     "manager@trustmedis.com",
			"password":  "password123",
			"full_name": "Manager User",
			"role":      "manager",
		},
		{
			"email":     "employee@trustmedis.com",
			"password":  "password123",
			"full_name": "Employee User",
			"role":      "employee",
		},
	}

	for _, user := range testUsers {
		// Try to login first
		loginToken := suite.login(user["email"], user["password"])
		if loginToken != "" {
			suite.tokens[user["email"]] = loginToken
			continue
		}

		// If login fails, try to register
		registerResp := suite.registerUser(user)
		if registerResp != nil {
			// Then login
			loginToken := suite.login(user["email"], user["password"])
			suite.tokens[user["email"]] = loginToken
		}
	}
}

func (suite *E2ETestSuite) registerUser(user map[string]string) map[string]interface{} {
	body, _ := json.Marshal(user)
	resp, err := suite.httpClient.Post(
		suite.baseURL+"/auth/register",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result
}

func (suite *E2ETestSuite) login(email, password string) string {
	loginReq := map[string]string{
		"email":    email,
		"password": password,
	}
	body, _ := json.Marshal(loginReq)
	resp, err := suite.httpClient.Post(
		suite.baseURL+"/auth/login",
		"application/json",
		bytes.NewReader(body),
	)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	token, ok := result["access_token"].(string)
	if !ok {
		return ""
	}
	return token
}

func (suite *E2ETestSuite) makeRequest(method, endpoint string, body interface{}, userEmail string) *http.Response {
	var reqBody io.Reader
	if body != nil {
		bodyBytes, _ := json.Marshal(body)
		reqBody = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, suite.baseURL+endpoint, reqBody)
	require.NoError(suite.T(), err)

	req.Header.Set("Content-Type", "application/json")
	if userEmail != "" {
		token := suite.tokens[userEmail]
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	resp, err := suite.httpClient.Do(req)
	require.NoError(suite.T(), err)

	return resp
}

func (suite *E2ETestSuite) parseResponse(resp *http.Response, v interface{}) {
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(v)
}

// ===== Authentication E2E Tests =====

func (suite *E2ETestSuite) TestAuthenticationFlow() {
	// 1. Register new user
	registerReq := map[string]string{
		"email":     "e2etest@trustmedis.com",
		"password":  "password123",
		"full_name": "E2E Test User",
	}

	resp := suite.makeRequest("POST", "/auth/register", registerReq, "")
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var registerResp map[string]interface{}
	suite.parseResponse(resp, &registerResp)
	assert.NotNil(suite.T(), registerResp["id"])

	// 2. Login with registered credentials
	loginReq := map[string]string{
		"email":    "e2etest@trustmedis.com",
		"password": "password123",
	}

	resp = suite.makeRequest("POST", "/auth/login", loginReq, "")
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var loginResp map[string]interface{}
	suite.parseResponse(resp, &loginResp)
	assert.NotEmpty(suite.T(), loginResp["access_token"])
	assert.NotEmpty(suite.T(), loginResp["refresh_token"])

	accessToken := loginResp["access_token"].(string)
	suite.tokens["e2etest@trustmedis.com"] = accessToken
}

func (suite *E2ETestSuite) TestLoginWithInvalidCredentials() {
	loginReq := map[string]string{
		"email":    "admin@trustmedis.com",
		"password": "wrongpassword",
	}

	resp := suite.makeRequest("POST", "/auth/login", loginReq, "")
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
	resp.Body.Close()
}

func (suite *E2ETestSuite) TestLoginWithNonExistentUser() {
	loginReq := map[string]string{
		"email":    "nonexistent@trustmedis.com",
		"password": "password123",
	}

	resp := suite.makeRequest("POST", "/auth/login", loginReq, "")
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
	resp.Body.Close()
}

func (suite *E2ETestSuite) TestRefreshToken() {
	// First login
	loginReq := map[string]string{
		"email":    "employee@trustmedis.com",
		"password": "password123",
	}

	resp := suite.makeRequest("POST", "/auth/login", loginReq, "")
	var loginResp map[string]interface{}
	suite.parseResponse(resp, &loginResp)

	refreshToken := loginResp["refresh_token"].(string)

	// Now refresh the token
	req, _ := http.NewRequest("POST", suite.baseURL+"/auth/refresh", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", refreshToken))

	resp, _ = suite.httpClient.Do(req)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var refreshResp map[string]interface{}
	suite.parseResponse(resp, &refreshResp)
	assert.NotEmpty(suite.T(), refreshResp["access_token"])
}

// ===== Attendance E2E Tests =====

func (suite *E2ETestSuite) TestCompleteAttendanceFlow() {
	userEmail := "employee@trustmedis.com"

	// 1. Check-in
	checkInReq := map[string]interface{}{
		"latitude":  -6.1753,
		"longitude": 106.8249,
		"device":    "mobile",
	}

	resp := suite.makeRequest("POST", "/attendance/check-in", checkInReq, userEmail)
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var checkInResp map[string]interface{}
	suite.parseResponse(resp, &checkInResp)
	assert.NotNil(suite.T(), checkInResp["id"])
	assert.NotNil(suite.T(), checkInResp["check_in_time"])

	// 2. Get today's attendance
	resp = suite.makeRequest("GET", "/attendance/today", nil, userEmail)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var todayResp map[string]interface{}
	suite.parseResponse(resp, &todayResp)
	assert.NotNil(suite.T(), todayResp["check_in_time"])
	assert.Nil(suite.T(), todayResp["check_out_time"])

	// 3. Wait a moment
	time.Sleep(500 * time.Millisecond)

	// 4. Check-out
	checkOutReq := map[string]interface{}{
		"latitude":  -6.1753,
		"longitude": 106.8249,
		"device":    "mobile",
	}

	resp = suite.makeRequest("POST", "/attendance/check-out", checkOutReq, userEmail)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var checkOutResp map[string]interface{}
	suite.parseResponse(resp, &checkOutResp)
	assert.NotNil(suite.T(), checkOutResp["check_out_time"])

	// 5. Get updated today's attendance
	resp = suite.makeRequest("GET", "/attendance/today", nil, userEmail)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var finalResp map[string]interface{}
	suite.parseResponse(resp, &finalResp)
	assert.NotNil(suite.T(), finalResp["check_in_time"])
	assert.NotNil(suite.T(), finalResp["check_out_time"])
}

func (suite *E2ETestSuite) TestAttendanceCheckInTwice() {
	userEmail := "employee@trustmedis.com"

	// First check-in should succeed
	checkInReq := map[string]interface{}{
		"latitude":  -6.1753,
		"longitude": 106.8249,
		"device":    "mobile",
	}

	resp := suite.makeRequest("POST", "/attendance/check-in", checkInReq, userEmail)
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	// Second check-in should fail (already checked in)
	resp = suite.makeRequest("POST", "/attendance/check-in", checkInReq, userEmail)
	assert.Equal(suite.T(), http.StatusConflict, resp.StatusCode)
	resp.Body.Close()
}

func (suite *E2ETestSuite) TestAttendanceCheckOutWithoutCheckIn() {
	userEmail := "manager@trustmedis.com"

	// Try to check-out without checking in first
	checkOutReq := map[string]interface{}{
		"latitude":  -6.1753,
		"longitude": 106.8249,
		"device":    "mobile",
	}

	resp := suite.makeRequest("POST", "/attendance/check-out", checkOutReq, userEmail)
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
	resp.Body.Close()
}

func (suite *E2ETestSuite) TestGetAttendanceHistory() {
	userEmail := "employee@trustmedis.com"

	// Get attendance history for last 30 days
	resp := suite.makeRequest(
		"GET",
		"/attendance/history?from=2024-01-01&to=2024-01-31",
		nil,
		userEmail,
	)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var historyResp map[string]interface{}
	suite.parseResponse(resp, &historyResp)
	assert.NotNil(suite.T(), historyResp["data"])
	assert.NotNil(suite.T(), historyResp["pagination"])
}

func (suite *E2ETestSuite) TestGetAttendanceHistoryWithInvalidDates() {
	userEmail := "employee@trustmedis.com"

	// Invalid from date
	resp := suite.makeRequest(
		"GET",
		"/attendance/history?from=invalid&to=2024-01-31",
		nil,
		userEmail,
	)
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
	resp.Body.Close()

	// Invalid to date
	resp = suite.makeRequest(
		"GET",
		"/attendance/history?from=2024-01-01&to=invalid",
		nil,
		userEmail,
	)
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
	resp.Body.Close()
}

// ===== Authorization E2E Tests =====

func (suite *E2ETestSuite) TestEmployeeCantViewOtherAttendance() {
	employeeEmail := "employee@trustmedis.com"
	anotherUserID := "550e8400-e29b-41d4-a716-446655440003"

	resp := suite.makeRequest(
		"GET",
		fmt.Sprintf("/attendance/users/%s?from=2024-01-01&to=2024-01-31", anotherUserID),
		nil,
		employeeEmail,
	)
	assert.Equal(suite.T(), http.StatusForbidden, resp.StatusCode)
	resp.Body.Close()
}

func (suite *E2ETestSuite) TestManagerCanViewEmployeeAttendance() {
	managerEmail := "manager@trustmedis.com"
	employeeID := "550e8400-e29b-41d4-a716-446655440002"

	resp := suite.makeRequest(
		"GET",
		fmt.Sprintf("/attendance/users/%s?from=2024-01-01&to=2024-01-31", employeeID),
		nil,
		managerEmail,
	)
	// Should be 200 OK or 404 Not Found (if employee doesn't have records)
	assert.Contains(suite.T(), []int{http.StatusOK, http.StatusNotFound}, resp.StatusCode)
	resp.Body.Close()
}

func (suite *E2ETestSuite) TestAdminCanViewAnyAttendance() {
	adminEmail := "admin@trustmedis.com"
	anyUserID := "550e8400-e29b-41d4-a716-446655440002"

	resp := suite.makeRequest(
		"GET",
		fmt.Sprintf("/attendance/users/%s?from=2024-01-01&to=2024-01-31", anyUserID),
		nil,
		adminEmail,
	)
	assert.Contains(suite.T(), []int{http.StatusOK, http.StatusNotFound}, resp.StatusCode)
	resp.Body.Close()
}

// ===== Report E2E Tests =====

func (suite *E2ETestSuite) TestGetMonthlyReport() {
	adminEmail := "admin@trustmedis.com"

	resp := suite.makeRequest(
		"GET",
		"/reports/monthly?year=2024&month=1",
		nil,
		adminEmail,
	)
	// Should return 200 or 404 if no data
	assert.Contains(suite.T(), []int{http.StatusOK, http.StatusNotFound}, resp.StatusCode)
	resp.Body.Close()
}

func (suite *E2ETestSuite) TestEmployeeCantAccessMonthlyReport() {
	employeeEmail := "employee@trustmedis.com"

	resp := suite.makeRequest(
		"GET",
		"/reports/monthly?year=2024&month=1",
		nil,
		employeeEmail,
	)
	assert.Equal(suite.T(), http.StatusForbidden, resp.StatusCode)
	resp.Body.Close()
}

func (suite *E2ETestSuite) TestGetReportWithoutAuthentication() {
	// Try to access report without token
	req, _ := http.NewRequest("GET", suite.baseURL+"/reports/monthly?year=2024&month=1", nil)
	resp, _ := suite.httpClient.Do(req)

	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
	resp.Body.Close()
}

// ===== Stress/Load Tests =====

func (suite *E2ETestSuite) TestMultipleCheckInsAndOuts() {
	userEmail := "employee@trustmedis.com"

	for i := 0; i < 3; i++ {
		// Check-out if already checked in
		checkOutReq := map[string]interface{}{
			"latitude":  -6.1753,
			"longitude": 106.8249,
			"device":    "mobile",
		}
		suite.makeRequest("POST", "/attendance/check-out", checkOutReq, userEmail)

		// Check-in
		checkInReq := map[string]interface{}{
			"latitude":  -6.1753,
			"longitude": 106.8249,
			"device":    "mobile",
		}
		resp := suite.makeRequest("POST", "/attendance/check-in", checkInReq, userEmail)
		assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)
		resp.Body.Close()

		time.Sleep(100 * time.Millisecond)
	}
}

// ===== Data Consistency E2E Tests =====

func (suite *E2ETestSuite) TestAttendanceDataConsistency() {
	userEmail := "employee@trustmedis.com"

	// 1. Check-in
	checkInReq := map[string]interface{}{
		"latitude":  -6.1753,
		"longitude": 106.8249,
		"device":    "mobile",
	}
	resp := suite.makeRequest("POST", "/attendance/check-in", checkInReq, userEmail)
	var checkInResp map[string]interface{}
	suite.parseResponse(resp, &checkInResp)
	recordID := checkInResp["id"]

	// 2. Get today's record
	resp = suite.makeRequest("GET", "/attendance/today", nil, userEmail)
	var todayResp map[string]interface{}
	suite.parseResponse(resp, &todayResp)

	// Verify it's the same record
	assert.Equal(suite.T(), recordID, todayResp["id"])

	// 3. Check-out and verify
	checkOutReq := map[string]interface{}{
		"latitude":  -6.1753,
		"longitude": 106.8249,
		"device":    "mobile",
	}
	resp = suite.makeRequest("POST", "/attendance/check-out", checkOutReq, userEmail)
	var checkOutResp map[string]interface{}
	suite.parseResponse(resp, &checkOutResp)

	// Verify same record with check-out time
	assert.Equal(suite.T(), recordID, checkOutResp["id"])
	assert.NotNil(suite.T(), checkOutResp["check_out_time"])
}

// ===== Error Handling E2E Tests =====

func (suite *E2ETestSuite) TestInvalidJSONRequest() {
	userEmail := "employee@trustmedis.com"
	token := suite.tokens[userEmail]

	req, _ := http.NewRequest("POST", suite.baseURL+"/attendance/check-in", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, _ := suite.httpClient.Do(req)
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
	resp.Body.Close()
}

func (suite *E2ETestSuite) TestMissingRequiredFields() {
	userEmail := "employee@trustmedis.com"

	// Check-in without latitude
	checkInReq := map[string]interface{}{
		"longitude": 106.8249,
		"device":    "mobile",
	}
	resp := suite.makeRequest("POST", "/attendance/check-in", checkInReq, userEmail)
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
	resp.Body.Close()
}

func (suite *E2ETestSuite) TestConcurrentRequests() {
	userEmail := "employee@trustmedis.com"

	// Make multiple concurrent requests
	done := make(chan bool, 5)

	for i := 0; i < 5; i++ {
		go func() {
			resp := suite.makeRequest("GET", "/attendance/today", nil, userEmail)
			assert.Contains(suite.T(), []int{http.StatusOK, http.StatusNotFound}, resp.StatusCode)
			resp.Body.Close()
			done <- true
		}()
	}

	for i := 0; i < 5; i++ {
		<-done
	}
}
