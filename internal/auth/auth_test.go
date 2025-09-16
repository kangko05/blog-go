package auth

import (
	"blog-go/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestService(t *testing.T) (*Service, func()) {
	cfg := &config.Config{
		JwtSecret: "test-secret-key",
		DbPath:    ":memory:",
	}

	service, err := NewService(cfg, nil)
	assert.NoError(t, err)

	cleanup := func() {}

	return service, cleanup
}

func TestRegister(t *testing.T) {
	service, cleanup := setupTestService(t)
	defer cleanup()

	t.Run("register success", func(t *testing.T) {
		err := service.Register("testuser", "testpass")
		assert.Nil(t, err)
	})

	t.Run("register fail - short username|password", func(t *testing.T) {
		assert.NotNil(t, service.Register("a", "testpass")) // short username
		assert.NotNil(t, service.Register("", "testpass"))  // empty username

		assert.NotNil(t, service.Register("username", "t")) // short password
		assert.NotNil(t, service.Register("username", ""))  // empty password
	})

	t.Run("duplicate user", func(t *testing.T) {
		username, password := "dupuser", "duppassword"
		assert.Nil(t, service.Register(username, password))    // no error for first register
		assert.NotNil(t, service.Register(username, password)) // error for next registers
	})
}

func TestLogin(t *testing.T) {
	service, cleanup := setupTestService(t)
	defer cleanup()

	username, password := "testlogin", "testpass"
	assert.Nil(t, service.Register(username, password))

	t.Run("login success", func(t *testing.T) {
		tokenString, err := service.Login(username, password)
		assert.Nil(t, err)
		assert.Greater(t, len(tokenString), 0)
	})

	t.Run("login fail", func(t *testing.T) {
		_, err := service.Login(username, "invalid password")
		assert.NotNil(t, err)
	})

	t.Run("invalid user", func(t *testing.T) {
		_, err := service.Login("invalidUser", password)
		assert.NotNil(t, err)
	})
}

func TestLogout(t *testing.T) {
	service, cleanup := setupTestService(t)
	defer cleanup()

	username, password := "testlogout", "testpass"
	assert.Nil(t, service.Register(username, password))

	token, err := service.Login(username, password)
	assert.Nil(t, err)

	t.Run("logout success", func(t *testing.T) {
		assert.Nil(t, service.Logout(token)) // no error
	})

	t.Run("logout invalid token", func(t *testing.T) {
		assert.NotNil(t, service.Logout("invalid-token"))
	})
}

func TestVerifyToken(t *testing.T) {
	service, cleanup := setupTestService(t)
	defer cleanup()

	username, password := "testverify", "testpass"
	assert.Nil(t, service.Register(username, password))

	t.Run("verify valid token", func(t *testing.T) {
		token, err := service.Login(username, password)
		assert.Nil(t, err)

		err = service.VerifyToken(token)
		assert.Nil(t, err)
	})

	t.Run("verify invalid token format", func(t *testing.T) {
		err := service.VerifyToken("invalid.token.format")
		assert.NotNil(t, err)

		err = service.VerifyToken("not-a-token")
		assert.NotNil(t, err)

		err = service.VerifyToken("")
		assert.NotNil(t, err)
	})

	t.Run("verify revoked token", func(t *testing.T) {
		token, err := service.Login(username, password)
		assert.Nil(t, err)

		err = service.Logout(token)
		assert.Nil(t, err)

		err = service.VerifyToken(token)
		assert.NotNil(t, err)
	})

	t.Run("verify token with wrong secret", func(t *testing.T) {
		wrongCfg := &config.Config{
			JwtSecret: "wrong-secret",
			DbPath:    ":memory:",
		}
		wrongService, _ := NewService(wrongCfg, nil)

		wrongService.Register(username, password)
		wrongToken, err := wrongService.Login(username, password)
		assert.Nil(t, err)

		err = service.VerifyToken(wrongToken)
		assert.NotNil(t, err)
	})
}

func TestTokenLifecycle(t *testing.T) {
	service, cleanup := setupTestService(t)
	defer cleanup()

	username, password := "lifecycle", "testpass"
	assert.Nil(t, service.Register(username, password))

	token1, err := service.Login(username, password)
	assert.Nil(t, err)
	assert.NotEmpty(t, token1)

	assert.Nil(t, service.VerifyToken(token1))

	token2, err := service.Login(username, password)
	assert.Nil(t, err)
	assert.NotEmpty(t, token2)

	assert.Nil(t, service.VerifyToken(token1))
	assert.Nil(t, service.VerifyToken(token2))

	assert.Nil(t, service.Logout(token1))

	assert.NotNil(t, service.VerifyToken(token1))
	assert.Nil(t, service.VerifyToken(token2))

	assert.Nil(t, service.Logout(token2))
	assert.NotNil(t, service.VerifyToken(token2))
}
