package handlers

import (
	"github.com/go-chi/chi/v5"
	inmemorydb "github.com/sudankdk/bookstore/internal/data/Inmemorydb"
	book "github.com/sudankdk/bookstore/internal/domain/usecase/Book"
)

func Routes() *chi.Mux {
	repo := inmemorydb.NewInMemoryBookRepo()
	service := book.NewBookHandler(repo)
	handler := NewBookHandler(service)

	r := chi.NewRouter()
	r.Post("/books", handler.Create)
	return r
}
