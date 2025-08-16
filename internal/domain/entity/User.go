package entity

import (
	"time"
)

type User struct {
	Id        string    `json:"id"`
	Username  string    `json:"username" validate:"required,max=200,min=3"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// func (b *User) Validate() error {
// 	validate := validator.New()
// 	return validate.Struct(b)
// }
