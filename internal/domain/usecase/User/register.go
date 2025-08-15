package user

import (
	"github.com/sudankdk/bookstore/internal/domain/entity"
	"github.com/sudankdk/bookstore/internal/domain/repo"
	userdto "github.com/sudankdk/bookstore/internal/dto/UserDTO"
	usermapper "github.com/sudankdk/bookstore/internal/mapper/UserMapper"
)

type UserRepoUsecase struct {
	repo repo.UserRepoAbstract
}

func NewUserRepo(userRepo repo.UserRepoAbstract) *UserRepoUsecase {
	return &UserRepoUsecase{repo: userRepo}
}

func (u *UserRepoUsecase) Register(dto userdto.UserDTO) (entity.User, error) {
	entity := usermapper.CreateDTOtoEntity(dto)
	return u.repo.Register(entity)
}
