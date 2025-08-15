package userdto

type UserDTO struct {
	Username string `json:"username" validate:"required,max=200,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"-"`
}
