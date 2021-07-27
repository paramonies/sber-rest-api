package service

import (
	"github.com/paramonies/sber-rest-api/internal/app/model"
	"github.com/paramonies/sber-rest-api/internal/app/repository"
)

type UserService struct {
	repo repository.UserRepositoryIf
}

func NewUserService(repo repository.UserRepositoryIf) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(user model.User) (model.User, error) {
	user, err := s.repo.CreateUser(user)
	return user, err
}

func (s *UserService) GetUserById(id int) (model.User, error) {
	user, err := s.repo.GetUserById(id)
	return user, err
}

func (s *UserService) UpdateUser(user model.UpdateUser) (model.User, error) {
	updatedUser, err := s.repo.UpdateUser(user)
	return updatedUser, err
}

func (s *UserService) DeleteUser(id int) error {
	return s.repo.DeleteUser(id)
}

func (s *UserService) GetListUsers(page, limit int) ([]model.User, error) {
	return s.repo.GetListUsers(page, limit)
}
