package auth

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	if err := Init(nil, "jwt secret here"); err != nil {
		panic(err)
	}

	exitCode := m.Run()
	defer repo.Close()
	os.Exit(exitCode)
}

func TestRegister(t *testing.T) {
	t.Run("register success", func(t *testing.T) {
		tusername, tpass := "testuser", "testpass"

		assert.Nil(t, Register(tusername, tpass))

		user, err := repo.GetUser(tusername)
		assert.Nil(t, err)
		assert.Equal(t, user.Name, tusername)
	})

	t.Run("register fail - short username|password", func(t *testing.T) {
		assert.NotNil(t, Register("a", "testpass")) // short username
		assert.NotNil(t, Register("", "testpass"))  // empty username

		assert.NotNil(t, Register("username", "t")) // short password
		assert.NotNil(t, Register("username", ""))  // empty password
	})

	t.Run("duplicate user", func(t *testing.T) {
		username, password := "dupuser", "duppassword"
		assert.Nil(t, Register(username, password))    // no error for first register
		assert.NotNil(t, Register(username, password)) // error for next registers
	})
}

func TestLogin(t *testing.T) {
	username, password := "testlogin", "testpass"
	assert.Nil(t, Register(username, password))

	t.Run("login success", func(t *testing.T) {
		tokenString, err := Login(username, password)
		assert.Nil(t, err)
		assert.Greater(t, len(tokenString), 0)
	})

	t.Run("login fail", func(t *testing.T) {
		_, err := Login(username, "invalid password")
		assert.NotNil(t, err)
	})

	t.Run("invalid user", func(t *testing.T) {
		_, err := Login("invalidUser", password)
		assert.NotNil(t, err)
	})
}

func TestLogout(t *testing.T) {
	username, password := "testlogout", "testpass"
	assert.Nil(t, Register(username, password))

	token, err := Login(username, password)
	assert.Nil(t, err)

	t.Run("logout sucess", func(t *testing.T) {
		assert.Nil(t, Logout(token)) // no error
	})
}
