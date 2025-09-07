package auth

import "fmt"

type AuthRepository interface {
	CheckUserExists(username string) bool
	GetUser(username string) (*User, error)
	SaveUser(user *User) error
	SaveToken(token *Token) error
	DeleteToken(tokenString string) error
}

const JWT_SECRET string = "secret key here"

var authRepo AuthRepository

func Init(ar AuthRepository) {
	authRepo = ar
}

func RegisterUser(user *User) error {
	if err := checkUserFormat(user); err != nil {
		return err
	}

	if authRepo.CheckUserExists(user.Name) {
		return fmt.Errorf("user already exists")
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return err
	}

	if err := authRepo.SaveUser(NewUser(user.Name, string(hashedPassword))); err != nil {
		return err
	}

	return nil
}

func Login(user *User) (*Token, error) {
	if err := checkUserFormat(user); err != nil {
		return nil, err
	}

	foundUser, err := authRepo.GetUser(user.Name)
	if err != nil {
		return nil, err
	}

	if err := comparePasswords(foundUser.Password, user.Password); err != nil {
		return nil, err
	}

	token, err := issueToken(user.Name)
	if err != nil {
		return nil, err
	}

	if err := authRepo.SaveToken(token); err != nil {
		return nil, err
	}

	return token, nil
}

func Logout(tokenString string) error {
	if err := validateToken(tokenString); err != nil {
		return err
	}

	return authRepo.DeleteToken(tokenString)
}
