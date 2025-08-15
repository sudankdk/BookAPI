package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"

	"github.com/go-chi/chi/v5"
	"github.com/sudankdk/bookstore/internal/data/sqldb"
	book "github.com/sudankdk/bookstore/internal/domain/usecase/Book"
	user "github.com/sudankdk/bookstore/internal/domain/usecase/User"
)

func Routes() *chi.Mux {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
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
	handler := NewBookHandler(service, rdb)

	serive2 := user.NewUserRepo(repo)
	handler2 := NewUserHandler(serive2)

	r := chi.NewRouter()
	r.Post("/books", handler.Create)
	r.Get("/books/{id}", handler.Get)
	r.Get("/books", handler.List)
	r.Delete("/books/{id}", handler.Delete)
	r.Patch("/books/{id}", handler.Update)

	r.Post("/register", handler2.Register)

	return r
}
