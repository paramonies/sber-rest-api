package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/paramonies/sber-rest-api/internal/app/model"
)

type UserRepositoryIf interface {
	CreateUser(user model.User) (model.User, error)
	GetUserById(id int) (model.User, error)
	UpdateUser(user model.UpdateUser) (model.User, error)
	DeleteUser(id int) error
	GetListUsers(page, limit int) ([]model.User, error)
}

type Repository struct {
	UserRepositoryIf
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		UserRepositoryIf: NewUserRepository(db),
	}
}
