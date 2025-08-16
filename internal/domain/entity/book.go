package entity

import (
	"time"
)

type Book struct {
	Id          string    `json:"id" `
	Name        string    `json:"name" validate:"required,min=2,max=200"`
	Author      string    `json:"author" validate:"required"`
	PriceCents  int64     `json:"price_cents" validate:"required,gt=0"`
	ISBN        string    `json:"isbn"`
	PublishedAt time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// func (b *Book) Validate() error {
// 	validate := validator.New()
// 	return validate.Struct(b)
// }
