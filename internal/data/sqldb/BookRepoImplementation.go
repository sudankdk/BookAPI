package sqldb

import (
	"database/sql"
	"fmt"
	"strings"
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

// Delete implements repo.BookAbstractRepo.
func (r *SQLBookRepo) Delete(id string) error {
	query := `Delete from books where id = $1`
	res, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no book found with id %s", id)
	}

	return nil
}

// Update implements repo.BookAbstractRepo.
func (r *SQLBookRepo) Update(id string, data map[string]any) (entity.Book, error) {
	var book entity.Book
	fields := []string{}
	values := []any{}
	i := 1

	for k, v := range data {
		fields = append(fields, fmt.Sprintf("%s=$%d", k, i))
		values = append(values, v)
		i++
	}
	fields = append(fields, fmt.Sprintf("updated_at=$%d", i))
	values = append(values, time.Now())
	i++

	query := fmt.Sprintf(
		"UPDATE books SET %s WHERE id=$%d RETURNING id, name, author, price_cents, isbn, published_at, updated_at",
		strings.Join(fields, ", "),
		i,
	)

	values = append(values, id)
	err := r.DB.QueryRow(query, values...).Scan(
		&book.Id,
		&book.Name,
		&book.Author,
		&book.PriceCents,
		&book.ISBN,
		&book.PublishedAt,
		&book.UpdatedAt,
	)
	if err != nil {
		return book, err
	}

	return book, nil
}
