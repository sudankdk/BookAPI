package book

import "github.com/sudankdk/bookstore/internal/domain/entity"

func (c *CreateBookUsecase) Update(id string, data map[string]any) (entity.Book, error) {
	return c.repo.Update(id, data)
}
