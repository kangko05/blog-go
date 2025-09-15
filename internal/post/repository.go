package post

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Repository interface {
	SavePost(md *Post) error
	GetPost(id int) (*Post, error)
	UpdatePost(md *Post) error
	DeletePost(id int) error
	ListPosts() ([]*Post, error)

	Close() error // cleanup code
}

// ============================================================================

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
		`CREATE TABLE posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME,
			updated_at DATETIME
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

func (mr *memRepo) SavePost(post *Post) error {
	result, err := mr.command(
		"INSERT INTO posts(title, content, created_at, updated_at) VALUES(?,?,?,?)",
		post.Title, post.Content, post.CreatedAt, post.UpdatedAt,
	)

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	post.Id = int(id)

	return err
}

func (mr *memRepo) GetPost(id int) (*Post, error) {
	var post Post

	err := mr.queryRow(
		"SELECT id,title,content,created_at,updated_at FROM posts WHERE id=?",
		id).Scan(&post.Id, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (mr *memRepo) UpdatePost(post *Post) error {
	_, err := mr.command(
		"UPDATE posts SET title=?, content=?, updated_at=? WHERE id=?",
		post.Title, post.Content, post.UpdatedAt, post.Id,
	)

	return err
}

func (mr *memRepo) DeletePost(id int) error {
	_, err := mr.command("DELETE FROM posts WHERE id=?", id)
	return err
}

func (mr *memRepo) ListPosts() ([]*Post, error) {
	rows, err := mr.conn.Query("SELECT id,title,content,created_at,updated_at FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post

	for rows.Next() {
		var post Post

		err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			log.Println(err)
		}

		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
