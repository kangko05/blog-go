package repo

import (
	"blog-go/internal/post"
	"log"
)

type PostRepository struct{ db *Database }

func NewPostRepository(db *Database) (*PostRepository, error) {
	pr := &PostRepository{db: db}

	if err := pr.createTables(); err != nil {
		return nil, err
	}

	return pr, nil
}

func (pr *PostRepository) createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME,
			updated_at DATETIME
		)`,
	}

	for _, query := range queries {
		_, err := pr.db.command(query)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pr *PostRepository) SavePost(post *post.Post) error {
	result, err := pr.db.command(
		"INSERT INTO posts(title, content, created_at, updated_at) VALUES(?,?,?,?)",
		post.Title, post.Content, post.CreatedAt, post.UpdatedAt,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	post.Id = int(id)

	return err
}

func (pr *PostRepository) GetPost(id int) (*post.Post, error) {
	var post post.Post

	err := pr.db.queryRow(
		"SELECT id,title,content,created_at,updated_at FROM posts WHERE id=?",
		id).Scan(&post.Id, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (pr *PostRepository) UpdatePost(post *post.Post) error {
	_, err := pr.db.command(
		"UPDATE posts SET title=?, content=?, updated_at=? WHERE id=?",
		post.Title, post.Content, post.UpdatedAt, post.Id,
	)

	return err
}

func (pr *PostRepository) DeletePost(id int) error {
	_, err := pr.db.command("DELETE FROM posts WHERE id=?", id)
	return err
}

func (pr *PostRepository) ListPosts() ([]*post.Post, error) {
	rows, err := pr.db.conn.Query("SELECT id,title,content,created_at,updated_at FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*post.Post

	for rows.Next() {
		var post post.Post

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
