# SISU Integration Testing Guide

This document describes the integration testing infrastructure and how to write tests for the SISU application.

## Quick Start

```bash
# Run all tests
go test ./...

# Run only integration tests
go test ./commands

# Run specific test
go test ./commands -run TestCompleteAdmissionCycle

# Run with verbose output
go test -v ./commands

# Run with coverage
go test ./commands -cover
```

## Test Infrastructure

### Test Database (`database/testing.go`)

The `NewTestDatabase()` function creates an isolated test database for each test:

```go
func TestMyFeature(t *testing.T) {
    db := database.NewTestDatabase(t)
    // Database is automatically cleaned up when test finishes

    // Use db.Database to access database methods
    registrations, err := db.FetchRegistrations()
    require.NoError(t, err)
}
```

**Features:**
- Creates temporary SQLite database
- Runs all migrations automatically
- Automatic cleanup via `t.Cleanup()`
- Each test gets its own isolated database
- Foreign keys enabled
- WAL mode for better concurrency

### Test Utilities (`testutil/`)

#### Fixtures (`testutil/fixtures.go`)

Create test domain objects with sensible defaults:

```go
// Create custom selection
selection := testutil.NewTestSelection(func(s *types.Selection) {
    s.Year = 2026
    s.Kind = types.SelectionKindWaitlist
})

// Create custom candidate
candidate := testutil.NewTestCandidate(func(c *types.Candidate) {
    c.CPF = "12345678901"
    c.Name = "João Silva"
})
```

Available builders:
- `NewTestSelection()`
- `NewTestCandidate()`
- `NewTestCourse()`
- `NewTestCall()`
- `NewTestRegistration()`

#### Helpers (`testutil/helpers.go`)

Common test operations:

```go
// Load test data
testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")
testutil.LoadWaitlistSelection(t, db.Database, "testdata/waitlist_small.csv")

// Create and manage calls
// (resolves Semester 1 of year 2025 internally via FetchSemesterByYearAndNumber)
callID := testutil.CreateCall(t, db.Database)
testutil.CloseCall(t, db.Database, callID)
testutil.CloseCallWithEnrollment(t, db.Database, callID) // Enrolls all first

// Update registration statuses
testutil.EnrollRegistration(t, db.Database, regID)
testutil.MarkRegistrationAbsent(t, db.Database, regID)
testutil.ClearRegistrationStatus(t, db.Database, regID)

// Delete selections
testutil.DeleteApprovedSelection(t, db.Database)
testutil.DeleteWaitlistSelection(t, db.Database)
```

#### Assertions (`testutil/assertions.go`)

Domain-specific assertions:

```go
// Assert specific states
testutil.AssertRegistrationStatus(t, db.Database, regID, types.RegistrationStatusEnrolled)
testutil.AssertCallStatus(t, db.Database, callID, types.CallStatusDone)
testutil.AssertCourseOccupiedSeats(t, db.Database, courseID, 10)

// Assert existence
selection := testutil.AssertSelectionExists(t, db.Database, types.SelectionKindApproved)
callID := testutil.AssertCallCreated(t, db.Database, 2) // Returns call ID

// Assert counts
testutil.AssertRegistrationCount(t, db.Database, 5)
testutil.AssertRegistrationsInCall(t, db.Database, callID, 3)

// Assert empty state
testutil.AssertDatabaseEmpty(t, db.Database)
```

### Test Data (`commands/testdata/`)

Pre-created CSV files for testing:

- **approved_small.csv** - 5 approved students across 2 courses (morning: 10 seats, evening: 15 seats)
- **waitlist_small.csv** - 3 waitlisted students
- **invalid_missing_fields.csv** - Malformed CSV for error testing

## Writing Integration Tests

### Basic Test Structure

Follow the Arrange-Act-Assert pattern:

```go
func TestMyFeature(t *testing.T) {
    // Arrange
    db := database.NewTestDatabase(t)
    testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")

    cmd := commands.MyCommand{
        Param: "value",
    }

    // Act
    err := cmd.Execute(db.Database)

    // Assert
    require.NoError(t, err)
    testutil.AssertSomeState(t, db.Database)
}
```

### Testing Error Cases

Test business rule violations:

```go
func TestMyFeature_ErrBusinessRule(t *testing.T) {
    db := database.NewTestDatabase(t)

    // Setup that will cause error
    testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")

    // Try to load approved again (should fail)
    cmd := commands.LoadSelectionCommand{
        Kind: types.SelectionKindApproved,
        // ...
    }

    err := cmd.Execute(db.Database)

    require.Error(t, err)
    assert.ErrorAs(t, err, &commands.ErrApprovedSelectionAlreadyExists{})
}
```

### Table-Driven Tests

Test multiple scenarios:

```go
func TestMyFeature_MultipleScenarios(t *testing.T) {
    tests := []struct {
        name          string
        setup         func(t *testing.T, db *database.Database)
        expectedError any
    }{
        {
            name: "scenario 1",
            setup: func(t *testing.T, db *database.Database) {
                // Setup code
            },
            expectedError: nil, // Success case
        },
        {
            name: "scenario 2 - error",
            setup: func(t *testing.T, db *database.Database) {
                // Setup that causes error
            },
            expectedError: &commands.ErrSomeError{},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            db := database.NewTestDatabase(t)
            tt.setup(t, db.Database)

            // Test execution
            err := doSomething()

            if tt.expectedError != nil {
                require.Error(t, err)
                assert.ErrorAs(t, err, tt.expectedError)
            } else {
                require.NoError(t, err)
            }
        })
    }
}
```

### Testing Complete Workflows

Test end-to-end scenarios:

```go
func TestCompleteWorkflow(t *testing.T) {
    db := database.NewTestDatabase(t)

    // Step 1: Load data
    t.Log("Loading approved selection")
    testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")

    // Step 2: Process call
    t.Log("Processing call #1")
    call1, _ := db.FetchCallByNumber(1)
    testutil.EnrollAllInCall(t, db.Database, call1.ID)
    testutil.CloseCall(t, db.Database, call1.ID)

    // Step 3: Create next call
    t.Log("Creating call #2")
    call2 := testutil.CreateCall(t, db.Database)

    // Assert final state
    testutil.AssertCallStatus(t, db.Database, call2, types.CallStatusCalling)
}
```

## Best Practices

### 1. Test Isolation
- Each test gets its own database
- Use `t.Parallel()` for independent tests
- Don't share state between tests

### 2. Clear Test Names
- Use descriptive names: `TestFeature_Scenario` or `TestFeature_ErrCondition`
- Group related tests with common prefix
- Use subtests for variations

### 3. Meaningful Assertions
- Use domain-specific assertions from `testutil`
- Add helpful messages: `assert.Equal(t, expected, actual, "what should be true")`
- Assert the most important conditions first

### 4. Test Data Management
- Prefer small, focused test data files
- Use existing test CSV files when possible
- Create new files only for specific scenarios

### 5. Error Testing
- Test all custom error types
- Verify error messages are helpful
- Test both error and success paths

### 6. Documentation
- Add comments for complex test setup
- Use `t.Log()` for workflow steps
- Document why tests are skipped

## Common Patterns

### Testing Transactions
Verify rollback on error:

```go
func TestTransactionRollback(t *testing.T) {
    db := database.NewTestDatabase(t)

    // Attempt operation that should fail
    cmd := commands.SomeCommand{/* invalid params */}
    err := cmd.Execute(db.Database)

    require.Error(t, err)

    // Verify nothing was persisted
    testutil.AssertDatabaseEmpty(t, db.Database)
}
```

### Testing State Machines
Verify all transitions:

```go
func TestStatusTransitions(t *testing.T) {
    db := database.NewTestDatabase(t)
    testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")

    regs, _ := db.FetchRegistrations()
    reg := regs[0]

    // Test valid transition
    testutil.EnrollRegistration(t, db.Database, reg.ID)
    testutil.AssertRegistrationStatus(t, db.Database, reg.ID, types.RegistrationStatusEnrolled)

    // Test reverse transition
    testutil.ClearRegistrationStatus(t, db.Database, reg.ID)
    testutil.AssertRegistrationStatus(t, db.Database, reg.ID, types.RegistrationStatusApproved)
}
```

### Testing Business Rules
One test per rule:

```go
func TestBusinessRule_CannotDoXWhenY(t *testing.T) {
    db := database.NewTestDatabase(t)

    // Setup condition Y
    setupConditionY(t, db)

    // Try to do X
    err := doX(db)

    // Verify error
    require.Error(t, err)
    assert.ErrorAs(t, err, &commands.ErrCannotDoXWhenY{})
}
```

## Troubleshooting

### Test Failures
1. Read the error message carefully
2. Check test data files exist
3. Verify database migrations ran
4. Use `t.Log()` to debug workflow
5. Run test with `-v` flag for verbose output

### Database Issues
- Ensure migrations are up to date
- Check foreign key constraints
- Verify test data is valid CSV

### Common Errors
- **"cannot close call with pending registrations"** - Enroll or mark absent all students first
- **"transaction rolled back"** - Check CSV data format
- **"call not found"** - Verify call was created

## Running Tests in CI/CD

```bash
# Run all tests with coverage
go test ./... -coverprofile=coverage.out

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# Run tests with race detector
go test -race ./...

# Run tests with timeout
go test -timeout 30s ./...
```

## Test Coverage Goals

- **Commands**: 80%+ (business logic)
- **Database**: 70%+ (CRUD operations)
- **Overall**: 60%+ (entire codebase)

Coverage varies by package and changes as tests are added — run `go test ./... -cover` for current numbers rather than relying on a figure recorded here.
