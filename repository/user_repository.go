package repository

import (
	"felix1234567890/go-trello/models"

	"gorm.io/gorm"
)

type IUserRepository interface {
	GetUsers() ([]models.User, error)
}
type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
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
