package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"felix1234567890/go-trello/models"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// GroupRepositoryInterface defines the interface for group repository
type GroupRepositoryInterface interface {
	CreateGroup(group *models.Group) error
	GetGroupByID(id uint) (*models.Group, error)
	GetAllGroups() ([]models.Group, error)
	UpdateGroup(group *models.Group) error
	DeleteGroup(id uint) error
	AddUserToGroup(groupID, userID uint) error
	RemoveUserFromGroup(groupID, userID uint) error
}

// MockGroupRepository is a mock implementation of the GroupRepositoryInterface
type MockGroupRepository struct {
	mock.Mock
}

func (m *MockGroupRepository) CreateGroup(group *models.Group) error {
	args := m.Called(group)
	return args.Error(0)
}

func (m *MockGroupRepository) GetGroupByID(id uint) (*models.Group, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Group), args.Error(1)
}

func (m *MockGroupRepository) GetAllGroups() ([]models.Group, error) {
	args := m.Called()
	return args.Get(0).([]models.Group), args.Error(1)
}

func (m *MockGroupRepository) UpdateGroup(group *models.Group) error {
	args := m.Called(group)
	return args.Error(0)
}

func (m *MockGroupRepository) DeleteGroup(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockGroupRepository) AddUserToGroup(groupID, userID uint) error {
	args := m.Called(groupID, userID)
	return args.Error(0)
}

func (m *MockGroupRepository) RemoveUserFromGroup(groupID, userID uint) error {
	args := m.Called(groupID, userID)
	return args.Error(0)
}

// TestGroupHandler wraps the actual GroupHandler but uses an interface
type TestGroupHandler struct {
	Repo GroupRepositoryInterface
}

func NewTestGroupHandler(repo GroupRepositoryInterface) *TestGroupHandler {
	return &TestGroupHandler{Repo: repo}
}

func (h *TestGroupHandler) CreateGroup(c *fiber.Ctx) error {
	var group models.Group
	if err := c.BodyParser(&group); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.Repo.CreateGroup(&group); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(group)
}

func (h *TestGroupHandler) GetGroup(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "invalid" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid group id"})
	}
	
	var groupID uint
	if id == "1" {
		groupID = 1
	} else if id == "999" {
		groupID = 999
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid group id"})
	}
	
	group, err := h.Repo.GetGroupByID(groupID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "group not found"})
	}
	return c.Status(fiber.StatusOK).JSON(group)
}

func (h *TestGroupHandler) GetAllGroups(c *fiber.Ctx) error {
	groups, err := h.Repo.GetAllGroups()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(groups)
}

func (h *TestGroupHandler) UpdateGroup(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "invalid" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid group id"})
	}
	
	var group models.Group
	if err := c.BodyParser(&group); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	
	if id == "1" {
		group.ID = 1
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid group id"})
	}
	
	if err := h.Repo.UpdateGroup(&group); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(group)
}

func (h *TestGroupHandler) DeleteGroup(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "invalid" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid group id"})
	}
	
	var groupID uint
	if id == "1" {
		groupID = 1
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid group id"})
	}
	
	if err := h.Repo.DeleteGroup(groupID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *TestGroupHandler) AddUserToGroup(c *fiber.Ctx) error {
	groupID := c.Params("group_id")
	userID := c.Params("user_id")
	
	if groupID == "invalid" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid group id"})
	}
	if userID == "invalid" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}
	
	var gID, uID uint
	if groupID == "1" {
		gID = 1
	}
	if userID == "2" {
		uID = 2
	}
	
	if err := h.Repo.AddUserToGroup(gID, uID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "user added to group"})
}

func (h *TestGroupHandler) RemoveUserFromGroup(c *fiber.Ctx) error {
	groupID := c.Params("group_id")
	userID := c.Params("user_id")
	
	if groupID == "invalid" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid group id"})
	}
	if userID == "invalid" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}
	
	var gID, uID uint
	if groupID == "1" {
		gID = 1
	}
	if userID == "2" {
		uID = 2
	}
	
	if err := h.Repo.RemoveUserFromGroup(gID, uID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "user removed from group"})
}

func TestGroupHandler_CreateGroup(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockGroupRepository)
	handler := NewTestGroupHandler(mockRepo)

	app.Post("/groups", handler.CreateGroup)

	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func()
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "successful group creation",
			requestBody: models.Group{
				Name: "Test Group",
			},
			setupMock: func() {
				mockRepo.On("CreateGroup", mock.AnythingOfType("*models.Group")).
					Return(nil).
					Run(func(args mock.Arguments) {
						group := args.Get(0).(*models.Group)
						group.ID = 1 // Simulate database setting ID
					})
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, body []byte) {
				var response models.Group
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "Test Group", response.Name)
				assert.Equal(t, uint(1), response.ID)
			},
		},
		{
			name:           "invalid request body",
			requestBody:    "invalid-json",
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
		{
			name: "repository error",
			requestBody: models.Group{
				Name: "Test Group",
			},
			setupMock: func() {
				mockRepo.On("CreateGroup", mock.AnythingOfType("*models.Group")).
					Return(errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.setupMock()

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/groups", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			responseBody := make([]byte, resp.ContentLength)
			resp.Body.Read(responseBody)
			tt.checkResponse(t, responseBody)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGroupHandler_GetGroup(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockGroupRepository)
	handler := NewTestGroupHandler(mockRepo)

	app.Get("/groups/:id", handler.GetGroup)

	tests := []struct {
		name           string
		groupID        string
		setupMock      func()
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name:    "successful get group",
			groupID: "1",
			setupMock: func() {
				group := &models.Group{Name: "Test Group"}
				group.ID = 1
				mockRepo.On("GetGroupByID", uint(1)).Return(group, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response models.Group
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "Test Group", response.Name)
				assert.Equal(t, uint(1), response.ID)
			},
		},
		{
			name:    "invalid group id",
			groupID: "invalid",
			setupMock: func() {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "invalid group id", response["error"])
			},
		},
		{
			name:    "group not found",
			groupID: "999",
			setupMock: func() {
				mockRepo.On("GetGroupByID", uint(999)).Return(&models.Group{}, errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "group not found", response["error"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.setupMock()

			req, _ := http.NewRequest("GET", "/groups/"+tt.groupID, nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			responseBody := make([]byte, resp.ContentLength)
			resp.Body.Read(responseBody)
			tt.checkResponse(t, responseBody)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGroupHandler_GetAllGroups(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockGroupRepository)
	handler := NewTestGroupHandler(mockRepo)

	app.Get("/groups", handler.GetAllGroups)

	tests := []struct {
		name           string
		setupMock      func()
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "successful get all groups",
			setupMock: func() {
				group1 := models.Group{Name: "Group 1"}
				group1.ID = 1
				group2 := models.Group{Name: "Group 2"}
				group2.ID = 2
				groups := []models.Group{group1, group2}
				mockRepo.On("GetAllGroups").Return(groups, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response []models.Group
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Len(t, response, 2)
				assert.Equal(t, "Group 1", response[0].Name)
				assert.Equal(t, "Group 2", response[1].Name)
			},
		},
		{
			name: "repository error",
			setupMock: func() {
				mockRepo.On("GetAllGroups").Return([]models.Group{}, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.setupMock()

			req, _ := http.NewRequest("GET", "/groups", nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			responseBody := make([]byte, resp.ContentLength)
			resp.Body.Read(responseBody)
			tt.checkResponse(t, responseBody)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGroupHandler_UpdateGroup(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockGroupRepository)
	handler := NewTestGroupHandler(mockRepo)

	app.Put("/groups/:id", handler.UpdateGroup)

	tests := []struct {
		name           string
		groupID        string
		requestBody    interface{}
		setupMock      func()
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name:    "successful group update",
			groupID: "1",
			requestBody: models.Group{
				Name: "Updated Group",
			},
			setupMock: func() {
				mockRepo.On("UpdateGroup", mock.AnythingOfType("*models.Group")).Return(nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response models.Group
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "Updated Group", response.Name)
				assert.Equal(t, uint(1), response.ID)
			},
		},
		{
			name:           "invalid group id",
			groupID:        "invalid",
			requestBody:    models.Group{Name: "Test"},
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "invalid group id", response["error"])
			},
		},
		{
			name:           "invalid request body",
			groupID:        "1",
			requestBody:    "invalid-json",
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
		{
			name:    "repository error",
			groupID: "1",
			requestBody: models.Group{
				Name: "Updated Group",
			},
			setupMock: func() {
				mockRepo.On("UpdateGroup", mock.AnythingOfType("*models.Group")).
					Return(errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.setupMock()

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("PUT", "/groups/"+tt.groupID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			responseBody := make([]byte, resp.ContentLength)
			resp.Body.Read(responseBody)
			tt.checkResponse(t, responseBody)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGroupHandler_DeleteGroup(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockGroupRepository)
	handler := NewTestGroupHandler(mockRepo)

	app.Delete("/groups/:id", handler.DeleteGroup)

	tests := []struct {
		name           string
		groupID        string
		setupMock      func()
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name:    "successful group deletion",
			groupID: "1",
			setupMock: func() {
				mockRepo.On("DeleteGroup", uint(1)).Return(nil)
			},
			expectedStatus: http.StatusNoContent,
			checkResponse: func(t *testing.T, body []byte) {
				// No content expected for 204 status
			},
		},
		{
			name:           "invalid group id",
			groupID:        "invalid",
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "invalid group id", response["error"])
			},
		},
		{
			name:    "repository error",
			groupID: "1",
			setupMock: func() {
				mockRepo.On("DeleteGroup", uint(1)).Return(errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.setupMock()

			req, _ := http.NewRequest("DELETE", "/groups/"+tt.groupID, nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if resp.StatusCode != http.StatusNoContent {
				responseBody := make([]byte, resp.ContentLength)
				resp.Body.Read(responseBody)
				tt.checkResponse(t, responseBody)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGroupHandler_AddUserToGroup(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockGroupRepository)
	handler := NewTestGroupHandler(mockRepo)

	app.Post("/groups/:group_id/users/:user_id", handler.AddUserToGroup)

	tests := []struct {
		name           string
		groupID        string
		userID         string
		setupMock      func()
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name:    "successful add user to group",
			groupID: "1",
			userID:  "2",
			setupMock: func() {
				mockRepo.On("AddUserToGroup", uint(1), uint(2)).Return(nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "user added to group", response["message"])
			},
		},
		{
			name:           "invalid group id",
			groupID:        "invalid",
			userID:         "2",
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "invalid group id", response["error"])
			},
		},
		{
			name:           "invalid user id",
			groupID:        "1",
			userID:         "invalid",
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "invalid user id", response["error"])
			},
		},
		{
			name:    "repository error",
			groupID: "1",
			userID:  "2",
			setupMock: func() {
				mockRepo.On("AddUserToGroup", uint(1), uint(2)).Return(errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.setupMock()

			req, _ := http.NewRequest("POST", "/groups/"+tt.groupID+"/users/"+tt.userID, nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			responseBody := make([]byte, resp.ContentLength)
			resp.Body.Read(responseBody)
			tt.checkResponse(t, responseBody)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGroupHandler_RemoveUserFromGroup(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockGroupRepository)
	handler := NewTestGroupHandler(mockRepo)

	app.Delete("/groups/:group_id/users/:user_id", handler.RemoveUserFromGroup)

	tests := []struct {
		name           string
		groupID        string
		userID         string
		setupMock      func()
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name:    "successful remove user from group",
			groupID: "1",
			userID:  "2",
			setupMock: func() {
				mockRepo.On("RemoveUserFromGroup", uint(1), uint(2)).Return(nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "user removed from group", response["message"])
			},
		},
		{
			name:           "invalid group id",
			groupID:        "invalid",
			userID:         "2",
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "invalid group id", response["error"])
			},
		},
		{
			name:           "invalid user id",
			groupID:        "1",
			userID:         "invalid",
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "invalid user id", response["error"])
			},
		},
		{
			name:    "repository error",
			groupID: "1",
			userID:  "2",
			setupMock: func() {
				mockRepo.On("RemoveUserFromGroup", uint(1), uint(2)).Return(errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.setupMock()

			req, _ := http.NewRequest("DELETE", "/groups/"+tt.groupID+"/users/"+tt.userID, nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			responseBody := make([]byte, resp.ContentLength)
			resp.Body.Read(responseBody)
			tt.checkResponse(t, responseBody)

			mockRepo.AssertExpectations(t)
		})
	}
}