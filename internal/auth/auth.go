package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func register(repo Repository, username, password string) error {
	if err := checkUserFormat(username, password); err != nil {
		return err
	}

	user, err := NewUser(username, password)
	if err != nil {
		return err
	}

	return repo.SaveUser(user)
}

func login(repo Repository, jwtSecret, username, password string) (string, error) {
	foundUser, err := repo.GetUser(username)
	if err != nil {
		return "", err
	}

	// check password
	if err := checkPassword(foundUser.HashedPassword, password); err != nil {
		return "", err
	}

	return createToken(repo, jwtSecret, username)
}

func logout(repo Repository, tokenString string) error {
	_, err := repo.GetToken(tokenString)
	if err != nil {
		return err
	}

	return revokeToken(repo, tokenString)
}

func verifyToken(repo Repository, jwtSecret, tokenString string) error {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return err // parse error
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	// check db
	_, err = repo.GetToken(tokenString)
	if err != nil {
		return err
	}

	return nil
}
