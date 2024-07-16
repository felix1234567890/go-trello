package service

import (
	"felix1234567890/go-trello/models"
	"felix1234567890/go-trello/repository"
)

type UserService interface {
	GetUsers() ([]models.User, error)
	GetUserById(id string) (models.User, error)
	DeleteUser(id string) error
	UpdateUser(id string, req *models.UpdateUserRequest) error
	CreateUser(req *models.User) (uint, error)
	LoginUser(req *models.LoginUserRequest) (uint, error)
}
type UserServiceImpl struct {
	Repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserServiceImpl {
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

func (s *UserServiceImpl) LoginUser(LoginUserRequest *models.LoginUserRequest) (uint, error) {
	return s.Repo.Login(LoginUserRequest)
}
