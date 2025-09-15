package jwt

import (
	"fmt"
	"log"
	"maps"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret string

func Init(js string) {
	if js == "" {
		log.Println("[warn] using empty string as a secret key")
	}

	jwtSecret = js
}

func CreateToken(mapClaim map[string]any) (string, error) {
	now := time.Now()

	claims := jwt.MapClaims{
		"iat": now.Unix(),
		"exp": now.Add(time.Hour).Unix(),
	}

	maps.Copy(claims, mapClaim)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return signedString, nil
}

func VerifyToken(tokenString string) (jwt.Claims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token.Claims, nil
}
