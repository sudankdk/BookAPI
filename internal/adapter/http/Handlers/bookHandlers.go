package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	book "github.com/sudankdk/bookstore/internal/domain/usecase/Book"
	bookdto "github.com/sudankdk/bookstore/internal/dto/BookDTO"
	"github.com/sudankdk/bookstore/pkg/httpx/response"
)

type BookHandler struct {
	bookService *book.CreateBookUsecase
}

func NewBookHandler(s *book.CreateBookUsecase) *BookHandler {
	return &BookHandler{bookService: s}
}

func (h *BookHandler) Create(w http.ResponseWriter, r *http.Request) {
	var book bookdto.CreateBookDTO
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	bookEntity, err := h.bookService.Execute(book)
	if err != nil {
		response.WriteJSON(w, 400, response.APIResponse{
			Error:   err.Error(),
			Success: false,
		})
		return
	}

	response.WriteJSON(w, 201, response.APIResponse{
		Success: true,
		Data:    bookEntity,
	})
}

func (h *BookHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	book, err := h.bookService.GetById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	response.WriteJSON(w, 200, response.APIResponse{
		Success: true,
		Data:    book,
	})
}

func (h *BookHandler) List(w http.ResponseWriter, r *http.Request) {
	books, err := h.bookService.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	response.WriteJSON(w, 200, response.APIResponse{
		Success: true,
		Data:    books,
	})
}
