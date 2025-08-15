package sqldb

import (
	"time"

	"github.com/google/uuid"
	"github.com/sudankdk/bookstore/internal/domain/entity"
)

func (r *SQLBookRepo) Register(u entity.User) (entity.User, error) {
	r.DB.Exec(`
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    username TEXT NOT NULL CHECK(length(username) >= 3),
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
)`)

	if u.Id == "" {
		u.Id = uuid.NewString()
	}
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	query := `INSERT INTO users (id, username, email, password, created_at, updated_at)
			  VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := r.DB.Exec(query, u.Id, u.Username, u.Email, u.Password, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return entity.User{}, err
	}
	return u, nil

}
