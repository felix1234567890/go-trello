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

func (s *UserService) DeleteUser(id string) error {
	return s.Repo.DeleteUser(id)
}
func (s *UserService) UpdateUser(id string, req *models.UpdateUserRequest) error {
	return s.Repo.UpdateUser(id, req)
}
func (s *UserService) CreateUser(req *models.User) (uint, error) {
	return s.Repo.CreateUser(req)
}

func (s *UserService) Login(LoginUserRequest *models.LoginUserRequest) (uint, error) {
	return s.Repo.Login(LoginUserRequest)
}
