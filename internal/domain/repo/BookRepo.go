package repo

import "github.com/sudankdk/bookstore/internal/domain/entity"

type BookAbstractRepo interface {
	List() ([]entity.Book, error)
	Get(id string) (entity.Book, error)
	Create(b entity.Book) (entity.Book, error)
	Update(id string, data map[string]any) (entity.Book, error)
	Delete(id string) error
}
