package auth

type User struct {
	Name           string
	HashedPassword string
}

func NewUser(username, password string) (*User, error) {
	hashed, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	return &User{
		Name:           username,
		HashedPassword: string(hashed),
	}, nil
}
