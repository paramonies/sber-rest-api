package model

type User struct {
	Id       int    `json:"id,omitempty" db:"id"`
	Name     string `json:"name" db:"name" binding:"required"`
	Age      int    `json:"age" db:"age" binding:"required"`
	UserType int    `json:"user-type" db:"user_type_id" binding:"required"`
	Created  string `json:"created,omitempty" db:"created_at"`
	Updated  string `json:"updated,omitempty" db:"updated_at"`
	Items    []Item `json:"items,omitempty"`
}

type UpdateUser struct {
	Id       *int    `json:"id" binding:"required"`
	Name     *string `json:"name,omitempty"`
	Age      *int    `json:"age,omitempty"`
	UserType *int    `json:"user-type,omitempty"`
	Items    *[]Item `json:"items,omitempty"`
}
