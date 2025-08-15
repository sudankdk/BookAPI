package repo

import "github.com/sudankdk/bookstore/internal/domain/entity"

type UserRepoAbstract interface {
	Register(entity.User) (entity.User, error)
}
