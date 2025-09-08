package auth

import "fmt"

type mockAuthRepository struct {
	resultState bool
	user        *User
}

func newMockAuthRepository(resultState bool, testUser *User) (*mockAuthRepository, error) {
	hashed, err := hashPassword(testUser.Password)
	if err != nil {
		return nil, err
	}

	return &mockAuthRepository{
		resultState: resultState,
		user:        NewUser(testUser.Name, string(hashed)),
	}, nil
}

func (mar *mockAuthRepository) CheckUserExists(username string) (bool, error) {
	if mar.resultState {
		return mar.user.Name == username, nil
	} else {
		return mar.user.Name == username, fmt.Errorf("failed to check user")
	}
}

func (mar *mockAuthRepository) GetUser(username string) (*User, error) {
	if mar.user.Name == username {
		return mar.user, nil
	}

	return nil, fmt.Errorf("failed to find user")
}

func (mar *mockAuthRepository) SaveUser(user *User) error {
	if mar.resultState {
		return nil
	}

	return fmt.Errorf("failed to save user")
}

func (mar *mockAuthRepository) SaveToken(token *Token) error {
	if mar.resultState {
		return nil
	}

	return fmt.Errorf("failed to save token")
}

func (mar *mockAuthRepository) DeleteToken(tokenString string) error {
	if mar.resultState {
		return nil
	}

	return fmt.Errorf("failed to delete token")
}
