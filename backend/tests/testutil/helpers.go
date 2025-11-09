package testutil

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trustmedis/mini-attendance/internal/infra/kafka"
	"github.com/trustmedis/mini-attendance/internal/infra/observability"
	"github.com/trustmedis/mini-attendance/internal/infra/redis"
)

// TestHelpers provides utility functions for testing
type TestHelpers struct {
	T *testing.T
}

// NewTestHelpers creates a new test helpers instance
func NewTestHelpers(t *testing.T) *TestHelpers {
	return &TestHelpers{T: t}
}

// ===== Request/Response Helpers =====

// MakeRequest creates an HTTP request for testing
func (h *TestHelpers) MakeRequest(method, path string, body interface{}) *http.Request {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, _ := json.Marshal(body)
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, path, bodyReader)
	if err != nil {
		h.T.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	return req
}

// SetAuthHeader adds Authorization header to request
func (h *TestHelpers) SetAuthHeader(req *http.Request, token string) {
	req.Header.Set("Authorization", "Bearer "+token)
}

// ParseResponse parses JSON response body into a value
func (h *TestHelpers) ParseResponse(resp *http.Response, v interface{}) {
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil && err != io.EOF {
		h.T.Fatalf("Failed to parse response: %v", err)
	}
}

// ===== Test Data Builders =====

// UserBuilder helps build test user data
type UserBuilder struct {
	Email    string
	Password string
	FullName string
	Role     string
}

// NewUserBuilder creates a new user builder with defaults
func NewUserBuilder() *UserBuilder {
	return &UserBuilder{
		Email:    "test@trustmedis.com",
		Password: "password123",
		FullName: "Test User",
		Role:     "employee",
	}
}

// WithEmail sets the email
func (ub *UserBuilder) WithEmail(email string) *UserBuilder {
	ub.Email = email
	return ub
}

// WithPassword sets the password
func (ub *UserBuilder) WithPassword(password string) *UserBuilder {
	ub.Password = password
	return ub
}

// WithFullName sets the full name
func (ub *UserBuilder) WithFullName(fullName string) *UserBuilder {
	ub.FullName = fullName
	return ub
}

// WithRole sets the role
func (ub *UserBuilder) WithRole(role string) *UserBuilder {
	ub.Role = role
	return ub
}

// Build returns the built user data
func (ub *UserBuilder) Build() map[string]string {
	return map[string]string{
		"email":     ub.Email,
		"password":  ub.Password,
		"full_name": ub.FullName,
		"role":      ub.Role,
	}
}

// AttendanceBuilder helps build test attendance data
type AttendanceBuilder struct {
	Latitude  float64
	Longitude float64
	Device    string
}

// NewAttendanceBuilder creates a new attendance builder with defaults
func NewAttendanceBuilder() *AttendanceBuilder {
	return &AttendanceBuilder{
		Latitude:  -6.1753,
		Longitude: 106.8249,
		Device:    "mobile",
	}
}

// WithCoordinates sets the coordinates
func (ab *AttendanceBuilder) WithCoordinates(lat, lon float64) *AttendanceBuilder {
	ab.Latitude = lat
	ab.Longitude = lon
	return ab
}

// WithDevice sets the device
func (ab *AttendanceBuilder) WithDevice(device string) *AttendanceBuilder {
	ab.Device = device
	return ab
}

// Build returns the built attendance data
func (ab *AttendanceBuilder) Build() map[string]interface{} {
	return map[string]interface{}{
		"latitude":  ab.Latitude,
		"longitude": ab.Longitude,
		"device":    ab.Device,
	}
}

// ===== Mock Factories =====

// TokenMockFactory creates mock JWT tokens for testing
type TokenMockFactory struct {
	Secret string
}

// NewTokenMockFactory creates a new token factory
func NewTokenMockFactory(secret string) *TokenMockFactory {
	return &TokenMockFactory{
		Secret: secret,
	}
}

// CreateMockAccessToken creates a mock access token
func (tmf *TokenMockFactory) CreateMockAccessToken(userID, role string) string {
	// In production, this would use actual JWT generation
	// For testing, we can return a fixed format token
	return "mock.access.token." + userID
}

// CreateMockRefreshToken creates a mock refresh token
func (tmf *TokenMockFactory) CreateMockRefreshToken(userID string) string {
	return "mock.refresh.token." + userID
}

// ===== Test Router Setup =====

// SetupTestRouter creates a Gin router for testing
func SetupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// ===== Database Test Helpers =====

// DatabaseTestHelper provides database utilities for tests
type DatabaseTestHelper struct {
	T    *testing.T
	Pool *pgxpool.Pool
}

// NewDatabaseTestHelper creates a new database test helper
func NewDatabaseTestHelper(t *testing.T, pool *pgxpool.Pool) *DatabaseTestHelper {
	return &DatabaseTestHelper{
		T:    t,
		Pool: pool,
	}
}

// CleanupTestData removes test data from database
func (dth *DatabaseTestHelper) CleanupTestData(ctx context.Context, table string, condition string) {
	query := "DELETE FROM " + table
	if condition != "" {
		query += " WHERE " + condition
	}

	_, err := dth.Pool.Exec(ctx, query)
	if err != nil {
		dth.T.Logf("Failed to cleanup test data: %v", err)
	}
}

// InsertTestUser inserts a test user into database
func (dth *DatabaseTestHelper) InsertTestUser(ctx context.Context, userID, email, passwordHash, fullName, role string) error {
	_, err := dth.Pool.Exec(
		ctx,
		`INSERT INTO users (id, email, password_hash, full_name, role, is_active)
		 VALUES ($1, $2, $3, $4, $5, TRUE)`,
		userID, email, passwordHash, fullName, role,
	)
	return err
}

// GetUser retrieves a user from database
func (dth *DatabaseTestHelper) GetUser(ctx context.Context, email string) (map[string]interface{}, error) {
	var userID, passwordHash, fullName, role string
	var isActive bool

	err := dth.Pool.QueryRow(
		ctx,
		"SELECT id, password_hash, full_name, role, is_active FROM users WHERE email = $1",
		email,
	).Scan(&userID, &passwordHash, &fullName, &role, &isActive)

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":             userID,
		"email":          email,
		"password_hash":  passwordHash,
		"full_name":      fullName,
		"role":           role,
		"is_active":      isActive,
	}, nil
}

// ===== Assertion Helpers =====

// AssertStatusCode asserts that response has expected status code
func (h *TestHelpers) AssertStatusCode(resp *http.Response, expected int) {
	if resp.StatusCode != expected {
		h.T.Errorf("Expected status code %d, got %d", expected, resp.StatusCode)
	}
}

// AssertContentType asserts that response has expected content type
func (h *TestHelpers) AssertContentType(resp *http.Response, expected string) {
	contentType := resp.Header.Get("Content-Type")
	if contentType != expected {
		h.T.Errorf("Expected content type %s, got %s", expected, contentType)
	}
}

// AssertJSONField asserts that JSON response contains expected field with value
func (h *TestHelpers) AssertJSONField(data map[string]interface{}, field string, expected interface{}) {
	value, exists := data[field]
	if !exists {
		h.T.Errorf("Field %s not found in response", field)
		return
	}
	if value != expected {
		h.T.Errorf("Field %s: expected %v, got %v", field, expected, value)
	}
}

// ===== Time Helpers =====

// CreateMockTime creates a mock time for testing
func CreateMockTime() time.Time {
	return time.Date(2024, 1, 15, 8, 30, 0, 0, time.UTC)
}

// CreateMockDate creates a mock date string (YYYY-MM-DD)
func CreateMockDate() string {
	return "2024-01-15"
}

// ===== Context Helpers =====

// CreateTestContext creates a test context with timeout
func CreateTestContext(timeout time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	return ctx
}

// ===== Mock Dependencies =====

// MockLogger creates a mock logger for testing
func MockLogger() observability.Logger {
	return observability.InitLogger()
}

// MockRedisClient creates a mock Redis client (or nil for testing)
// In real tests, use testcontainers or embedded Redis
func MockRedisClient() *redis.RedisClient {
	return nil // Implement with actual mock if needed
}

// MockKafkaProducer creates a mock Kafka producer (or nil for testing)
// In real tests, use testcontainers or embedded Kafka
func MockKafkaProducer() *kafka.Producer {
	return nil // Implement with actual mock if needed
}

// ===== Retry Helpers =====

// RetryWithBackoff retries a function with exponential backoff
func RetryWithBackoff(t *testing.T, fn func() error, maxRetries int) error {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		if err := fn(); err == nil {
			return nil
		} else {
			lastErr = err
		}
		time.Sleep(time.Duration(1<<uint(i)) * 100 * time.Millisecond)
	}
	return lastErr
}

// ===== Table-Driven Test Helpers =====

// TestCase represents a single test case
type TestCase struct {
	Name      string
	Input     interface{}
	Expected  interface{}
	Error     bool
	ErrorType string
}

// RunTableDrivenTests runs a set of table-driven tests
func RunTableDrivenTests(t *testing.T, testCases []TestCase, testFn func(tc TestCase) error) {
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := testFn(tc)
			if tc.Error && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tc.Error && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})
	}
}
