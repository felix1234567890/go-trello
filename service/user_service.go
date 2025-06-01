package service

import (
	"felix1234567890/go-trello/models"
	"felix1234567890/go-trello/repository" // This will now refer to the package containing the interface
)

// UserService defines the business logic operations related to users.
type UserService interface {
	// GetUsers retrieves all users.
	GetUsers() ([]models.User, error)
	// GetUserById retrieves a user by their ID.
	GetUserById(id string) (models.User, error)
	// DeleteUser removes a user by their ID.
	DeleteUser(id string) error
	// UpdateUser updates an existing user's details.
	UpdateUser(id string, req *models.UpdateUserRequest) error
	// CreateUser handles the creation of a new user.
	// Input user's password should be plain text.
	// Returns the ID of the created user or an error.
	CreateUser(req *models.User) (uint, error)
	// LoginUser handles user authentication.
	// Returns the user's ID on successful login or an error.
	LoginUser(req *models.LoginUserRequest) (uint, error)
}

// UserServiceImpl is the concrete implementation of the UserService interface.
type UserServiceImpl struct {
	Repo repository.UserRepository // Changed to depend on the interface
}

// NewUserService creates a new UserServiceImpl.
func NewUserService(repo repository.UserRepository) UserService { // Changed to accept interface and return interface
	return &UserServiceImpl{
		Repo: repo,
	}
}

func (s *UserServiceImpl) GetUsers() ([]models.User, error) {
	return s.Repo.GetUsers()
}

func (s *UserServiceImpl) GetUserById(id string) (models.User, error) {
	return s.Repo.GetUserById(id)
}

func (s *UserServiceImpl) DeleteUser(id string) error {
	return s.Repo.DeleteUser(id)
}
func (s *UserServiceImpl) UpdateUser(id string, req *models.UpdateUserRequest) error {
	return s.Repo.UpdateUser(id, req)
}
func (s *UserServiceImpl) CreateUser(req *models.User) (uint, error) {
	return s.Repo.CreateUser(req)
}

func (s *UserServiceImpl) LoginUser(loginUserRequest *models.LoginUserRequest) (uint, error) {
	// Ensure the method called matches the repository interface
	return s.Repo.LoginUser(loginUserRequest) 
}
