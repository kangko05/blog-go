package post

import (
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
	Id        int
	Title     string
	Content   string
	Category  Category
	CreatedAt time.Time
	UpdatedAt time.Time
}

func newPost(cat Category, title, content string) *Post {
	now := time.Now()

	return &Post{
		Title:     title,
		Content:   content,
		Category:  cat,
		CreatedAt: now,
		UpdatedAt: now,
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

func createPost(repo Repository, cat Category, title, content string) (*Post, error) {
	md := newPost(cat, title, content)

	if err := repo.SavePost(md); err != nil {
		return nil, err
	}

	return md, nil
}

func getPost(repo Repository, id int) (*Post, error) {
	return repo.GetPost(id)
}

func updatePost(repo Repository, id int, title, content string) error {
	foundMd, err := repo.GetPost(id)
	if err != nil {
		return err
	}

	foundMd.Title = title
	foundMd.Content = content
	foundMd.UpdatedAt = time.Now()

	return repo.UpdatePost(foundMd)
}

func deletePost(repo Repository, id int) error {
	return repo.DeletePost(id)
}

func listAllPosts(repo Repository) ([]*Post, error) {
	return repo.ListPosts()
}
