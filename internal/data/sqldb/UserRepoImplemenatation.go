package sqldb

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sudankdk/bookstore/internal/domain/entity"
	userdto "github.com/sudankdk/bookstore/internal/dto/UserDTO"
	"github.com/sudankdk/bookstore/pkg/httpx/response"
	"github.com/sudankdk/bookstore/pkg/utils"
)

var JwtKey = []byte("my_secret_key")

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func (r *SQLBookRepo) Register(u entity.User) (entity.User, error) {

	if u.Id == "" {
		u.Id = uuid.NewString()
	}
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	passwordHashed, _ := utils.HashPassword(u.Password)
	u.Password = passwordHashed
	query := `INSERT INTO users (id, username, email, password, created_at, updated_at)
			  VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := r.DB.Exec(query, u.Id, u.Username, u.Email, u.Password, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return entity.User{}, err
	}
	return u, nil

}

func (r *SQLBookRepo) Login(dto userdto.UserLoginDTO) (response.LoginResponse, error) {
	//check if user exist
	user, err := r.get_user_by_email(dto.Email)
	if err != nil {
		return response.LoginResponse{}, err
	}
	if !utils.CheckPW(user.Password, dto.Password) {
		return response.LoginResponse{}, fmt.Errorf("invalid credentials")
	}
	expirationTime := time.Now().Add(45 * time.Minute)
	claims := &Claims{
		UserID: user.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	// Generate access token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(JwtKey)
	if err != nil {
		return response.LoginResponse{}, err
	}

	// Generate refresh token (valid for longer period, e.g., 7 days)
	refreshExpiration := time.Now().Add(7 * 24 * time.Hour)
	refreshClaims := &Claims{
		UserID: user.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiration),
		},
	}
	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := refreshTokenObj.SignedString(JwtKey)
	if err != nil {
		return response.LoginResponse{}, err
	}

	return response.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (r *SQLBookRepo) get_user_by_email(email string) (entity.User, error) {
	query := `SELECT id, username, email, password, created_at, updated_at FROM users WHERE email=$1`
	var user entity.User
	err := r.DB.QueryRow(query, email).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}
