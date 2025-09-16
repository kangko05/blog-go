package repo

import "blog-go/internal/auth"

type AuthRepository struct {
	db *Database
}

func NewAuthRepository(db *Database) (*AuthRepository, error) {
	ar := &AuthRepository{db: db}

	if err := ar.createTables(); err != nil {
		return nil, err
	}

	return ar, nil
}

func (ar *AuthRepository) createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			hashed_password TEXT NOT NULL
		)`,

		`CREATE TABLE IF NOT EXISTS tokens(
			token TEXT PRIMARY KEY
		)`,
	}

	for _, query := range queries {
		_, err := ar.db.command(query)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ar *AuthRepository) SaveUser(user *auth.User) error {
	_, err := ar.db.command("INSERT INTO users(username, hashed_password) VALUES(?,?)", user.Name, user.HashedPassword)
	if err != nil {
		return err
	}

	return nil
}

func (ar *AuthRepository) GetUser(username string) (*auth.User, error) {
	var user auth.User

	err := ar.db.queryRow("SELECT username,hashed_password FROM users WHERE username=?", username).Scan(&user.Name, &user.HashedPassword)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ar *AuthRepository) GetToken(tokenString string) (string, error) {
	var foundTokenString string
	err := ar.db.queryRow("SELECT token FROM tokens WHERE token=?", tokenString).Scan(&foundTokenString)
	if err != nil {
		return "", err
	}

	return foundTokenString, nil
}

func (ar *AuthRepository) SaveToken(tokenString string) error {
	_, err := ar.db.command("INSERT OR REPLACE INTO tokens(token) VALUES(?)", tokenString)
	return err
}

func (ar *AuthRepository) DeleteToken(tokenString string) error {
	_, err := ar.db.command("DELETE FROM tokens WHERE token=?", tokenString)
	return err
}
