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

func VerifyToken(tokenString string, cfg config.Config) (string, error) {
	secret := []byte(cfg.JwtSecret)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected Signing method")
		}
		return secret, nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("error extracting claims")
	}
	username, ok := claims["username"].(string)
	if !ok {
		return "", errors.New("error extracting username")
	}
	return username, nil
}
