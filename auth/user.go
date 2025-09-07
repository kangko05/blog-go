package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	MinUsernameLength     = 4
	MinUserPasswordLength = 8
)

type User struct {
	Name     string `json:"name"`
	Password string `json:"-"`
}

func NewUser(name, password string) *User {
	return &User{Name: name, Password: password}
}

// ============================================================================

// bcrypt wrapper
func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func comparePasswords(encrypted, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(encrypted), []byte(plain))
}

func checkUserFormat(user *User) error {
	if len(user.Name) <= 0 {
		return fmt.Errorf("empty username provided")
	}

	if len(user.Name) < MinUsernameLength {
		return fmt.Errorf("username too short")
	}

	if len(user.Password) <= 0 {
		return fmt.Errorf("empty password provided")
	}

	if len(user.Password) < MinUserPasswordLength {
		return fmt.Errorf("password too short")
	}

	return nil
}
