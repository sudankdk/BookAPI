package handlers

import (
	"encoding/json"
	"net/http"

	user "github.com/sudankdk/bookstore/internal/domain/usecase/User"
	userdto "github.com/sudankdk/bookstore/internal/dto/UserDTO"
	"github.com/sudankdk/bookstore/pkg/httpx/response"
)

type UserHandler struct {
	userService *user.UserRepoUsecase
}

func NewUserHandler(user *user.UserRepoUsecase) *UserHandler {
	return &UserHandler{userService: user}
}

func (u *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user userdto.UserDTO
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	userEntity, err := u.userService.Register(user)
	if err != nil {
		response.WriteJSON(w, 400, response.APIResponse{
			Error:   err.Error(),
			Success: false,
		})
		return
	}

	response.WriteJSON(w, 201, response.APIResponse{
		Success: true,
		Data:    userEntity,
	})
}
