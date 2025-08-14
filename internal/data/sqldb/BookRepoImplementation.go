package sqldb

import (
	"database/sql"
	"fmt"
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

func (r *SQLBookRepo) Get(id string) (entity.Book, error) {
	var book entity.Book
	query := `select * from books where id=$1 `
	row := r.DB.QueryRow(query, id)
	if err := row.Scan(&book.Id, &book.Name, &book.Author, &book.PriceCents, &book.ISBN, &book.CreatedAt, &book.UpdatedAt, &book.PublishedAt); err != nil {
		if err == sql.ErrNoRows {
			return book, fmt.Errorf("no book found with id %s", id)
		} else {
			return book, err
		}
	}
	return book, nil
}

func (r *SQLBookRepo) List() ([]entity.Book, error) {
	var books []entity.Book
	query := `select * from books`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var book entity.Book
		if err := rows.Scan(&book.Id, &book.Name, &book.Author, &book.PriceCents, &book.ISBN, &book.CreatedAt, &book.UpdatedAt, &book.PublishedAt); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return books, nil
}
