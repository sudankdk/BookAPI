package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sudankdk/bookstore/internal/data/sqldb"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Brearer token is missing , u bitch"))
			return
		}

		claims := &sqldb.Claims{}
		token, err := jwt.ParseWithClaims(tokenString[len("Bearer "):], claims, func(token *jwt.Token) (interface{}, error) {
			return sqldb.JwtKey, nil
		})

		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
