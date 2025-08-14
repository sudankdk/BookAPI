package bookdto

import "time"

type CreateBookDTO struct {
	Name        string    `json:"name"`
	Author      string    `json:"author"`
	PriceCents  int64     `json:"price_cents"`
	ISBN        string    `json:"isbn"`
	PublishedAt time.Time `json:"published_at"`
}
