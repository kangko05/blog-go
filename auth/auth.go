package auth

import "fmt"

const JWT_SECRET string = "secret key here"

var authRepo Repository

func Init(ar Repository) {
	authRepo = ar
}

func RegisterUser(user *User) error {
	if err := checkUserFormat(user); err != nil {
		return err
	}

	exists, err := authRepo.CheckUserExists(user.Name)
	if err != nil {
		return err
	}
	if exists {
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
