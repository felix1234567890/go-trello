package repository

import (
	"felix1234567890/go-trello/models"
	"felix1234567890/go-trello/utils"

	"gorm.io/gorm"
)

// UserRepository defines the interface for user persistence operations.
type UserRepository interface {
	// GetUsers retrieves all users from the database.
	GetUsers() ([]models.User, error)
	// GetUserById retrieves a single user by their ID.
	// Returns the found user or an error (e.g., gorm.ErrRecordNotFound if not found).
	GetUserById(id string) (models.User, error)
	// DeleteUser removes a user by their ID.
	// Returns an error if the operation fails or the user is not found.
	DeleteUser(id string) error
	// UpdateUser updates an existing user's details by their ID.
	// Takes a models.UpdateUserRequest containing fields to update.
	// Returns an error if the operation fails or the user is not found.
	UpdateUser(id string, req *models.UpdateUserRequest) error
	// CreateUser adds a new user to the database.
	// The input user's password should be plain text; it will be hashed before saving.
	// Returns the ID of the newly created user and an error if the operation fails.
	CreateUser(req *models.User) (uint, error)
	// LoginUser authenticates a user based on email and password.
	// Returns the user's ID on successful authentication, or an error otherwise.
	LoginUser(req *models.LoginUserRequest) (uint, error)
}

// userRepositoryImpl is the concrete implementation of UserRepository.
type userRepositoryImpl struct {
	DB *gorm.DB
}

// NewUserRepository creates a new instance of the concrete UserRepository implementation.
func NewUserRepository(db *gorm.DB) UserRepository { // Returns interface
	return &userRepositoryImpl{ // Ensure this returns a pointer if methods have pointer receivers
		DB: db,
	}
}

func (r *userRepositoryImpl) GetUsers() ([]models.User, error) {
	var users []models.User
	if err := r.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepositoryImpl) GetUserById(id string) (models.User, error) { // Corrected receiver
	var user models.User
	if err := r.DB.First(&user, id).Error; err != nil {
		return models.User{}, err // Return zero-value User on error
	}
	return user, nil
}

func (r *userRepositoryImpl) DeleteUser(id string) error {
	result := r.DB.Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *userRepositoryImpl) UpdateUser(id string, req *models.UpdateUserRequest) error { // Corrected receiver
	result := r.DB.Model(&models.User{}).Where("id = ?", id).Updates(req) // Pass req directly
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *userRepositoryImpl) CreateUser(req *models.User) (uint, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return 0, err
	}
	req.Password = hashedPassword
	result := r.DB.Create(&req)
	if result.Error != nil {
		return 0, result.Error
	}
	return req.ID, nil
}

// LoginUser attempts to authenticate a user.
// Renamed from Login to LoginUser for consistency with UserService interface.
func (r *userRepositoryImpl) LoginUser(loginUserRequest *models.LoginUserRequest) (uint, error) { // Corrected receiver
	var user models.User
	result := r.DB.Where("email = ?", loginUserRequest.Email).First(&user)
	if result.Error != nil {
		return 0, result.Error // Could be gorm.ErrRecordNotFound if email not found
	}
	if err := utils.CheckPasswordHash(loginUserRequest.Password, user.Password); err != nil {
		return 0, err // Password mismatch
	}
	return user.ID, nil
}
