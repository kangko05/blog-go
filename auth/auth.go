package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CheckUserExists(username string) bool
	GetUser(username string) (*User, error)
	SaveUser(user *User) error
}

var userRepo UserRepository

// takes dependencies for the package
func Init(ur UserRepository) {
	userRepo = ur
}

func RegisterUser(user *User) error {
	if err := checkUserFormat(user); err != nil {
		return err
	}

	if userRepo.CheckUserExists(user.Name) {
		return fmt.Errorf("user already exists")
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return err
	}

	if err := userRepo.SaveUser(NewUser(user.Name, string(hashedPassword))); err != nil {
		return err
	}

	return nil
}

func Login(user *User) (*Token, error) {
	if err := checkUserFormat(user); err != nil {
		return nil, err
	}

	foundUser, err := userRepo.GetUser(user.Name)
	if err != nil {
		return nil, err
	}

	if err := comparePasswords(foundUser.Password, user.Password); err != nil {
		return nil, err
	}

	return NewToken(), nil
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

	if len(user.Name) < 4 {
		return fmt.Errorf("username too short")
	}

	if len(user.Password) <= 0 {
		return fmt.Errorf("empty password provided")
	}

	if len(user.Password) < 8 {
		return fmt.Errorf("password too short")
	}

	return nil
}
