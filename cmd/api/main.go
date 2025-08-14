package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	handlers "github.com/sudankdk/bookstore/internal/adapter/http/Handlers"
)

func main() {

	r := chi.NewRouter()

	r.Mount("/api", handlers.Routes())

	log.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
