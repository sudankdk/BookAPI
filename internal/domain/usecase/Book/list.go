package book

import "github.com/sudankdk/bookstore/internal/domain/entity"

func (c *CreateBookUsecase) List() ([]entity.Book, error) {
	return c.repo.List()
}
