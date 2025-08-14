package book

import (
	"github.com/sudankdk/bookstore/internal/domain/entity"
	"github.com/sudankdk/bookstore/internal/domain/repo"
	bookdto "github.com/sudankdk/bookstore/internal/dto/BookDTO"
	bookmapper "github.com/sudankdk/bookstore/internal/mapper/BookMapper"
)

type CreateBookUsecase struct {
	repo repo.BookAbstractRepo
}

func NewBookHandler(BookRepo repo.BookAbstractRepo) *CreateBookUsecase {
	return &CreateBookUsecase{repo: BookRepo}
}

func (c *CreateBookUsecase) Execute(dto bookdto.CreateBookDTO) (entity.Book, error) {
	entity := bookmapper.CreateDTOtoEntity(dto)
	return c.repo.Create(entity)
}
