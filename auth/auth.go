package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CheckUserExists(username string) bool
	GetUser(username string) (*User, error)
	SaveUser(user *User) error
}

type TokenRepository interface {
	Save(token *jwt.Token) error
}

var (
	userRepo  UserRepository
	tokenRepo TokenRepository
)

func Init(ur UserRepository, tr TokenRepository) {
	userRepo = ur
	tokenRepo = tr
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

func Login(user *User) (*jwt.Token, error) {
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

	token := issueToken(user.Name)

	return token, nil
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
