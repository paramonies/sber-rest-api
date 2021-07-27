package service

import (
	"github.com/paramonies/sber-rest-api/internal/app/model"
	"github.com/paramonies/sber-rest-api/internal/app/repository"
)

type ItemService struct {
	repo repository.ItemRepositoryIf
}

func NewItemService(repo repository.ItemRepositoryIf) *ItemService {
	return &ItemService{repo: repo}
}

func (s *ItemService) CreateItem(item model.Item) (model.Item, error) {
	return s.repo.CreateItem(item)
}

func (s *ItemService) UpdateItem(itemId int, input model.UpdateItem) (model.Item, error) {
	return s.repo.UpdateItem(itemId, input)
}
