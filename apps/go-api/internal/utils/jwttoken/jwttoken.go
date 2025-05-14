package jwttoken

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SecretKey = []byte("secret-key")

func CreateToken(id string, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":       id,
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
