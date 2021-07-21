package service

import "github.com/paramonies/sber-rest-api/internal/app/repository"

type UserServiceIf interface {
	CreateUser()
}

type Service struct {
	UserServiceIf
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		UserServiceIf: NewUserService(repos.UserRepositoryIf),
	}
}
