package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	testUser := NewUser("testuser", "password")

	Init(newMockUserRepository(true), newMockTokenRepository(true))

	t.Run("register success", func(t *testing.T) {
		assert.Nil(t, RegisterUser(testUser))
	})

	t.Run("empty user info", func(t *testing.T) {
		assert.NotNil(t, RegisterUser(NewUser("", "password")))
		assert.NotNil(t, RegisterUser(NewUser("username", "")))
		assert.NotNil(t, RegisterUser(NewUser("", "")))
	})

	t.Run("detect repo error", func(t *testing.T) {
		Init(newMockUserRepository(false), newMockTokenRepository(true))
		assert.NotNil(t, RegisterUser(testUser))
	})
}
