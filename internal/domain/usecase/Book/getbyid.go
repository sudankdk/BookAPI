package book

import "github.com/sudankdk/bookstore/internal/domain/entity"

func (c *CreateBookUsecase) GetById(id string) (entity.Book, error) {
	return c.repo.Get(id)
}
