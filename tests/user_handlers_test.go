package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"felix1234567890/go-trello/handlers"
	"felix1234567890/go-trello/models"
	"felix1234567890/go-trello/utils"
	"net/http"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockUserService is a mock implementation of the UserService interface
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(user *models.User) (uint, error) {
	args := m.Called(user)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockUserService) GetUsers() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserService) GetUserById(id string) (models.User, error) {
	args := m.Called(id)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(id string, user *models.UpdateUserRequest) error {
	args := m.Called(id, user)
	return args.Error(0)
}

func (m *MockUserService) DeleteUser(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserService) LoginUser(user *models.LoginUserRequest) (uint, error) {
	args := m.Called(user)
	return args.Get(0).(uint), args.Error(1)
}

func setupTestApp() *fiber.App {
	// Initialize secret key for JWT
	os.Setenv("JWT_SECRET", "test-secret-key")
	os.Setenv("SECRET_KEY", "test-secret-key")
	utils.InitSecretKey()
	return fiber.New()
}

func TestUserHandler_CreateUser(t *testing.T) {
	app := setupTestApp()
	mockService := new(MockUserService)
	handler := handlers.NewUserHandler(mockService)

	app.Post("/users", handler.CreateUser)

	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func()
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "successful user creation",
			requestBody: models.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func() {
				mockService.On("CreateUser", mock.AnythingOfType("*models.User")).
					Return(uint(1), nil)
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "token")
			},
		},
		{
			name: "invalid request body",
			requestBody: "invalid-json",
			setupMock: func() {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "Invalid request body", response["message"])
			},
		},
		{
			name: "validation errors",
			requestBody: models.CreateUserRequest{
				Username: "abc",
				Email:    "invalid-email",
				Password: "123",
			},
			setupMock: func() {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "errors")
			},
		},
		{
			name: "service error",
			requestBody: models.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func() {
				mockService.On("CreateUser", mock.AnythingOfType("*models.User")).
					Return(uint(0), errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "Failed to create user", response["message"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.ExpectedCalls = nil
			tt.setupMock()

			var body []byte
			if tt.requestBody != nil {
				body, _ = json.Marshal(tt.requestBody)
			}

			req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			responseBody := make([]byte, resp.ContentLength)
			resp.Body.Read(responseBody)
			tt.checkResponse(t, responseBody)

			mockService.AssertExpectations(t)
		})
	}
}

func TestUserHandler_Login(t *testing.T) {
	app := setupTestApp()
	mockService := new(MockUserService)
	handler := handlers.NewUserHandler(mockService)

	app.Post("/login", handler.Login)

	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func()
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "successful login",
			requestBody: models.LoginUserRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func() {
				mockService.On("LoginUser", mock.AnythingOfType("*models.LoginUserRequest")).
					Return(uint(1), nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "token")
			},
		},
		{
			name: "invalid credentials",
			requestBody: models.LoginUserRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			setupMock: func() {
				mockService.On("LoginUser", mock.AnythingOfType("*models.LoginUserRequest")).
					Return(uint(0), utils.ErrInvalidCredentials)
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "Invalid username or password", response["message"])
			},
		},
		{
			name: "service error",
			requestBody: models.LoginUserRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func() {
				mockService.On("LoginUser", mock.AnythingOfType("*models.LoginUserRequest")).
					Return(uint(0), errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "Login failed", response["message"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.ExpectedCalls = nil
			tt.setupMock()

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			responseBody := make([]byte, resp.ContentLength)
			resp.Body.Read(responseBody)
			tt.checkResponse(t, responseBody)

			mockService.AssertExpectations(t)
		})
	}
}

func TestUserHandler_GetUsers(t *testing.T) {
	app := setupTestApp()
	mockService := new(MockUserService)
	handler := handlers.NewUserHandler(mockService)

	app.Get("/users", handler.GetUsers)

	tests := []struct {
		name           string
		setupMock      func()
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "successful get users",
			setupMock: func() {
				users := []models.User{
					{Username: "user1", Email: "user1@example.com"},
					{Username: "user2", Email: "user2@example.com"},
				}
				mockService.On("GetUsers").Return(users, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "users")
				users := response["users"].([]interface{})
				assert.Len(t, users, 2)
			},
		},
		{
			name: "service error",
			setupMock: func() {
				mockService.On("GetUsers").Return([]models.User{}, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "Failed to retrieve users", response["message"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.ExpectedCalls = nil
			tt.setupMock()

			req, _ := http.NewRequest("GET", "/users", nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			responseBody := make([]byte, resp.ContentLength)
			resp.Body.Read(responseBody)
			tt.checkResponse(t, responseBody)

			mockService.AssertExpectations(t)
		})
	}
}

func TestUserHandler_GetUserById(t *testing.T) {
	app := setupTestApp()
	mockService := new(MockUserService)
	handler := handlers.NewUserHandler(mockService)

	app.Get("/users/:id", handler.GetUserById)

	tests := []struct {
		name           string
		userID         string
		setupMock      func()
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name:   "successful get user by id",
			userID: "1",
			setupMock: func() {
				user := models.User{Username: "testuser", Email: "test@example.com"}
				mockService.On("GetUserById", "1").Return(user, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "user")
				user := response["user"].(map[string]interface{})
				assert.Equal(t, "testuser", user["username"])
			},
		},
		{
			name:   "user not found",
			userID: "999",
			setupMock: func() {
				mockService.On("GetUserById", "999").Return(models.User{}, gorm.ErrRecordNotFound)
			},
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "User with id 999 not found", response["message"])
			},
		},
		{
			name:   "service error",
			userID: "1",
			setupMock: func() {
				mockService.On("GetUserById", "1").Return(models.User{}, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "Failed to retrieve user", response["message"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.ExpectedCalls = nil
			tt.setupMock()

			req, _ := http.NewRequest("GET", "/users/"+tt.userID, nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			responseBody := make([]byte, resp.ContentLength)
			resp.Body.Read(responseBody)
			tt.checkResponse(t, responseBody)

			mockService.AssertExpectations(t)
		})
	}
}

func TestUserHandler_UpdateUser(t *testing.T) {
	app := setupTestApp()
	mockService := new(MockUserService)
	handler := handlers.NewUserHandler(mockService)

	app.Put("/users/:id", handler.UpdateUser)

	tests := []struct {
		name           string
		userID         string
		requestBody    interface{}
		setupMock      func()
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name:   "successful user update",
			userID: "1",
			requestBody: models.UpdateUserRequest{
				Username: "updateduser",
				Email:    "updated@example.com",
			},
			setupMock: func() {
				mockService.On("UpdateUser", "1", mock.AnythingOfType("*models.UpdateUserRequest")).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "User updated successfully", response["message"])
			},
		},
		{
			name:           "invalid request body",
			userID:         "1",
			requestBody:    "invalid-json",
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "Invalid request body", response["message"])
			},
		},
		{
			name:   "user not found",
			userID: "999",
			requestBody: models.UpdateUserRequest{
				Username: "updateduser",
			},
			setupMock: func() {
				mockService.On("UpdateUser", "999", mock.AnythingOfType("*models.UpdateUserRequest")).
					Return(gorm.ErrRecordNotFound)
			},
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "User with id 999 not found for update", response["message"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.ExpectedCalls = nil
			tt.setupMock()

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("PUT", "/users/"+tt.userID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			responseBody := make([]byte, resp.ContentLength)
			resp.Body.Read(responseBody)
			tt.checkResponse(t, responseBody)

			mockService.AssertExpectations(t)
		})
	}
}

func TestUserHandler_DeleteUser(t *testing.T) {
	app := setupTestApp()
	mockService := new(MockUserService)
	handler := handlers.NewUserHandler(mockService)

	app.Delete("/users/:id", handler.DeleteUser)

	tests := []struct {
		name           string
		userID         string
		setupMock      func()
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name:   "successful user deletion",
			userID: "1",
			setupMock: func() {
				mockService.On("DeleteUser", "1").Return(nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "User deleted successfully", response["message"])
			},
		},
		{
			name:   "user not found",
			userID: "999",
			setupMock: func() {
				mockService.On("DeleteUser", "999").Return(gorm.ErrRecordNotFound)
			},
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "User with id 999 not found", response["message"])
			},
		},
		{
			name:   "service error",
			userID: "1",
			setupMock: func() {
				mockService.On("DeleteUser", "1").Return(errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "Failed to delete user", response["message"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.ExpectedCalls = nil
			tt.setupMock()

			req, _ := http.NewRequest("DELETE", "/users/"+tt.userID, nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			responseBody := make([]byte, resp.ContentLength)
			resp.Body.Read(responseBody)
			tt.checkResponse(t, responseBody)

			mockService.AssertExpectations(t)
		})
	}
}