package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/paramonies/sber-rest-api/internal/app/model"
)

type ItemRepository struct {
	db *sqlx.DB
}

func NewItemRepository(db *sqlx.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) CreateItem(item model.Item) (model.Item, error) {
	emptyItem := model.Item{}
	_, err := verifyIfUserExists(r.db, item.UserId)
	if err != nil {
		return emptyItem, err
	}

	var itemId int
	createQuery := fmt.Sprintf("INSERT INTO %s (name, user_id) VALUES ($1, $2) RETURNING id", ITEMSTABLE)
	row := r.db.QueryRow(createQuery, item.Name, item.UserId)
	if err := row.Scan(&itemId); err != nil {
		return emptyItem, err
	}
	item.Id = itemId
	item.Created = time.Now().Format("2006-01-02 15:04:05")
	item.Updated = time.Now().Format("2006-01-02 15:04:05")

	return item, nil
}

func (r *ItemRepository) UpdateItem(itemId int, item model.UpdateItem) (model.Item, error) {
	emptyItem := model.Item{}
	_, err := verifyIfItemExists(r.db, itemId)
	if err != nil {
		return emptyItem, err
	}

	_, err = verifyIfUserExists(r.db, *item.UserId)
	if err != nil {
		return emptyItem, err
	}

	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if item.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *item.Name)
		argId++
	}

	if item.UserId != nil {
		setValues = append(setValues, fmt.Sprintf("user_id=$%d", argId))
		args = append(args, *item.UserId)
		argId++
	}

	//update updated_at db field
	setValues = append(setValues, fmt.Sprintf("updated_at=$%d", argId))
	args = append(args, time.Now())
	argId++

	setQuery := strings.Join(setValues, ", ")
	updateItemQuery := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d", ITEMSTABLE, setQuery, argId)
	args = append(args, itemId)

	_, err = r.db.Exec(updateItemQuery, args...)
	if err != nil {
		return emptyItem, err
	}

	//get item fot response
	var itemResponse model.Item
	selectItemQuery := fmt.Sprintf("Select * From %s WHERE id = $1", ITEMSTABLE)
	err = r.db.Get(&itemResponse, selectItemQuery, itemId)
	if err != nil {
		return emptyItem, err
	}

	itemResponse.Created = formatTime(itemResponse.Created)
	itemResponse.Updated = formatTime(itemResponse.Updated)

	return itemResponse, nil
}
