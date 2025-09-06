package auth

type User struct {
	Name     string
	Password string
}

func NewUser(name, password string) *User {
	return &User{Name: name, Password: password}
}
