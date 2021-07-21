package service

import "github.com/paramonies/sber-rest-api/internal/app/repository"

type UserService struct {
	repo repository.UserRepositoryIf
}

func NewUserService(repo repository.UserRepositoryIf) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser() {

}
