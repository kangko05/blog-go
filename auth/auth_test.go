package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	testUser := NewUser("testuser", "password")

	mar, err := newMockAuthRepository(true, testUser)
	assert.Nil(t, err)
	Init(mar)

	t.Run("register success", func(t *testing.T) {
		assert.Nil(t, RegisterUser(NewUser("testUser2", "password")))
	})

	t.Run("empty user info", func(t *testing.T) {
		assert.NotNil(t, RegisterUser(NewUser("", "password")))
		assert.NotNil(t, RegisterUser(NewUser("username", "")))
		assert.NotNil(t, RegisterUser(NewUser("", "")))
	})

	t.Run("detect repo error", func(t *testing.T) {
		mar, err := newMockAuthRepository(false, NewUser("ThisHasToFail", "ThisHasToFail"))
		assert.Nil(t, err)
		Init(mar)

		assert.NotNil(t, RegisterUser(testUser))
	})
}

func TestLogin(t *testing.T) {
	testUser := NewUser("testuser", "password")

	mar, err := newMockAuthRepository(true, testUser)
	assert.Nil(t, err)
	Init(mar)

	t.Run("login success", func(t *testing.T) {
		token, err := Login(testUser)
		assert.Nil(t, err)
		assert.Nil(t, validateToken(token.String))
	})

	t.Run("login fail", func(t *testing.T) {
		// empty user info provided
		token, err := Login(NewUser("", "password"))
		assert.NotNil(t, err)
		assert.Nil(t, token)

		token, err = Login(NewUser("user", ""))
		assert.NotNil(t, err)
		assert.Nil(t, token)

		// wrong password
		token, err = Login(NewUser("testuser", "password!"))
		assert.NotNil(t, err)
		assert.Nil(t, token)

		// user not registered
		token, err = Login(NewUser("newUser", "password"))
		assert.NotNil(t, err)
		assert.Nil(t, token)
	})

}
