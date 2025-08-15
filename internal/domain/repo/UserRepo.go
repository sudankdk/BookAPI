package repo

import (
	"github.com/sudankdk/bookstore/internal/domain/entity"
	userdto "github.com/sudankdk/bookstore/internal/dto/UserDTO"
	"github.com/sudankdk/bookstore/pkg/httpx/response"
)

type UserRepoAbstract interface {
	Register(entity.User) (entity.User, error)
	Login(dto userdto.UserLoginDTO) (response.LoginResponse, error)
}
