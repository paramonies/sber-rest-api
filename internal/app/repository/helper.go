package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	USERSTABLE     = "users"
	ITEMSTABLE     = "items"
	USERTYPESTABLE = "user_types"
)

func verifyIfUserExists(repoDB *sqlx.DB, id int) (int, error) {
	findUserQuery := fmt.Sprintf("SELECT id FROM %s WHERE id = $1", USERSTABLE)
	row := repoDB.QueryRow(findUserQuery, id)
	var userId int
	if err := row.Scan(&userId); err != nil {
		switch {
		case err == sql.ErrNoRows:
			return -1, fmt.Errorf("user with id %d not found", id)
		default:
			return -1, err
		}
	}
	return userId, nil
}

func verifyIfItemExists(repoDB *sqlx.DB, id int) (int, error) {
	findItemQuery := fmt.Sprintf("SELECT id FROM %s WHERE id = $1", ITEMSTABLE)
	row := repoDB.QueryRow(findItemQuery, id)
	var itemId int
	if err := row.Scan(&itemId); err != nil {
		switch {
		case err == sql.ErrNoRows:
			return -1, fmt.Errorf("item with id %d not found", id)
		default:
			return -1, err
		}
	}
	return itemId, nil
}

func formatTime(timeStr string) string {
	timeT, _ := time.Parse("2006-01-02T15:04:05Z", timeStr)
	return timeT.Format("2006-01-02 15:04:05")
}
