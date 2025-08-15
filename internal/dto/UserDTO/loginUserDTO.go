package userdto

type UserLoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"-"`
}
