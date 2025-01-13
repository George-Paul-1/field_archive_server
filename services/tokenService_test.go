package services

import (
	"errors"
	"field_archive/server/internal/config"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestCreateToken(t *testing.T) {
	Mockcfg := config.Config{JwtSecret: "Test"}
	hmacSampleSecret := []byte("Test")
	u := "Mockuser"
	token, err := CreateToken(u, Mockcfg)
	if err != nil {
		t.Errorf("Error in token creation, %v", err)
	}
	check, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected Signing method")
		}
		return hmacSampleSecret, nil
	})
	if err != nil {
		t.Errorf("%v", err)
	}
	var claims jwt.MapClaims
	if c, ok := check.Claims.(jwt.MapClaims); ok {
		claims = c
	} else {
		t.Errorf("Error extracting claims %v", claims["username"])
	}
	if u != claims["username"] {
		t.Errorf("Usernames do not match %v should be %v", claims["username"], u)
	}
}
