package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func issueToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Hour).Unix(),
	})

	if err := tokenRepo.Save(token); err != nil {
		return "", err
	}

	signedString, err := token.SignedString("secret key here")
	if err != nil {
		return "", err
	}

	return signedString, nil
}
