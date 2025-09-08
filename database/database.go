package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var globalDb *sql.DB

func Init(connectionPath string) error {
	conn, err := sql.Open("sqlite3", connectionPath)
	if err != nil {
		return err
	}

	globalDb = conn

	if err := createTables(); err != nil {
		return err
	}

	return nil
}

func Close() error {
	if globalDb != nil {
		return globalDb.Close()
	}

	return nil
}

func command(query string, args ...any) error {
	tx, err := globalDb.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(query, args...)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func queryRow(query string, args ...any) *sql.Row {
	return globalDb.QueryRow(query, args...)
}

func query(query string, args ...any) (*sql.Rows, error) {
	return globalDb.Query(query, args...)
}

// create necessary tables ====================================================
func createTables() error {
	if globalDb == nil {
		return fmt.Errorf("attempt to create table before db initialization")
	}

	queries := []string{

		`CREATE TABLE IF NOT EXISTS users(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL
		)`,

		`CREATE TABLE IF NOT EXISTS tokens(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			token_string VARCHAR(255) UNIQUE NOT NULL,
			user_id INTEGER,
			issued_at INTEGER NOT NULL,
			expires_at INTEGER NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users (id)
		)`,

		`CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)`,
		`CREATE INDEX IF NOT EXISTS idx_tokens_expires ON tokens(expires_at)`,
	}

	for _, query := range queries {
		if err := command(query); err != nil {
			return err
		}
	}

	return nil
}
