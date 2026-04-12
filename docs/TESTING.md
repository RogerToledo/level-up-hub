# Testing Guide

## Overview

This project uses a comprehensive testing strategy with unit tests, integration tests, and automated CI/CD pipelines.

## Test Structure

```
internal/
├── account/
│   ├── service.go
│   └── service_test.go          # Unit tests
├── pagination/
│   ├── pagination.go
│   └── pagination_test.go       # Unit tests
├── mocks/
│   └── repository_mock.go       # Mock implementations
└── testutil/
    └── db.go                     # Test helpers
```

## Running Tests

### All Tests

```bash
# Run all tests
make test

# Run with verbose output
make test-verbose

# Run with coverage report
make cover
```

### Specific Package

```bash
# Test single package
go test ./internal/pagination -v

# Test with coverage
go test ./internal/account -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Integration Tests

```bash
# Requires PostgreSQL running
docker-compose up -d postgres

# Run integration tests
go test -tags=integration ./... -v
```

## Writing Tests

### Unit Tests

Unit tests should test individual functions in isolation using mocks.

**Example:**

```go
package account

import (
    "context"
    "testing"
    
    "github.com/me/level-up-hub/internal/mocks"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestService_FindUserByID(t *testing.T) {
    t.Run("successfully finds user", func(t *testing.T) {
        // Arrange
        mockRepo := new(mocks.MockQuerier)
        service := NewService(mockRepo)
        ctx := context.Background()
        
        userID := uuid.New()
        expectedUser := repository.FindUserByIDRow{
            ID:       userID,
            Username: "testuser",
            Email:    "test@example.com",
        }

        mockRepo.On("FindUserByID", ctx, userID).Return(expectedUser, nil)

        // Act
        user, err := service.FindUserByID(ctx, userID)

        // Assert
        assert.NoError(t, err)
        assert.Equal(t, expectedUser.Username, user.Username)
        mockRepo.AssertExpectations(t)
    })
}
```

### Table-Driven Tests

Use table-driven tests for testing multiple scenarios:

```go
func TestGetPaginationParams(t *testing.T) {
    tests := []struct {
        name           string
        queryParams    map[string]string
        expectedPage   int
        expectedSize   int
    }{
        {
            name:         "default values",
            queryParams:  map[string]string{},
            expectedPage: 1,
            expectedSize: 20,
        },
        {
            name: "custom values",
            queryParams: map[string]string{
                "page":      "2",
                "page_size": "10",
            },
            expectedPage: 2,
            expectedSize: 10,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Integration Tests

Integration tests interact with real database:

```go
// +build integration

package account_test

import (
    "testing"
    "github.com/me/level-up-hub/internal/testutil"
)

func TestIntegration_CreateUser(t *testing.T) {
    // Setup test database
    pool := testutil.SetupTestDB(t)
    defer testutil.CleanupTestData(t, pool, "users")
    
    // Test with real database
    repo := repository.New(pool)
    service := account.NewService(repo)
    
    // ... test implementation
}
```

### HTTP Handler Tests

Test HTTP handlers using `httptest`:

```go
func TestHandler_Login(t *testing.T) {
    gin.SetMode(gin.TestMode)
    
    mockService := new(mocks.MockService)
    handler := NewHandler(mockService, &config.Config{})
    
    w := httptest.NewRecorder()
    c, r := gin.CreateTestContext(w)
    
    r.POST("/login", handler.Login)
    
    body := `{"email":"test@example.com","password":"pass123"}`
    req := httptest.NewRequest("POST", "/login", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    
    r.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
}
```

## Mocking

### Creating Mocks

Use `testify/mock` for creating mocks:

```go
type MockQuerier struct {
    mock.Mock
}

func (m *MockQuerier) FindUserByID(ctx context.Context, id uuid.UUID) (User, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(User), args.Error(1)
}
```

### Using Mocks

```go
mockRepo := new(mocks.MockQuerier)

// Setup expectations
mockRepo.On("FindUserByID", ctx, userID).Return(user, nil)

// Call code under test
result, err := service.GetUser(ctx, userID)

// Verify expectations
mockRepo.AssertExpectations(t)
```

## Test Coverage

### View Coverage

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View in terminal
go tool cover -func=coverage.out

# View in browser
go tool cover -html=coverage.out
```

### Coverage Goals

- **Unit tests:** > 80% coverage
- **Critical business logic:** > 90% coverage
- **Handler/API layer:** > 70% coverage

### Coverage by Package

```bash
# Check coverage per package
go test -cover ./...

# Output:
# ok      github.com/me/level-up-hub/internal/account     coverage: 85.2%
# ok      github.com/me/level-up-hub/internal/pagination  coverage: 92.1%
```

## CI/CD Pipeline

### GitHub Actions

Our CI pipeline runs automatically on:
- Push to `main` or `develop`
- Pull requests to `main` or `develop`

**Pipeline includes:**

1. **Test Job**
   - Spins up PostgreSQL service
   - Runs database migrations
   - Executes all tests
   - Uploads coverage to Codecov

2. **Lint Job**
   - Runs `golangci-lint`
   - Checks code style and quality

3. **Build Job**
   - Compiles the binary
   - Uploads artifact

4. **Security Job**
   - Runs `gosec` security scanner

### Local CI Simulation

```bash
# Run all CI checks locally
make check

# This runs:
# - go fmt
# - golangci-lint
# - go test
```

## Test Best Practices

### ✅ DO

- Write tests first (TDD)
- Test one thing per test
- Use descriptive test names
- Use table-driven tests for multiple scenarios
- Mock external dependencies
- Clean up test data
- Use `t.Helper()` in test utilities
- Test error cases
- Test edge cases and boundaries

### ❌ DON'T

- Don't test implementation details
- Don't depend on test execution order
- Don't use sleep for synchronization
- Don't leave test data in database
- Don't skip tests without good reason
- Don't ignore race conditions
- Don't test third-party libraries

## Benchmarking

### Writing Benchmarks

```go
func BenchmarkPaginationParams(b *testing.B) {
    gin.SetMode(gin.TestMode)
    
    for i := 0; i < b.N; i++ {
        w := httptest.NewRecorder()
        c, _ := gin.CreateTestContext(w)
        req := httptest.NewRequest("GET", "/test?page=1&page_size=20", nil)
        c.Request = req
        
        _ = GetPaginationParams(c)
    }
}
```

### Running Benchmarks

```bash
# Run all benchmarks
go test -bench=. ./...

# Run specific benchmark
go test -bench=BenchmarkPaginationParams ./internal/pagination

# With memory statistics
go test -bench=. -benchmem ./...
```

## Test Utilities

### Database Helper

```go
// Setup test database
pool := testutil.SetupTestDB(t)

// Cleanup after test
defer testutil.CleanupTestData(t, pool, "users", "activities")
```

### Context Helper

```go
// Create test context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

## Troubleshooting

### Tests Hanging

**Problem:** Tests don't complete

**Solution:**
```go
// Add timeout to tests
func TestSomething(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // Your test code
}
```

### Database Connection Issues

**Problem:** Cannot connect to test database

**Solution:**
```bash
# Start PostgreSQL container
docker-compose up -d postgres

# Check if running
docker ps | grep postgres

# Check environment variables
echo $DB_URL_DEV
```

### Race Conditions

**Problem:** `go test -race` fails

**Solution:**
```go
// Use proper synchronization
var mu sync.Mutex
mu.Lock()
defer mu.Unlock()

// Or use channels
done := make(chan bool)
go func() {
    // work
    done <- true
}()
<-done
```

### Mock Not Called

**Problem:** `mock: expected call but got none`

**Solution:**
```go
// Ensure mock is called with exact arguments
mockRepo.On("FindUserByID", ctx, userID).Return(user, nil)

// Use mock.Anything for flexible matching
mockRepo.On("FindUserByID", mock.Anything, mock.Anything).Return(user, nil)

// Verify expectations at end
defer mockRepo.AssertExpectations(t)
```

## Resources

- [Testing in Go](https://go.dev/doc/tutorial/add-a-test)
- [testify documentation](https://github.com/stretchr/testify)
- [Table-Driven Tests](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
- [Go Test Coverage](https://go.dev/blog/cover)
- [golangci-lint](https://golangci-lint.run/)

## Example Test Output

```bash
$ make test

🧪 Executando testes...
=== RUN   TestGetPaginationParams
=== RUN   TestGetPaginationParams/default_values_when_no_params_provided
=== RUN   TestGetPaginationParams/valid_page_and_size
=== RUN   TestGetPaginationParams/negative_page_defaults_to_1
--- PASS: TestGetPaginationParams (0.00s)
    --- PASS: TestGetPaginationParams/default_values_when_no_params_provided (0.00s)
    --- PASS: TestGetPaginationParams/valid_page_and_size (0.00s)
    --- PASS: TestGetPaginationParams/negative_page_defaults_to_1 (0.00s)
PASS
coverage: 92.1% of statements
ok      github.com/me/level-up-hub/internal/pagination  0.234s
```
