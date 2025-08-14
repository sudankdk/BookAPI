package sqldb

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/sudankdk/bookstore/internal/domain/entity"
)

type SQLBookRepo struct {
	DB *sql.DB
}

func NewSqlBookRepo(db *sql.DB) *SQLBookRepo {
	return &SQLBookRepo{DB: db}
}

func (r *SQLBookRepo) Create(b entity.Book) (entity.Book, error) {
	if b.Id == "" {
		b.Id = uuid.NewString()
	}
	now := time.Now()
	b.CreatedAt = now
	b.UpdatedAt = now

	query := `INSERT INTO books (id, name, author, price_cents, isbn, published_at, created_at, updated_at)
			  VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err := r.DB.Exec(query, b.Id, b.Name, b.Author, b.PriceCents, b.ISBN, b.PublishedAt, b.CreatedAt, b.UpdatedAt)
	if err != nil {
		return entity.Book{}, err
	}
	return b, nil

}
