package service

import (
	"errors"
	"felix1234567890/go-trello/models"
	"felix1234567890/go-trello/repository" // To define the interface
	"testing"
	// "time" // No longer needed if expectedUserAfterCreation is removed

	"github.com/stretchr/testify/assert" // Using testify for assertions
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	CreateUserFunc func(user *models.User) (uint, error) // Corrected signature
	// Add other methods from UserRepository interface as needed for other tests
	GetUsersFunc    func() ([]models.User, error)
	GetUserByIdFunc func(id string) (models.User, error) 
	UpdateUserFunc  func(id string, user *models.UpdateUserRequest) error // Changed to return error
	DeleteUserFunc  func(id string) error
	LoginUserFunc   func(user *models.LoginUserRequest) (uint, error) // This field will be used by LoginUser method
}

// Ensure MockUserRepository implements repository.UserRepository
var _ repository.UserRepository = (*MockUserRepository)(nil) // Corrected check

func (m *MockUserRepository) CreateUser(user *models.User) (uint, error) { // Return uint, error to match interface
	if m.CreateUserFunc != nil {
		// The mock's CreateUserFunc should simulate what the real repo's CreateUser does:
		// 1. Hashing is done by real repo's CreateUser, so mock doesn't need to re-hash.
		// 2. It returns (id, error).
		// The service test for CreateUser is testing the *service's* logic,
		// not the repo's hashing logic.
		// So, the mock's CreateUserFunc should just return a predefined ID and error.
		// The previous mock returned (*models.User, error) which was for the old repo signature.
		// Let's assume the mock's CreateUserFunc is now defined to return (uint, error)
		return m.CreateUserFunc(user) 
	}
	panic("CreateUserFunc not implemented in mock")
}


func (m *MockUserRepository) GetUsers() ([]models.User, error) {
	if m.GetUsersFunc != nil {
		return m.GetUsersFunc()
	}
	panic("GetUsersFunc not implemented")
}

func (m *MockUserRepository) GetUserById(id string) (models.User, error) { // Changed to return models.User
	if m.GetUserByIdFunc != nil {
		return m.GetUserByIdFunc(id)
	}
	panic("GetUserByIdFunc not implemented")
}

func (m *MockUserRepository) UpdateUser(id string, user *models.UpdateUserRequest) error { // Changed to return error
	if m.UpdateUserFunc != nil {
		return m.UpdateUserFunc(id, user)
	}
	panic("UpdateUserFunc not implemented")
}

func (m *MockUserRepository) DeleteUser(id string) error {
	if m.DeleteUserFunc != nil {
		return m.DeleteUserFunc(id)
	}
	panic("DeleteUserFunc not implemented")
}

// LoginUser implements the UserRepository interface for the mock.
func (m *MockUserRepository) LoginUser(user *models.LoginUserRequest) (uint, error) {
	if m.LoginUserFunc != nil {
		return m.LoginUserFunc(user)
	}
	panic("LoginUserFunc not implemented for mock")
}


func TestUserServiceImpl_CreateUser(t *testing.T) {
	validUserInput := &models.User{
		// ID will be 0 by default
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	// This is what the *service* is expected to receive from the *repository*
	// The repository's CreateUser method returns (uint, error)
	expectedIDFromRepo := uint(1)

	tests := []struct {
		name            string
		userInput       *models.User
		setupMockCreate func(*models.User) (uint, error) // Mock for CreateUser in repo
		expectedID      uint
		expectedSvcErr  error // Expected error from the service call
	}{
		{
			name:      "successful creation",
			userInput: validUserInput,
			setupMockCreate: func(user *models.User) (uint, error) {
				// This mock simulates the repository's CreateUser.
				// It receives a user (password should be plain text here as service passes it as is)
				// It should return the ID and nil error for success.
				// The actual hashing and DB save is "done" by the real repo.
				assert.Equal(t, "password123", user.Password, "Password passed to repo should be plain text")
				return expectedIDFromRepo, nil
			},
			expectedID:     expectedIDFromRepo,
			expectedSvcErr: nil,
		},
		{
			name:      "repository returns error",
			userInput: validUserInput,
			setupMockCreate: func(user *models.User) (uint, error) {
				return 0, errors.New("repository error")
			},
			expectedID:     0,
			expectedSvcErr: errors.New("repository error"),
		},
		{
			name:      "repository returns 0 ID without error (current service behavior passes this through)",
			userInput: validUserInput,
			setupMockCreate: func(user *models.User) (uint, error) {
				return 0, nil // Simulate repo returning 0 ID but no error
			},
			expectedID:     0,
			expectedSvcErr: nil, // Reflecting current service behavior: it doesn't error on (0, nil) from repo.
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{}
			mockRepo.CreateUserFunc = tt.setupMockCreate

			userService := NewUserService(mockRepo)
			
			id, err := userService.CreateUser(tt.userInput)

			assert.Equal(t, tt.expectedID, id)
			if tt.expectedSvcErr != nil {
				assert.Error(t, err, "Expected an error but got nil for test case: %s", tt.name)
				if err != nil { // Check err is not nil before calling err.Error()
					assert.Contains(t, err.Error(), tt.expectedSvcErr.Error(), "Error message mismatch for test case: %s", tt.name)
				}
			} else {
				assert.NoError(t, err, "Expected no error but got one for test case: %s", tt.name)
			}
		})
	}
}
