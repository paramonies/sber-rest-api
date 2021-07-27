package service

import (
	"github.com/paramonies/sber-rest-api/internal/app/model"
	"github.com/paramonies/sber-rest-api/internal/app/repository"
)

type UserServiceIf interface {
	CreateUser(user model.User) (model.User, error)
	GetUserById(id int) (model.User, error)
	UpdateUser(user model.UpdateUser) (model.User, error)
	DeleteUser(id int) error
	GetListUsers(page, limit int) ([]model.User, error)
}

type Service struct {
	UserServiceIf
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		UserServiceIf: NewUserService(repos.UserRepositoryIf),
	}
}
