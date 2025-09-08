package database

import (
	"blog-go/auth"
	"database/sql"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository() (*AuthRepository, error) {
	return &AuthRepository{db: globalDb}, nil
}

func (ar *AuthRepository) DeleteToken(tokenString string) error {
	return command("DELETE FROM tokens WHERE token_string=?", tokenString)
}

func (ar *AuthRepository) SaveToken(token *auth.Token) error {
	return command(
		"INSERT INTO tokens(token_string,issued_at,expires_at) VALUES(?,?,?)",
		token.String, token.IssuedAt, token.ExpiresAt,
	)
}

func (ar *AuthRepository) SaveUser(user *auth.User) error {
	return command(
		"INSERT INTO users(username,password) VALUES(?,?)",
		user.Name, user.Password,
	)
}

func (ar *AuthRepository) GetUser(username string) (*auth.User, error) {
	var name, password string

	query := "SELECT username,password FROM users WHERE username=?"
	err := queryRow(query, username).Scan(&name, &password)
	if err != nil {
		return nil, err
	}

	return auth.NewUser(name, password), nil
}

func (ar *AuthRepository) CheckUserExists(username string) (bool, error) {
	var cnt int

	err := queryRow("SELECT COUNT(*) FROM users WHERE username=?", username).Scan(&cnt)
	if err != nil {
		return false, err
	}

	return cnt > 0, nil
}
