package utils

import (
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func ExtractUserIDFromAuthHeader(secret, headerval, queryval string) (string, error) {
	raw := headerval
	if raw == "" {
		raw = queryval
	}

	if raw == "" {
		return "", errors.New("missing auth")
	}
	if after, ok := strings.CutPrefix(raw, "Bearer "); ok {
		raw = after
	}
	token, err := jwt.Parse(raw, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}
	userID, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("userID not found in token")
	}
	return userID, nil
}
