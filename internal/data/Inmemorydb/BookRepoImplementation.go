package inmemorydb

import (
	"sync"

	"github.com/sudankdk/bookstore/internal/domain/entity"
)

type InMemoryBookRepo struct {
	books map[string]entity.Book
	mu    sync.Mutex
}

func NewInMemoryBookRepo() *InMemoryBookRepo {
	return &InMemoryBookRepo{
		books: make(map[string]entity.Book),
	}
}

func (r *InMemoryBookRepo) Create(b entity.Book) (entity.Book, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.books[b.Id] = b
	return b, nil
}
