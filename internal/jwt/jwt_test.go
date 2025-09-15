package jwt

import (
	"os"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func extractJwtClaims(tokenString string) (jwt.Claims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	return token.Claims, nil
}

func TestMain(m *testing.M) {
	Init("secret key comes here")

	exitCode := m.Run()

	//

	os.Exit(exitCode)
}

func TestCreateToken(t *testing.T) {
	t.Run("empty claims", func(t *testing.T) {
		tokenString, err := CreateToken(nil)
		assert.Nil(t, err)
		assert.NotEqual(t, len(tokenString), 0)
	})

	t.Run("with claims", func(t *testing.T) {
		tkey, tval := "test", "test value"

		tclaim := map[string]any{
			tkey: tval,
		}

		tokenString, err := CreateToken(tclaim)
		assert.Nil(t, err)
		assert.NotEqual(t, len(tokenString), 0)

		claims, err := extractJwtClaims(tokenString)
		assert.Nil(t, err)

		mapClaim, ok := claims.(jwt.MapClaims)
		assert.True(t, ok)
		assert.Equal(t, mapClaim[tkey], tval)
	})
}

func TestVerifyToken(t *testing.T) {
	tokenString, err := CreateToken(map[string]any{"test": "test val"})

	assert.Nil(t, err)

	t.Run("verify success", func(t *testing.T) {
		claims, err := VerifyToken(tokenString)
		assert.Nil(t, err)
		assert.NotNil(t, claims)

		mc, ok := claims.(jwt.MapClaims)
		assert.True(t, ok)
		assert.Equal(t, mc["test"], "test val")
	})

	t.Run("verify fail (invalid token)", func(t *testing.T) {
		claims, err := VerifyToken("invalid token")
		assert.NotNil(t, err)
		assert.Nil(t, claims)
	})
}
