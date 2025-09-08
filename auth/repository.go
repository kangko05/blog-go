package auth

type Repository interface {
	CheckUserExists(username string) (bool, error)
	GetUser(username string) (*User, error)
	SaveUser(user *User) error
	SaveToken(token *Token) error
	DeleteToken(tokenString string) error
}
