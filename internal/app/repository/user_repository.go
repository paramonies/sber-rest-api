package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/paramonies/sber-rest-api/internal/app/model"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user model.User) (model.User, error) {
	emptyUser := model.User{}

	findUserNameQuery := fmt.Sprintf("SELECT count(*) FROM %s WHERE name = $1", USERSTABLE)
	row := r.db.QueryRow(findUserNameQuery, user.Name)
	var count int
	if err := row.Scan(&count); err != nil {
		return emptyUser, err
	} else if count != 0 {
		return emptyUser, fmt.Errorf("user %s already exixts", user.Name)
	}

	selectUserTypeQuery := fmt.Sprintf("SELECT id FROM %s WHERE id = $1", USERTYPESTABLE)
	row = r.db.QueryRow(selectUserTypeQuery, user.UserType)
	var utId int
	if err := row.Scan(&utId); err != nil {
		switch {
		case err == sql.ErrNoRows:
			return emptyUser, fmt.Errorf("user-type for %d not found", user.UserType)
		default:
			return emptyUser, err
		}
	}

	tx, err := r.db.Begin()
	if err != nil {
		return emptyUser, err
	}
	var userId int
	createUserQuery := fmt.Sprintf("INSERT INTO %s (name, age, user_type_id) VALUES ($1, $2, $3) RETURNING id", USERSTABLE)
	row = r.db.QueryRow(createUserQuery, user.Name, user.Age, user.UserType)
	if err := row.Scan(&userId); err != nil {
		tx.Rollback()
		return emptyUser, err
	}
	user.Id = userId
	user.Created = time.Now().Format("2006-01-02 15:04:05")
	user.Updated = time.Now().Format("2006-01-02 15:04:05")
	var itemId int
	for i, v := range user.Items {
		createItemQuery := fmt.Sprintf("INSERT INTO %s (name, user_id) VALUES ($1, $2) RETURNING id", ITEMSTABLE)
		row := r.db.QueryRow(createItemQuery, v.Name, userId)
		if err := row.Scan(&itemId); err != nil {
			tx.Rollback()
			return emptyUser, err
		}
		user.Items[i].Id = itemId
		user.Items[i].UserId = userId
		user.Items[i].Created = time.Now().Format("2006-01-02 15:04:05")
		user.Items[i].Updated = time.Now().Format("2006-01-02 15:04:05")
	}

	return user, tx.Commit()
}

func (r *UserRepository) GetUserById(id int) (model.User, error) {
	return getUserResponse(r, id)
}

func (r *UserRepository) UpdateUser(user model.UpdateUser) (model.User, error) {
	emptyUser := model.User{}

	userId, err := verifyIfUserExists(r.db, *(user.Id))
	if err != nil {
		return emptyUser, err
	}

	//verify if user-type exists
	selectUserTypeQuery := fmt.Sprintf("SELECT id FROM %s WHERE id = $1", USERTYPESTABLE)
	row := r.db.QueryRow(selectUserTypeQuery, *user.UserType)
	var utId int
	if err := row.Scan(&utId); err != nil {
		switch {
		case err == sql.ErrNoRows:
			return emptyUser, fmt.Errorf("user-type for %d not found", *user.UserType)
		default:
			return emptyUser, err
		}
	}

	var tx *sql.Tx

	if user.Items != nil {
		tx, err = r.db.Begin()
		if err != nil {
			return emptyUser, err
		}
	}

	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if user.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *user.Name)
		argId++
	}

	if user.Age != nil {
		setValues = append(setValues, fmt.Sprintf("age=$%d", argId))
		args = append(args, *user.Age)
		argId++
	}

	if user.UserType != nil {
		setValues = append(setValues, fmt.Sprintf("user_type_id=$%d", argId))
		args = append(args, *user.UserType)
		argId++
	}

	//update updated_at db field
	setValues = append(setValues, fmt.Sprintf("updated_at=$%d", argId))
	args = append(args, time.Now())
	argId++

	setQuery := strings.Join(setValues, ", ")
	updateUserQuery := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d", USERSTABLE, setQuery, argId)
	args = append(args, *user.Id)

	_, err = r.db.Exec(updateUserQuery, args...)
	if err != nil {
		if user.Items != nil {
			tx.Rollback()
		}
		return emptyUser, err
	}

	//Create Items for updated user
	if user.Items != nil {
		for _, v := range *user.Items {
			createItemQuery := fmt.Sprintf("INSERT INTO %s (name, user_id) VALUES ($1, $2) RETURNING id", ITEMSTABLE)
			_, err := r.db.Exec(createItemQuery, v.Name, userId)
			if err != nil {
				tx.Rollback()
				return emptyUser, err
			}
		}
		tx.Commit()
	}

	//make response User

	return getUserResponse(r, *user.Id)
}

func (r *UserRepository) DeleteUser(id int) error {
	if _, err := verifyIfUserExists(r.db, id); err != nil {
		return err
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	deleteItemsQuery := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1", ITEMSTABLE)
	_, err = r.db.Exec(deleteItemsQuery, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	deleteUserQuery := fmt.Sprintf("DELETE FROM %s WHERE id = $1", USERSTABLE)
	_, err = r.db.Exec(deleteUserQuery, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *UserRepository) GetListUsers(page, limit int) ([]model.User, error) {
	// emptyUser := model.User{}
	var usersList []model.User
	selectUsersQuery := fmt.Sprintf("Select * FROM %s ORDER BY created_at DESC LIMIT $1 OFFSET ($2-1)*$1", USERSTABLE)
	rows, err := r.db.Queryx(selectUsersQuery, limit, page)
	if err != nil {
		return nil, err
	}
	var user model.User
	var item model.Item
	var listItems []model.Item
	for rows.Next() {
		rows.StructScan(&user)
		// fmt.Printf("%d %s %d \n", user.Id, user.Name, user.Age)
		user.Created = formatTime(user.Created)
		user.Updated = formatTime(user.Updated)

		//get items for user
		selectItemsQuery := fmt.Sprintf("Select * From %s WHERE user_id = $1", ITEMSTABLE)
		rows, err := r.db.Queryx(selectItemsQuery, user.Id)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			rows.StructScan(&item)
			item.Created = formatTime(item.Created)
			item.Updated = formatTime(item.Updated)
			listItems = append(listItems, item)
		}

		user.Items = listItems
		usersList = append(usersList, user)
	}
	return usersList, nil
}

func getUserResponse(r *UserRepository, id int) (model.User, error) {
	emptyUser := model.User{}
	var userResponse model.User
	selectUserQuery := fmt.Sprintf("Select * From %s WHERE id = $1", USERSTABLE)
	err := r.db.Get(&userResponse, selectUserQuery, id)
	if err != nil {
		return emptyUser, err
	}

	userResponse.Created = formatTime(userResponse.Created)
	userResponse.Updated = formatTime(userResponse.Updated)

	var item model.Item
	var listItems []model.Item
	selectItemsQuery := fmt.Sprintf("Select * From %s WHERE user_id = $1", ITEMSTABLE)
	rows, err := r.db.Queryx(selectItemsQuery, id)
	if err != nil {
		return emptyUser, err
	}
	for rows.Next() {
		rows.StructScan(&item)
		item.Created = formatTime(item.Created)
		item.Updated = formatTime(item.Updated)
		listItems = append(listItems, item)
	}

	userResponse.Items = listItems
	return userResponse, nil
}
