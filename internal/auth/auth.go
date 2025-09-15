package auth

import (
	"fmt"
	"log"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

var (
	repo      Repository
	jwtSecret string
)

func Init(r Repository, js string) error {
	repo = r
	jwtSecret = js

	if r == nil {
		log.Println("[warn] got nil repository, using in-memory repo")
		memRepo, err := connectSqlite()
		if err != nil {
			return err
		}

		repo = memRepo
	}

	return nil
}

// ===========================================================================

func Register(username, password string) error {
	if err := checkUserFormat(username, password); err != nil {
		return err
	}

	user, err := NewUser(username, password)
	if err != nil {
		return err
	}

	return repo.SaveUser(user)
}

func Login(username, password string) (string, error) {
	foundUser, err := repo.GetUser(username)
	if err != nil {
		return "", err
	}

	// check password
	if err := checkPassword(foundUser.HashedPassword, password); err != nil {
		return "", err
	}

	return createToken(username)
}

func Logout(tokenString string) error {
	return revokeToken(tokenString)
}

func VerifyToken(tokenString string) error {
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

func TestVerifyToken(t *testing.T) {
	username, password := "testverify", "testpass"
	assert.Nil(t, Register(username, password))

	t.Run("verify valid token", func(t *testing.T) {
		token, err := Login(username, password)
		assert.Nil(t, err)

		err = VerifyToken(token)
		assert.Nil(t, err)
	})

	t.Run("verify invalid token format", func(t *testing.T) {
		err := VerifyToken("invalid.token.format")
		assert.NotNil(t, err)

		err = VerifyToken("not-a-token")
		assert.NotNil(t, err)

		err = VerifyToken("")
		assert.NotNil(t, err)
	})

	t.Run("verify revoked token", func(t *testing.T) {
		token, err := Login(username, password)
		assert.Nil(t, err)

		err = Logout(token)
		assert.Nil(t, err)

		err = VerifyToken(token)
		assert.NotNil(t, err)
	})

	t.Run("verify token with wrong secret", func(t *testing.T) {
		wrongToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VybmFtZSI6InRlc3QifQ.wrong_signature"

		err := VerifyToken(wrongToken)
		assert.NotNil(t, err)
	})
}
