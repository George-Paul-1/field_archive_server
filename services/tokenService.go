package services

import (
	"errors"
	"field_archive/server/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(username string, cfg config.Config) (string, error) {
	if username == "" {
		return "", errors.New("username cannot be blank for token creation")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Minute * 10).Unix(), // 10 minute expiry time
		})
	secret := []byte(cfg.JwtSecret)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
