package tests

import (
	"bytes"
	"encoding/json"
	"felix1234567890/go-trello/handlers"
	"felix1234567890/go-trello/models"
	"felix1234567890/go-trello/repository"
	"felix1234567890/go-trello/routes"
	"felix1234567890/go-trello/utils"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type E2ETestSuite struct {
	suite.Suite
	app        *fiber.App
	db         *gorm.DB
	testUser   *models.User
	testGroup  *models.Group
	authToken  string
}

func (suite *E2ETestSuite) SetupSuite() {
	// Setup test database using SQLite
	var err error
	suite.db, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	suite.Require().NoError(err)

	// Migrate the schema
	err = suite.db.AutoMigrate(&models.User{}, &models.Group{})
	suite.Require().NoError(err)

	// Initialize secret key for JWT
	os.Setenv("JWT_SECRET", "test-secret-key")
	os.Setenv("SECRET_KEY", "test-secret-key")
	utils.InitSecretKey()

	// Setup Fiber app
	suite.app = fiber.New()
	suite.setupRoutes()
}

func (suite *E2ETestSuite) setupRoutes() {
	globalPrefix := suite.app.Group("/api")
	
	// Setup user routes
	userRoutes := globalPrefix.Group("/users")
	routes.SetupUserRoutes(userRoutes, suite.db)

	// Setup group routes
	groupRepo := repository.NewGroupRepository(suite.db)
	groupHandler := handlers.NewGroupHandler(groupRepo)
	routes.SetupGroupRoutes(globalPrefix, groupHandler)
}

func (suite *E2ETestSuite) SetupTest() {
	// Clean up database before each test
	suite.db.Exec("DELETE FROM user_groups")
	suite.db.Exec("DELETE FROM users")
	suite.db.Exec("DELETE FROM groups")

	// Create test user
	suite.testUser = &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	// Create test group
	suite.testGroup = &models.Group{
		Name: "Test Group",
	}

	suite.authToken = ""
}

func (suite *E2ETestSuite) TearDownSuite() {
	if suite.db != nil {
		sqlDB, _ := suite.db.DB()
		sqlDB.Close()
	}
}

// Helper method to make HTTP requests
func (suite *E2ETestSuite) makeRequest(method, url string, body interface{}, token ...string) (*http.Response, []byte) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	suite.Require().NoError(err)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if len(token) > 0 && token[0] != "" {
		req.Header.Set("Authorization", "Bearer "+token[0])
	}

	resp, err := suite.app.Test(req, -1)
	suite.Require().NoError(err)

	responseBody, err := io.ReadAll(resp.Body)
	suite.Require().NoError(err)

	return resp, responseBody
}

// Test User Handlers

func (suite *E2ETestSuite) TestCreateUser() {
	createUserReq := models.CreateUserRequest{
		Username: "newuser",
		Email:    "newuser@example.com",
		Password: "password123",
	}

	resp, body := suite.makeRequest("POST", "/api/users/", createUserReq)

	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var response map[string]interface{}
	err := json.Unmarshal(body, &response)
	suite.Require().NoError(err)

	assert.Contains(suite.T(), response, "token")
	assert.NotEmpty(suite.T(), response["token"])
}

func (suite *E2ETestSuite) TestCreateUserValidationErrors() {
	// Test with invalid email
	createUserReq := models.CreateUserRequest{
		Username: "test",
		Email:    "invalid-email",
		Password: "123",
	}

	resp, body := suite.makeRequest("POST", "/api/users/", createUserReq)

	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	err := json.Unmarshal(body, &response)
	suite.Require().NoError(err)

	assert.Contains(suite.T(), response, "errors")
}

func (suite *E2ETestSuite) TestLogin() {
	// First create a user
	createUserReq := models.CreateUserRequest{
		Username: "loginuser",
		Email:    "login@example.com",
		Password: "password123",
	}
	suite.makeRequest("POST", "/api/users/", createUserReq)

	// Now test login
	loginReq := models.LoginUserRequest{
		Email:    "login@example.com",
		Password: "password123",
	}

	resp, body := suite.makeRequest("POST", "/api/users/login", loginReq)

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err := json.Unmarshal(body, &response)
	suite.Require().NoError(err)

	assert.Contains(suite.T(), response, "token")
	assert.NotEmpty(suite.T(), response["token"])

	// Store token for authenticated tests
	suite.authToken = response["token"].(string)
}

func (suite *E2ETestSuite) TestLoginInvalidCredentials() {
	loginReq := models.LoginUserRequest{
		Email:    "nonexistent@example.com",
		Password: "wrongpassword",
	}

	resp, _ := suite.makeRequest("POST", "/api/users/login", loginReq)

	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
}

func (suite *E2ETestSuite) TestGetUsers() {
	// Create a few users first
	for i := 1; i <= 3; i++ {
		createUserReq := models.CreateUserRequest{
			Username: fmt.Sprintf("user%d", i),
			Email:    fmt.Sprintf("user%d@example.com", i),
			Password: "password123",
		}
		suite.makeRequest("POST", "/api/users/", createUserReq)
	}

	resp, body := suite.makeRequest("GET", "/api/users/", nil)

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err := json.Unmarshal(body, &response)
	suite.Require().NoError(err)

	assert.Contains(suite.T(), response, "users")
	users := response["users"].([]interface{})
	assert.Len(suite.T(), users, 3)
}

func (suite *E2ETestSuite) TestGetUserById() {
	// Create a user first
	createUserReq := models.CreateUserRequest{
		Username: "getuser",
		Email:    "getuser@example.com",
		Password: "password123",
	}
	suite.makeRequest("POST", "/api/users/", createUserReq)

	// Get the user ID from database
	var user models.User
	suite.db.Where("email = ?", "getuser@example.com").First(&user)

	resp, body := suite.makeRequest("GET", fmt.Sprintf("/api/users/%d", user.ID), nil)

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err := json.Unmarshal(body, &response)
	suite.Require().NoError(err)

	assert.Contains(suite.T(), response, "user")
	userData := response["user"].(map[string]interface{})
	assert.Equal(suite.T(), "getuser", userData["username"])
	assert.Equal(suite.T(), "getuser@example.com", userData["email"])
}

func (suite *E2ETestSuite) TestGetUserByIdNotFound() {
	resp, _ := suite.makeRequest("GET", "/api/users/999", nil)
	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)
}

func (suite *E2ETestSuite) TestUpdateUser() {
	// Create a user first
	createUserReq := models.CreateUserRequest{
		Username: "updateuser",
		Email:    "updateuser@example.com",
		Password: "password123",
	}
	suite.makeRequest("POST", "/api/users/", createUserReq)

	// Get the user ID from database
	var user models.User
	suite.db.Where("email = ?", "updateuser@example.com").First(&user)

	updateUserReq := models.UpdateUserRequest{
		Username: "updateduser",
		Email:    "updated@example.com",
	}

	resp, body := suite.makeRequest("PUT", fmt.Sprintf("/api/users/%d", user.ID), updateUserReq)

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err := json.Unmarshal(body, &response)
	suite.Require().NoError(err)

	assert.Equal(suite.T(), "User updated successfully", response["message"])

	// Verify the update in database
	var updatedUser models.User
	suite.db.First(&updatedUser, user.ID)
	assert.Equal(suite.T(), "updateduser", updatedUser.Username)
	assert.Equal(suite.T(), "updated@example.com", updatedUser.Email)
}

func (suite *E2ETestSuite) TestUpdateUserNotFound() {
	updateUserReq := models.UpdateUserRequest{
		Username: "updateduser",
	}

	resp, _ := suite.makeRequest("PUT", "/api/users/999", updateUserReq)
	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)
}

func (suite *E2ETestSuite) TestDeleteUser() {
	// Create a user first
	createUserReq := models.CreateUserRequest{
		Username: "deleteuser",
		Email:    "deleteuser@example.com",
		Password: "password123",
	}
	suite.makeRequest("POST", "/api/users/", createUserReq)

	// Get the user ID from database
	var user models.User
	suite.db.Where("email = ?", "deleteuser@example.com").First(&user)

	resp, body := suite.makeRequest("DELETE", fmt.Sprintf("/api/users/%d", user.ID), nil)

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err := json.Unmarshal(body, &response)
	suite.Require().NoError(err)

	assert.Equal(suite.T(), "User deleted successfully", response["message"])

	// Verify deletion in database
	var deletedUser models.User
	result := suite.db.First(&deletedUser, user.ID)
	assert.Error(suite.T(), result.Error)
}

func (suite *E2ETestSuite) TestDeleteUserNotFound() {
	resp, _ := suite.makeRequest("DELETE", "/api/users/999", nil)
	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)
}

func (suite *E2ETestSuite) TestGetMe() {
	// Create and login user first
	createUserReq := models.CreateUserRequest{
		Username: "meuser",
		Email:    "me@example.com",
		Password: "password123",
	}
	suite.makeRequest("POST", "/api/users/", createUserReq)

	loginReq := models.LoginUserRequest{
		Email:    "me@example.com",
		Password: "password123",
	}
	loginResp, loginBody := suite.makeRequest("POST", "/api/users/login", loginReq)
	assert.Equal(suite.T(), http.StatusOK, loginResp.StatusCode)

	var loginResponse map[string]interface{}
	err := json.Unmarshal(loginBody, &loginResponse)
	suite.Require().NoError(err)
	token := loginResponse["token"].(string)

	// Test GetMe with token
	resp, body := suite.makeRequest("GET", "/api/users/me", nil, token)

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	suite.Require().NoError(err)

	assert.Contains(suite.T(), response, "user")
	userData := response["user"].(map[string]interface{})
	assert.Equal(suite.T(), "meuser", userData["username"])
	assert.Equal(suite.T(), "me@example.com", userData["email"])
}

// Test Group Handlers

func (suite *E2ETestSuite) TestCreateGroup() {
	group := models.Group{
		Name: "New Test Group",
	}

	resp, body := suite.makeRequest("POST", "/api/groups/", group)

	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var response models.Group
	err := json.Unmarshal(body, &response)
	suite.Require().NoError(err)

	assert.Equal(suite.T(), "New Test Group", response.Name)
	assert.NotZero(suite.T(), response.ID)
}

func (suite *E2ETestSuite) TestCreateGroupInvalidBody() {
	resp, _ := suite.makeRequest("POST", "/api/groups/", "invalid-json")
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

func (suite *E2ETestSuite) TestGetAllGroups() {
	// Create a few groups first
	for i := 1; i <= 3; i++ {
		group := models.Group{
			Name: fmt.Sprintf("Group %d", i),
		}
		suite.makeRequest("POST", "/api/groups/", group)
	}

	resp, body := suite.makeRequest("GET", "/api/groups/", nil)

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var groups []models.Group
	err := json.Unmarshal(body, &groups)
	suite.Require().NoError(err)

	assert.Len(suite.T(), groups, 3)
}

func (suite *E2ETestSuite) TestGetGroupById() {
	// Create a group first
	group := models.Group{
		Name: "Get Test Group",
	}
	createResp, createBody := suite.makeRequest("POST", "/api/groups/", group)
	assert.Equal(suite.T(), http.StatusCreated, createResp.StatusCode)

	var createdGroup models.Group
	err := json.Unmarshal(createBody, &createdGroup)
	suite.Require().NoError(err)

	resp, body := suite.makeRequest("GET", fmt.Sprintf("/api/groups/%d", createdGroup.ID), nil)

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response models.Group
	err = json.Unmarshal(body, &response)
	suite.Require().NoError(err)

	assert.Equal(suite.T(), "Get Test Group", response.Name)
	assert.Equal(suite.T(), createdGroup.ID, response.ID)
}

func (suite *E2ETestSuite) TestGetGroupByIdNotFound() {
	resp, _ := suite.makeRequest("GET", "/api/groups/999", nil)
	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)
}

func (suite *E2ETestSuite) TestGetGroupByIdInvalidId() {
	resp, _ := suite.makeRequest("GET", "/api/groups/invalid", nil)
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

func (suite *E2ETestSuite) TestUpdateGroup() {
	// Create a group first
	group := models.Group{
		Name: "Update Test Group",
	}
	createResp, createBody := suite.makeRequest("POST", "/api/groups/", group)
	assert.Equal(suite.T(), http.StatusCreated, createResp.StatusCode)

	var createdGroup models.Group
	err := json.Unmarshal(createBody, &createdGroup)
	suite.Require().NoError(err)

	// Update the group
	updateGroup := models.Group{
		Name: "Updated Group Name",
	}

	resp, body := suite.makeRequest("PUT", fmt.Sprintf("/api/groups/%d", createdGroup.ID), updateGroup)

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response models.Group
	err = json.Unmarshal(body, &response)
	suite.Require().NoError(err)

	assert.Equal(suite.T(), "Updated Group Name", response.Name)
	assert.Equal(suite.T(), createdGroup.ID, response.ID)

	// Verify the update in database
	var updatedGroup models.Group
	suite.db.First(&updatedGroup, createdGroup.ID)
	assert.Equal(suite.T(), "Updated Group Name", updatedGroup.Name)
}

func (suite *E2ETestSuite) TestUpdateGroupInvalidId() {
	group := models.Group{
		Name: "Test Group",
	}

	resp, _ := suite.makeRequest("PUT", "/api/groups/invalid", group)
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

func (suite *E2ETestSuite) TestDeleteGroup() {
	// Create a group first
	group := models.Group{
		Name: "Delete Test Group",
	}
	createResp, createBody := suite.makeRequest("POST", "/api/groups/", group)
	assert.Equal(suite.T(), http.StatusCreated, createResp.StatusCode)

	var createdGroup models.Group
	err := json.Unmarshal(createBody, &createdGroup)
	suite.Require().NoError(err)

	resp, _ := suite.makeRequest("DELETE", fmt.Sprintf("/api/groups/%d", createdGroup.ID), nil)

	assert.Equal(suite.T(), http.StatusNoContent, resp.StatusCode)

	// Verify deletion in database
	var deletedGroup models.Group
	result := suite.db.First(&deletedGroup, createdGroup.ID)
	assert.Error(suite.T(), result.Error)
}

func (suite *E2ETestSuite) TestDeleteGroupInvalidId() {
	resp, _ := suite.makeRequest("DELETE", "/api/groups/invalid", nil)
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

func (suite *E2ETestSuite) TestAddUserToGroup() {
	// Create a user first
	createUserReq := models.CreateUserRequest{
		Username: "groupuser",
		Email:    "groupuser@example.com",
		Password: "password123",
	}
	suite.makeRequest("POST", "/api/users/", createUserReq)

	var user models.User
	suite.db.Where("email = ?", "groupuser@example.com").First(&user)

	// Create a group
	group := models.Group{
		Name: "User Group Test",
	}
	createResp, createBody := suite.makeRequest("POST", "/api/groups/", group)
	assert.Equal(suite.T(), http.StatusCreated, createResp.StatusCode)

	var createdGroup models.Group
	err := json.Unmarshal(createBody, &createdGroup)
	suite.Require().NoError(err)

	// Add user to group
	resp, body := suite.makeRequest("POST", fmt.Sprintf("/api/groups/%d/users/%d", createdGroup.ID, user.ID), nil)

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	suite.Require().NoError(err)

	assert.Equal(suite.T(), "user added to group", response["message"])

	// Verify in database
	var groupWithUsers models.Group
	suite.db.Preload("Users").First(&groupWithUsers, createdGroup.ID)
	assert.Len(suite.T(), groupWithUsers.Users, 1)
	assert.Equal(suite.T(), user.ID, groupWithUsers.Users[0].ID)
}

func (suite *E2ETestSuite) TestAddUserToGroupInvalidIds() {
	resp, _ := suite.makeRequest("POST", "/api/groups/invalid/users/invalid", nil)
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

func (suite *E2ETestSuite) TestRemoveUserFromGroup() {
	// Create a user first
	createUserReq := models.CreateUserRequest{
		Username: "removeuser",
		Email:    "removeuser@example.com",
		Password: "password123",
	}
	suite.makeRequest("POST", "/api/users/", createUserReq)

	var user models.User
	suite.db.Where("email = ?", "removeuser@example.com").First(&user)

	// Create a group
	group := models.Group{
		Name: "Remove User Group Test",
	}
	createResp, createBody := suite.makeRequest("POST", "/api/groups/", group)
	assert.Equal(suite.T(), http.StatusCreated, createResp.StatusCode)

	var createdGroup models.Group
	err := json.Unmarshal(createBody, &createdGroup)
	suite.Require().NoError(err)

	// Add user to group first
	suite.makeRequest("POST", fmt.Sprintf("/api/groups/%d/users/%d", createdGroup.ID, user.ID), nil)

	// Now remove user from group
	resp, body := suite.makeRequest("DELETE", fmt.Sprintf("/api/groups/%d/users/%d", createdGroup.ID, user.ID), nil)

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	suite.Require().NoError(err)

	assert.Equal(suite.T(), "user removed from group", response["message"])

	// Verify in database
	var groupWithUsers models.Group
	suite.db.Preload("Users").First(&groupWithUsers, createdGroup.ID)
	assert.Len(suite.T(), groupWithUsers.Users, 0)
}

func (suite *E2ETestSuite) TestRemoveUserFromGroupInvalidIds() {
	resp, _ := suite.makeRequest("DELETE", "/api/groups/invalid/users/invalid", nil)
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

// Test Runner
func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}