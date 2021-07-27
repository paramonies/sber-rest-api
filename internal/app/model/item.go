package model

type Item struct {
	Id      int    `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	UserId  int    `json:"user-id" db:"user_id"`
	Created string `json:"created" db:"created_at"`
	Updated string `json:"updated" db:"updated_at"`
}
