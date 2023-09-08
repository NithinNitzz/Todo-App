package jwtutil

import (
	"github.com/dgrijalva/jwt-go"
	_ "github.com/labstack/echo/v4"
	"time"
)

const secretKey = "todoapp"

// GenerateToken generates a JWT token for the given username
func GenerateToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
