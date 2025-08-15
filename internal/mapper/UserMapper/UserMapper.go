package usermapper

import (
	"github.com/sudankdk/bookstore/internal/domain/entity"
	userdto "github.com/sudankdk/bookstore/internal/dto/UserDTO"
)

func CreateDTOtoEntity(dto userdto.UserDTO) entity.User {
	return entity.User{
		Username: dto.Username,
		Email:    dto.Email,
		Password: dto.Password,
	}
}
