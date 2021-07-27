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

type ItemRepositoryIf interface {
	CreateItem(item model.Item) (model.Item, error)
	UpdateItem(itemId int, item model.UpdateItem) (model.Item, error)
}

type Repository struct {
	UserRepositoryIf
	ItemRepositoryIf
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		UserRepositoryIf: NewUserRepository(db),
		ItemRepositoryIf: NewItemRepository(db),
	}
}
