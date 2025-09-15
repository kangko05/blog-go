package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	MinUsernameLen = 4
	MinPasswordLen = 8
)

func checkUserFormat(username, password string) error {
	if len(username) < MinUsernameLen {
		return fmt.Errorf("username must be longer than %d", MinUsernameLen)
	}

	if len(password) < MinPasswordLen {
		return fmt.Errorf("password must be longer than %d", MinPasswordLen)
	}

	return nil
}

func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func checkPassword(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

// ===========================================================================

func createToken(username string) (string, error) {
	// create token
	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat":      now.Unix(),
		"exp":      now.Add(time.Hour).Unix(),
		"username": username,
	})

	signedString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	if err := repo.SaveToken(signedString); err != nil {
		return "", err
	}

	return signedString, nil
}

func revokeToken(tokenString string) error {
	return repo.DeleteToken(tokenString)
}
