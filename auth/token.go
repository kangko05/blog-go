package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	String    string
	IssuedAt  int64
	ExpiresAt int64
}

// =============================================================================

func issueToken(username string) (*Token, error) {
	now := time.Now()
	iat := now.Unix()
	exp := now.Add(time.Hour).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"iat":      now.Unix(),
		"exp":      exp,
		"nonce":    now.UnixNano(),
	})

	signedString, err := token.SignedString([]byte(JWT_SECRET))
	if err != nil {
		return nil, err
	}

	return &Token{
		String:    signedString,
		IssuedAt:  iat,
		ExpiresAt: exp,
	}, nil
}

func validateToken(tokenString string) error {
	parsedToken, err := parseToken(tokenString)
	if err != nil {
		return err
	}

	if !parsedToken.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func parseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return []byte(JWT_SECRET), nil
	})
}
