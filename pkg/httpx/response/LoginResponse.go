package response

import "github.com/sudankdk/bookstore/internal/domain/entity"

type LoginResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	User         entity.User `json:"user"`
}
