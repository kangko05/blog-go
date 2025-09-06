package auth

import "fmt"

type mockUserRepository struct {
	resultState bool
}

func newMockUserRepository(resultState bool) *mockUserRepository {
	return &mockUserRepository{resultState: resultState}
}

func (mur *mockUserRepository) CheckUserExists(_ string) bool {
	return !mur.resultState
}

func (mur *mockUserRepository) GetUser(username string) (*User, error) {
	if mur.resultState {
		return NewUser(username, "hashed password"), nil
	}

	return nil, fmt.Errorf("failed to get user")
}

func (mur *mockUserRepository) SaveUser(user *User) error {
	if mur.resultState {
		return nil
	}

	return fmt.Errorf("failed to save user")
}

// ===========================================================================

type mockTokenRepository struct {
	resultState bool
}

func newMockTokenRepository(resultState bool) *mockTokenRepository {
	return &mockTokenRepository{resultState: resultState}
}
