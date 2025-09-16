package auth

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Repository interface {
	SaveUser(user *User) error
	GetUser(username string) (*User, error)

	GetToken(tokenString string) (string, error)
	SaveToken(tokenString string) error
	DeleteToken(tokenString string) error
}

// ===========================================================================

type memRepo struct {
	conn *sql.DB
}

func connectSqlite() (*memRepo, error) {
	conn, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	mr := &memRepo{conn: conn}

	if err := mr.createTables(); err != nil {
		return nil, err
	}

	return mr, nil
}

func (mr *memRepo) Close() error {
	return mr.conn.Close()
}

func (mr *memRepo) createTables() error {
	queries := []string{
		`CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			hashed_password TEXT NOT NULL
		)`,

		`CREATE TABLE tokens(
			token TEXT PRIMARY KEY
		)`,
	}

	for _, query := range queries {
		_, err := mr.command(query)
		if err != nil {
			return err
		}
	}

	return nil
}

func (mr *memRepo) command(query string, args ...any) (sql.Result, error) {
	tx, err := mr.conn.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	result, err := tx.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func (mr *memRepo) queryRow(q string, args ...any) *sql.Row {
	return mr.conn.QueryRow(q, args...)
}

func (mr *memRepo) SaveUser(user *User) error {
	_, err := mr.command("INSERT INTO users(username, hashed_password) VALUES(?,?)", user.Name, user.HashedPassword)
	if err != nil {
		return err
	}

	return nil
}

func (mr *memRepo) GetUser(username string) (*User, error) {
	var user User

	err := mr.queryRow("SELECT username,hashed_password FROM users WHERE username=?", username).Scan(&user.Name, &user.HashedPassword)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (mr *memRepo) GetToken(tokenString string) (string, error) {
	var foundTokenString string
	err := mr.queryRow("SELECT token FROM tokens WHERE token=?", tokenString).Scan(&foundTokenString)
	if err != nil {
		return "", err
	}

	return foundTokenString, nil
}

func (mr *memRepo) SaveToken(tokenString string) error {
	_, err := mr.command("INSERT OR REPLACE INTO tokens(token) VALUES(?)", tokenString)
	return err
}

func (mr *memRepo) DeleteToken(tokenString string) error {
	_, err := mr.command("DELETE FROM tokens WHERE token=?", tokenString)
	return err
}
