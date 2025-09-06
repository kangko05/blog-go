package auth

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockUserRepository struct {
	resultState bool
}

func newMockUserRepository(resultState bool) *mockUserRepository {
	return &mockUserRepository{resultState: resultState}
}

func (mur *mockUserRepository) CheckUserExists(_ string) bool {
	return !mur.resultState
}

func (mur *mockUserRepository) GetUser(username string) (*User, error) {
	if mur.resultState {
		return NewUser(username, "hashed password"), nil
	}

	return nil, fmt.Errorf("failed to get user")
}

func (mur *mockUserRepository) SaveUser(user *User) error {
	if mur.resultState {
		return nil
	}

	return fmt.Errorf("failed to save user")
}

func TestRegisterUser(t *testing.T) {
	testUser := NewUser("testuser", "password")

	Init(newMockUserRepository(true))

	t.Run("register success", func(t *testing.T) {
		assert.Nil(t, RegisterUser(testUser))
	})

	t.Run("empty user info", func(t *testing.T) {
		assert.NotNil(t, RegisterUser(NewUser("", "password")))
		assert.NotNil(t, RegisterUser(NewUser("username", "")))
		assert.NotNil(t, RegisterUser(NewUser("", "")))
	})

	t.Run("detect repo error", func(t *testing.T) {
		Init(newMockUserRepository(false))
		assert.NotNil(t, RegisterUser(testUser))
	})
}
