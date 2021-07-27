package model

type Item struct {
	Id      int    `json:"id" db:"id"`
	Name    string `json:"name" binding:"required" db:"name"`
	UserId  int    `json:"user-id" binding:"required" db:"user_id"`
	Created string `json:"created" db:"created_at"`
	Updated string `json:"updated" db:"updated_at"`
}

type UpdateItem struct {
	Name   *string `json:"name" binding:"required" db:"name"`
	UserId *int    `json:"user-id" binding:"required" db:"user_id"`
}
