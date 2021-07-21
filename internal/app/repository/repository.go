package repository

import "github.com/jmoiron/sqlx"

type UserRepositoryIf interface {
	CreateUser()
}

type Repository struct {
	UserRepositoryIf
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		UserRepositoryIf: NewUserRepository(db),
	}
}
