package service

import (
	"felix1234567890/go-trello/models"
	"felix1234567890/go-trello/repository"
)

type UserService struct {
	Repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		Repo: repo,
	}
}
func (s *UserService) GetUsers() ([]models.User, error) {
	return s.Repo.GetUsers()
}

func (s *UserService) GetUserById(id string) (models.User, error) {
	return s.Repo.GetUserById(id)
}
