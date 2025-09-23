package post

import (
	"slices"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type Category string

const (
	ALL      Category = "all"
	NOTES    Category = "notes"
	PROJECTS Category = "proj"
)

type Post struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Category  Category  `json:"category"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Tags      []string  `json:"tags"`
}

func newPost(cat Category, title, content string, tags []string) *Post {
	now := time.Now()

	return &Post{
		Title:     title,
		Content:   content,
		Category:  cat,
		CreatedAt: now,
		UpdatedAt: now,
		Tags:      tags,
	}
}

// markdown -> html
func (p *Post) Html() string {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	np := parser.NewWithExtensions(extensions)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank | html.Safelink
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return string(markdown.ToHTML([]byte(p.Content), np, renderer))
}

// ============================================================================

func createPost(repo Repository, cat Category, title, content string, tags []string) (*Post, error) {
	md := newPost(cat, title, content, tags)

	if err := repo.SavePost(md); err != nil {
		return nil, err
	}

	return md, nil
}

func getPost(repo Repository, id int) (*Post, error) {
	return repo.GetPost(id)
}

func updatePost(repo Repository, id int, title, content string, tags []string) error {
	foundMd, err := repo.GetPost(id)
	if err != nil {
		return err
	}

	foundMd.Title = title
	foundMd.Content = content
	foundMd.UpdatedAt = time.Now()

	for _, t := range tags {
		if !slices.Contains(foundMd.Tags, t) {
			foundMd.Tags = append(foundMd.Tags, t)
		}
	}

	return repo.UpdatePost(foundMd)
}

func deletePost(repo Repository, id int) error {
	return repo.DeletePost(id)
}

func listAllPosts(repo Repository) ([]*Post, error) {
	return repo.ListPosts()
}
