# Go-Trello Testing Suite

This directory contains comprehensive tests for the Go-Trello application, covering all handler methods with both unit tests and end-to-end (e2e) tests.

## Test Structure

### Test Files

- `e2e_test.go` - End-to-end integration tests using a real test database
- `user_handlers_test.go` - Unit tests for user handler methods with mocked dependencies  
- `group_handlers_test.go` - Unit tests for group handler methods with mocked dependencies
- `test_helpers.go` - Common test utilities and helper functions

### Handler Methods Covered

#### User Handlers (7 methods):
1. **CreateUser** - POST `/api/users/`
   - Tests successful user creation with JWT token generation
   - Tests validation errors (invalid email, short username/password)
   - Tests service layer errors
   
2. **Login** - POST `/api/users/login`
   - Tests successful authentication with token generation
   - Tests invalid credentials
   - Tests service layer errors

3. **GetUsers** - GET `/api/users/`
   - Tests retrieving all users
   - Tests service layer errors

4. **GetUserById** - GET `/api/users/:id`
   - Tests retrieving specific user by ID
   - Tests user not found (404)
   - Tests service layer errors

5. **UpdateUser** - PUT `/api/users/:id`
   - Tests successful user updates
   - Tests validation errors
   - Tests user not found (404)
   - Tests service layer errors

6. **DeleteUser** - DELETE `/api/users/:id`
   - Tests successful user deletion
   - Tests user not found (404)
   - Tests service layer errors

7. **GetMe** - GET `/api/users/me` (requires authentication)
   - Tests authenticated user profile retrieval
   - Tests with valid JWT tokens

#### Group Handlers (7 methods):
1. **CreateGroup** - POST `/api/groups/`
   - Tests successful group creation
   - Tests invalid request body
   - Tests repository errors

2. **GetGroup** - GET `/api/groups/:id`
   - Tests retrieving specific group by ID
   - Tests invalid group ID format
   - Tests group not found

3. **GetAllGroups** - GET `/api/groups/`
   - Tests retrieving all groups
   - Tests repository errors

4. **UpdateGroup** - PUT `/api/groups/:id`
   - Tests successful group updates
   - Tests invalid group ID format
   - Tests invalid request body
   - Tests repository errors

5. **DeleteGroup** - DELETE `/api/groups/:id`
   - Tests successful group deletion
   - Tests invalid group ID format
   - Tests repository errors

6. **AddUserToGroup** - POST `/api/groups/:group_id/users/:user_id`
   - Tests successfully adding user to group
   - Tests invalid group/user ID formats
   - Tests repository errors

7. **RemoveUserFromGroup** - DELETE `/api/groups/:group_id/users/:user_id`
   - Tests successfully removing user from group
   - Tests invalid group/user ID formats
   - Tests repository errors

## Test Types

### Unit Tests
- Use mocked dependencies (services, repositories)
- Test individual handler functions in isolation
- Fast execution
- Test error conditions and edge cases

### Integration/E2E Tests  
- Use real test database (SQLite in-memory)
- Test complete request/response cycles
- Test database interactions
- Test authentication flows

## Running Tests

### Prerequisites
```bash
# Install dependencies
make deps
```

### All Tests
```bash
# Run all tests with coverage
make test

# Run all tests with detailed coverage report
make test-coverage
```

### Specific Test Types
```bash
# Run only unit tests (fast)
make test-unit

# Run only handler unit tests
make test-handlers

# Run only e2e tests
make test-e2e

# Run all integration tests
make test-integration
```

### Handler-Specific Tests
```bash
# Test user handlers only
make test-user-handlers

# Test group handlers only  
make test-group-handlers
```

### Coverage Reports
```bash
# Generate HTML coverage report
make test-coverage

# Show coverage by function
make test-coverage-func
```

## Test Configuration

### Environment Variables
Tests use these environment variables (automatically set in test environment):
- `JWT_SECRET` - Secret key for JWT token generation
- Database configuration (for e2e tests using SQLite)

### Test Database
- E2E tests use SQLite in-memory database
- Database is automatically migrated and cleaned between tests
- No external database required for testing

## Test Utilities

### TestDatabase Helper
- `SetupTestDB()` - Creates in-memory test database
- `CleanupTestDB()` - Cleans test data between tests
- `SeedTestUser()` - Creates test user data
- `SeedTestGroup()` - Creates test group data
- Various assertion helpers for database state

### Mock Objects
- `MockUserService` - Mocks user service interface
- `MockGroupRepository` - Mocks group repository interface
- Uses testify/mock for behavior verification

## Continuous Integration

### CI Commands
```bash
# Full CI pipeline
make ci-build

# CI test with coverage
make ci-test

# CI linting
make ci-lint
```

## Test Scenarios Covered

### Authentication & Authorization
- JWT token generation and validation
- User login/logout flows
- Protected endpoint access

### Data Validation
- Request body parsing errors
- Validation rule enforcement
- Invalid data format handling

### Error Handling
- Database connection errors
- Record not found errors
- Internal server errors
- Bad request errors

### CRUD Operations
- Create, Read, Update, Delete for users
- Create, Read, Update, Delete for groups
- User-group relationship management

### Edge Cases
- Invalid ID formats
- Non-existent resources
- Malformed JSON requests
- Empty response handling

## Best Practices Demonstrated

1. **Test Isolation** - Each test is independent and doesn't affect others
2. **Mocking** - External dependencies are mocked for unit tests
3. **Coverage** - All handler methods and error paths are tested
4. **Realistic Scenarios** - E2E tests use realistic data flows
5. **Helper Functions** - Common test setup is abstracted into reusable helpers
6. **Clear Assertions** - Tests have clear, descriptive assertions
7. **Error Testing** - Both success and failure paths are covered

## Running Specific Tests

```bash
# Run a specific test function
go test ./tests -run TestUserHandler_CreateUser -v

# Run tests matching a pattern
go test ./tests -run "TestUserHandler.*" -v

# Run with race detection
go test ./tests -race -v

# Run with short flag (excludes slow tests)
go test ./tests -short -v
```

## Debugging Tests

```bash
# Run tests with verbose output
go test ./tests -v

# Run tests with coverage and save profile
go test ./tests -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run single test with detailed output
go test ./tests -run TestUserHandler_CreateUser -v -count=1
```