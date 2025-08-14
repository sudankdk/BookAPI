package bookmapper

import (
	"github.com/sudankdk/bookstore/internal/domain/entity"
	bookdto "github.com/sudankdk/bookstore/internal/dto/BookDTO"
)

func CreateDTOtoEntity(dto bookdto.CreateBookDTO) entity.Book {
	return entity.Book{
		Name:        dto.Name,
		Author:      dto.Author,
		PriceCents:  dto.PriceCents,
		ISBN:        dto.ISBN,
		PublishedAt: dto.PublishedAt,
	}
}
