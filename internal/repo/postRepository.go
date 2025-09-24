package repo

import (
	"blog-go/internal/post"
	"encoding/json"
	"fmt"
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
			category TEXT NOT NULL,
			created_at DATETIME,
			updated_at DATETIME,
			tags TEXT
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
	tagsJson, err := json.Marshal(post.Tags)
	if err != nil {
		return err
	}

	result, err := pr.db.command(
		"INSERT INTO posts(title, content, category, created_at, updated_at, tags) VALUES(?,?,?,?,?,?)",
		post.Title, post.Content, post.Category, post.CreatedAt, post.UpdatedAt, string(tagsJson),
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
	var tagStr string

	err := pr.db.queryRow(
		"SELECT id,title,content,category,created_at,updated_at,tags FROM posts WHERE id=?",
		id).Scan(&post.Id, &post.Title, &post.Content, &post.Category, &post.CreatedAt, &post.UpdatedAt, &tagStr)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(tagStr), &post.Tags); err != nil {
		return nil, err
	}

	return &post, nil
}

func (pr *PostRepository) UpdatePost(post *post.Post) error {
	tagsJson, err := json.Marshal(post.Tags)
	if err != nil {
		return err
	}

	_, err = pr.db.command(
		"UPDATE posts SET title=?, content=?, updated_at=?, tags=? WHERE id=?",
		post.Title, post.Content, post.UpdatedAt, string(tagsJson), post.Id,
	)

	return err
}

func (pr *PostRepository) DeletePost(id int) error {
	_, err := pr.db.command("DELETE FROM posts WHERE id=?", id)
	return err
}

func (pr *PostRepository) ListPosts() ([]*post.Post, error) {
	rows, err := pr.db.conn.Query("SELECT id,title,content,category,created_at,updated_at,tags FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*post.Post

	for rows.Next() {
		var post post.Post
		var tagsTxt string

		err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.Category, &post.CreatedAt, &post.UpdatedAt, &tagsTxt)
		if err != nil {
			log.Println(err)
		}

		if err := json.Unmarshal([]byte(tagsTxt), &post.Tags); err != nil {
			fmt.Println(err)
		}

		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
