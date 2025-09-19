package repo

import (
	"blog-go/internal/config"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	conn *sql.DB
}

func ConnectDatabase(cfg *config.Config) (*Database, error) {
	conn, err := sql.Open("sqlite3", cfg.DbPath)
	if err != nil {
		return nil, err
	}

	return &Database{conn: conn}, nil
}

func (db *Database) Close() error {
	return db.conn.Close()
}

func (db *Database) command(query string, args ...any) (sql.Result, error) {
	tx, err := db.conn.Begin()
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

func (db *Database) queryRow(q string, args ...any) *sql.Row {
	return db.conn.QueryRow(q, args...)
}
