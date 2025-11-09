# Testing Guide

This document describes the testing strategy and implementation for the Mini Attendance System.

## Table of Contents

1. [Testing Overview](#testing-overview)
2. [Unit Tests](#unit-tests)
3. [Integration Tests](#integration-tests)
4. [E2E Tests](#e2e-tests)
5. [Running Tests](#running-tests)
6. [Test Coverage](#test-coverage)
7. [Best Practices](#best-practices)
8. [CI/CD Integration](#cicd-integration)

---

## Testing Overview

The project uses a three-layer testing approach:

```
┌─────────────────────────────────────────┐
│         E2E Tests (API Level)           │
│  - Full integration with all services   │
│  - Tests against running server         │
│  - Verifies complete workflows          │
└─────────────────────────────────────────┘
               ↓
┌─────────────────────────────────────────┐
│    Integration Tests (Service Layer)    │
│  - Tests with real database/redis       │
│  - Service-to-service integration       │
│  - Mocked external services             │
└─────────────────────────────────────────┘
               ↓
┌─────────────────────────────────────────┐
│   Unit Tests (Function/Method Level)    │
│  - Isolated component testing           │
│  - All dependencies mocked              │
│  - Fast execution                       │
└─────────────────────────────────────────┘
```

### Testing Framework

- **Testing Framework**: Go's built-in `testing` package
- **Assertion Library**: `github.com/stretchr/testify`
- **Test Structure**: Table-driven tests with test suites
- **Mocking**: `github.com/golang/mock` for dependency mocking

---

## Unit Tests

### Purpose

Unit tests verify individual functions and methods in isolation.

### Location

- `internal/api/http/handlers/*_test.go` - Handler layer tests
- `internal/services/auth/auth_service_test.go` - Service layer tests
- `internal/infra/db/*_test.go` - Database layer tests

### Example: Auth Handler Unit Tests

**File**: `internal/api/http/handlers/auth_handler_test.go`

Tests include:

- **Login Tests**

  - `TestLogin_Success` - Valid credentials
  - `TestLogin_InvalidCredentials` - Wrong password
  - `TestLogin_MissingEmail` - Required field validation
  - `TestLogin_InvalidEmailFormat` - Format validation
  - `TestLogin_UserNotFound` - Non-existent user

- **Register Tests**

  - `TestRegister_Success` - Valid registration
  - `TestRegister_EmailAlreadyExists` - Duplicate email
  - `TestRegister_ShortPassword` - Validation error
  - `TestRegister_MissingFullName` - Required field

- **Token Tests**
  - `TestRefreshToken_Success` - Valid refresh
  - `TestRefreshToken_InvalidToken` - Invalid token
  - `TestRefreshToken_MissingToken` - Missing header

### Example: Auth Service Unit Tests

**File**: `internal/services/auth/auth_service_test.go`

Tests include:

- **Token Generation**

  - `TestGenerateAccessToken_Success`
  - `TestGenerateAccessToken_WithAdminRole`
  - `TestGenerateRefreshToken_Success`
  - `TestGenerateAccessToken_EmptyUserID` (error case)

- **Token Validation**

  - `TestValidateAccessToken_Success`
  - `TestValidateAccessToken_InvalidToken`
  - `TestValidateRefreshToken_Success`

- **Token Lifecycle**
  - `TestCompleteAuthenticationLifecycle`
  - `TestTokenRenewalFlow`

### Running Unit Tests

```bash
# Run all unit tests
go test ./...

# Run specific test package
go test ./internal/api/http/handlers/...

# Run specific test
go test ./internal/services/auth -run TestGenerateAccessToken_Success

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...
```

---

## Integration Tests

### Purpose

Integration tests verify that multiple components work together correctly.

### Characteristics

- Use real database (test database)
- Test service-to-service interactions
- Mock external APIs (Kafka, Redis if needed)
- Slower than unit tests but faster than E2E

### Example: Attendance Integration Tests

**File**: `tests/integration/attendance_test.go` (to be created)

Would include:

- Database setup and teardown
- User creation with actual database operations
- Attendance record creation and retrieval
- Validation of business logic

```go
func TestAttendanceFlow_Integration(t *testing.T) {
    // Setup test database
    db := setupTestDatabase(t)
    defer teardownTestDatabase(db)

    // Create test user
    userID := "test-user-123"
    insertTestUser(t, db, userID)

    // Check-in
    record, err := checkIn(db, userID, -6.1753, 106.8249)
    assert.NoError(t, err)
    assert.NotNil(t, record)

    // Retrieve and verify
    retrieved, _ := getAttendanceRecord(db, userID)
    assert.Equal(t, record.ID, retrieved.ID)
}
```

### Running Integration Tests

```bash
# Run integration tests with test database
docker compose -f docker-compose.test.yml up -d
go test -v -tags=integration ./tests/integration/...
docker compose -f docker-compose.test.yml down
```

---

## E2E Tests

### Purpose

End-to-end tests verify complete user workflows against a running server.

### Location

**File**: `tests/e2e/api_test.go`

### Test Scenarios

#### 1. Authentication Flow

```
Register User → Login → Receive Tokens → Refresh Token
```

Tests:

- `TestAuthenticationFlow` - Complete auth flow
- `TestLoginWithInvalidCredentials` - Error handling
- `TestRefreshToken` - Token refresh

#### 2. Attendance Complete Flow

```
Check-In → View Today's Record → Check-Out → Verify Record
```

Tests:

- `TestCompleteAttendanceFlow` - Full attendance lifecycle
- `TestAttendanceCheckInTwice` - Duplicate check-in prevention
- `TestAttendanceCheckOutWithoutCheckIn` - Error handling

#### 3. Authorization Tests

```
Employee → Can only see own data
Manager → Can see team data
Admin → Can see all data
```

Tests:

- `TestEmployeeCantViewOtherAttendance` - Access control
- `TestManagerCanViewEmployeeAttendance` - Manager permissions
- `TestAdminCanViewAnyAttendance` - Admin permissions

#### 4. Report Access

Tests:

- `TestGetMonthlyReport` - Admin can access
- `TestEmployeeCantAccessMonthlyReport` - Access control

#### 5. Error Handling

Tests:

- `TestInvalidJSONRequest` - Malformed request
- `TestMissingRequiredFields` - Validation
- `TestConcurrentRequests` - Concurrency handling

#### 6. Data Consistency

Tests:

- `TestAttendanceDataConsistency` - Data integrity

### Running E2E Tests

```bash
# Start application server
cd backend
docker compose up -d

# Run E2E tests
go test -v -tags=e2e ./tests/e2e/...

# Stop server
docker compose down
```

### E2E Test Structure

```go
type E2ETestSuite struct {
    suite.Suite
    baseURL    string
    httpClient *http.Client
    tokens     map[string]string
}

func (suite *E2ETestSuite) SetupSuite() {
    // Start server, wait for readiness
    // Register test users
}

func (suite *E2ETestSuite) TearDownSuite() {
    // Cleanup test data
}

func (suite *E2ETestSuite) TestScenario() {
    // Actual test implementation
}
```

---

## Running Tests

### All Tests

```bash
# Run all tests
go test ./...

# Run with coverage report
go test -cover ./...

# Run with detailed coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Specific Test Types

```bash
# Unit tests only
go test -short ./...

# Integration tests (requires test database)
go test -tags=integration ./tests/integration/...

# E2E tests (requires running server)
go test -tags=e2e ./tests/e2e/...
```

### Filter Tests

```bash
# Run tests matching pattern
go test -run "TestLogin" ./...

# Run tests NOT matching pattern
go test -run "!TestLogin" ./...

# Run specific suite
go test -run "AuthHandlerTestSuite" ./internal/api/http/handlers/
```

### Test Execution Options

```bash
# Verbose output
go test -v ./...

# Show function execution time
go test -v -count=1 ./...

# Run with timeout (30 seconds)
go test -timeout 30s ./...

# Run in parallel (use CPU cores)
go test -parallel 4 ./...

# Run once (disable test caching)
go test -count=1 ./...
```

---

## Test Coverage

### Measuring Coverage

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage report
go tool cover -html=coverage.out

# Coverage for specific package
go test -cover ./internal/api/http/handlers/

# Coverage threshold check
go test -coverprofile=coverage.out ./... && \
  go tool cover -func=coverage.out | grep total | awk '{print $3}'
```

### Coverage Goals

- **Unit Tests**: Aim for 80%+ coverage
- **Critical Paths**: Aim for 100% coverage (auth, attendance)
- **Handlers**: Aim for 75%+ coverage
- **Services**: Aim for 85%+ coverage
- **Overall**: Aim for 70%+ project-wide

### Coverage Report Example

```
total:  (statements) 75.4%
```

---

## Best Practices

### 1. Test Naming

```go
// ✓ Good: Descriptive, states what and expected outcome
func TestLogin_Success() {}
func TestLogin_InvalidCredentials() {}
func TestCheckIn_AlreadyCheckedInToday() {}

// ✗ Bad: Vague or unclear
func TestLogin() {}
func TestFail() {}
func Test1() {}
```

### 2. Table-Driven Tests

```go
// ✓ Good: Organized, maintainable, easy to add cases
testCases := []struct {
    name      string
    input     string
    expected  bool
    wantError bool
}{
    {"valid email", "test@example.com", true, false},
    {"invalid email", "invalid", false, true},
}

for _, tc := range testCases {
    t.Run(tc.name, func(t *testing.T) {
        result, err := validate(tc.input)
        assert.Equal(t, tc.expected, result)
        if tc.wantError {
            assert.Error(t, err)
        }
    })
}
```

### 3. Test Isolation

```go
// ✓ Good: Each test is independent
func TestFeatureA(t *testing.T) {
    setup := createCleanEnvironment()
    defer cleanup(setup)
    // Test feature A
}

// ✗ Bad: Tests depend on execution order
var globalState = 0
func TestFeature1(t *testing.T) {
    globalState = 1
}
func TestFeature2(t *testing.T) {
    assert.Equal(t, globalState, 1) // Depends on TestFeature1
}
```

### 4. Use Builders for Test Data

```go
// ✓ Good: Clear, readable, reusable
user := NewUserBuilder().
    WithEmail("test@example.com").
    WithRole("manager").
    Build()

// ✗ Bad: Hard to read, error-prone
user := map[string]interface{}{
    "email": "test@example.com",
    "password": "password123",
    "full_name": "Test",
    "role": "manager",
}
```

### 5. Assert Specific Errors

```go
// ✓ Good: Specific error checking
_, err := login("invalid@example.com", "password")
assert.Error(t, err)
assert.Contains(t, err.Error(), "user not found")

// ✗ Bad: Generic error checking
_, err := login("invalid@example.com", "password")
assert.Error(t, err)
```

### 6. Clean Up Resources

```go
// ✓ Good: Proper cleanup
func TestWithDB(t *testing.T) {
    db := setupTestDB()
    defer db.Close()

    defer func() {
        if err := recover(); err != nil {
            db.Close()
            panic(err)
        }
    }()

    // Test code
}
```

### 7. Use Subtests

```go
// ✓ Good: Clear test organization
func TestAuth(t *testing.T) {
    t.Run("Login", func(t *testing.T) {
        t.Run("Success", func(t *testing.T) { /* ... */ })
        t.Run("InvalidCredentials", func(t *testing.T) { /* ... */ })
    })
    t.Run("Register", func(t *testing.T) { /* ... */ })
}
```

---

## CI/CD Integration

### GitHub Actions Configuration

**File**: `.github/workflows/test.yml`

```yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest

    services:
      postgres:
        image: postgres:16-alpine
        env:
          POSTGRES_PASSWORD: postgres
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
      - uses: actions/checkout@v3

      - uses: actions/setup-go@v4
        with:
          go-version: '1.24.0'

      - name: Run unit tests
        run: go test -short -race -coverprofile=coverage.out ./...

      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage.out

      - name: Run E2E tests
        run: go test -tags=e2e ./tests/e2e/...
```

### Pre-commit Hooks

```bash
#!/bin/bash
# .githooks/pre-commit

go test -short ./... || exit 1
go fmt ./...
go vet ./...
```

---

## Test Data Management

### Using Seeds

```go
// Setup test data before tests
func setupTestData(t *testing.T, db *pgxpool.Pool) {
    ctx := context.Background()

    // Create test users
    users := []struct {
        id   string
        email string
        role  string
    }{
        {"user1", "admin@test.com", "admin"},
        {"user2", "manager@test.com", "manager"},
    }

    for _, user := range users {
        insertTestUser(t, ctx, db, user)
    }
}
```

### Cleanup

```go
func teardownTestData(t *testing.T, db *pgxpool.Pool) {
    ctx := context.Background()

    // Delete in reverse dependency order
    db.Exec(ctx, "DELETE FROM attendance_records")
    db.Exec(ctx, "DELETE FROM users")
}
```

---

## Troubleshooting

### Common Issues

**Issue**: Tests fail with "connection refused"

```bash
# Solution: Ensure database is running
docker compose up -d
```

**Issue**: Tests timeout

```bash
# Solution: Increase timeout
go test -timeout 60s ./...
```

**Issue**: Flaky tests (pass sometimes, fail sometimes)

```bash
# Solution: Check for test isolation issues, race conditions
go test -race ./...
```

**Issue**: Mock not working as expected

```bash
# Solution: Verify mock setup in SetupTest()
func (suite *AuthHandlerTestSuite) SetupTest() {
    gin.SetMode(gin.TestMode)
    suite.router = gin.New()
    suite.handlers = &handler.Handlers{}
}
```

---

## Testing Checklist

- [ ] Unit tests for all public functions
- [ ] Integration tests for database operations
- [ ] E2E tests for critical user workflows
- [ ] Tests for error cases and edge cases
- [ ] Tests for authorization and access control
- [ ] Tests for concurrent operations
- [ ] Coverage report generated and reviewed
- [ ] All tests passing locally
- [ ] Tests passing in CI/CD pipeline
- [ ] Performance regression tests (if applicable)

---

## Further Reading

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Table-Driven Tests](https://dave.cheney.net/2013/06/09/writing-table-driven-tests-in-go)
- [Go Best Practices](https://golang.org/doc/effective_go)
