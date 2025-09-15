package user

import "golang.org/x/crypto/bcrypt"

type Repository interface {
	GetUser(username string) (*User, error)
	SaveUser(user *User) error
}

type User struct {
	Name           string
	HashedPassword string
}

func New(username, password string) (*User, error) {
	hashed, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	return &User{
		Name:           username,
		HashedPassword: hashed,
	}, nil
}

func Validate(ur Repository, username, plainPassword string) error {
	foundUser, err := ur.GetUser(username)
	if err != nil {
		return err
	}

	return bcrypt.CompareHashAndPassword([]byte(foundUser.HashedPassword), []byte(plainPassword))
}

// ===========================================================================

func hashPassword(password string) (string, error) {
	hashedB, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedB), nil
}
