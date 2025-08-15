package user

import (
	userdto "github.com/sudankdk/bookstore/internal/dto/UserDTO"
	"github.com/sudankdk/bookstore/pkg/httpx/response"
)

func (u *UserRepoUsecase) Login(dto userdto.UserLoginDTO) (response.LoginResponse, error) {
	return u.repo.Login(dto)
}
