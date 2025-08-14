package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/go-chi/chi/v5"
	"github.com/sudankdk/bookstore/internal/data/sqldb"
	book "github.com/sudankdk/bookstore/internal/domain/usecase/Book"
)

func Routes() *chi.Mux {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env variables")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	// repo := inmemorydb.NewInMemoryBookRepo(db)
	repo := sqldb.NewSqlBookRepo(db)
	service := book.NewBookHandler(repo)
	handler := NewBookHandler(service)

	r := chi.NewRouter()
	r.Post("/books", handler.Create)
	return r
}
