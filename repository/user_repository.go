package repository

import (
	"felix1234567890/go-trello/models"
	"felix1234567890/go-trello/utils"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}
func (r *UserRepository) GetUsers() ([]models.User, error) {
	var users []models.User
	if err := r.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) GetUserById(id string) (models.User, error) {
	var user models.User
	if err := r.DB.First(&user, id).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *UserRepository) DeleteUser(id string) error {
	result := r.DB.Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *UserRepository) UpdateUser(id string, req *models.UpdateUserRequest) error {
	result := r.DB.Model(&models.User{}).Where("id = ?", id).Updates(&req)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *UserRepository) CreateUser(req *models.User) (uint, error) {
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

func (r *UserRepository) Login(LoginUserRequest *models.LoginUserRequest) (uint, error) {

	var user models.User
	result := r.DB.Where("email = ?", LoginUserRequest.Email).First(&user)
	if result.Error != nil {
		return 0, result.Error
	}
	if err := utils.CheckPasswordHash(LoginUserRequest.Password, user.Password); err != nil {
		return 0, err
	}
	return user.ID, nil
}
