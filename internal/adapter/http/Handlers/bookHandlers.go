package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	book "github.com/sudankdk/bookstore/internal/domain/usecase/Book"
	bookdto "github.com/sudankdk/bookstore/internal/dto/BookDTO"
	"github.com/sudankdk/bookstore/pkg/httpx/response"
)

type BookHandler struct {
	bookService *book.CreateBookUsecase
	Redis       *redis.Client
}

func NewBookHandler(s *book.CreateBookUsecase, redisclient *redis.Client) *BookHandler {
	return &BookHandler{bookService: s, Redis: redisclient}
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
	ctx := context.Background()
	//suru ma cachma xa ki xain check garnre
	cached, err := h.Redis.Get(ctx, "book"+id).Result()
	if err == nil {
		//cache hit vayo
		response.WriteJSON(w, 200, response.APIResponse{
			Success: true,
			Data:    cached,
		})
	}
	book, err := h.bookService.GetById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	h.Redis.Set(ctx, "book"+id, book, 10*time.Minute)
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

func (h *BookHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing book id", http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	var data map[string]any
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	book, err := h.bookService.Update(id, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if h.Redis != nil {
		h.Redis.Del(ctx, "book"+id)
	}

	response.WriteJSON(w, 200, response.APIResponse{
		Success: true,
		Data:    book,
	})

}

func (h *BookHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing book id", http.StatusBadRequest)
		return
	}
	if err := h.bookService.Delete(id); err != nil {
		response.WriteJSON(w, 200, response.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	ctx := context.Background()
	if h.Redis != nil {
		h.Redis.Del(ctx, "book"+id)
	}
	response.WriteJSON(w, 200, response.APIResponse{
		Success: true,
		Data:    nil,
	})

}
