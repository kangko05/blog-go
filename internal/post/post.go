package post

import (
	"log"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type Post struct {
	Id        int
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func newPost(title, content string) *Post {
	now := time.Now()

	return &Post{
		Title:     title,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// markdown -> html
func (p *Post) Html() string {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	np := parser.NewWithExtensions(extensions)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return string(markdown.ToHTML([]byte(p.Content), np, renderer))
}

// ============================================================================

var repo Repository

func Init(r Repository) error {
	repo = r

	if r == nil {
		log.Println("[warn] got nil repository, using in-memory repo")
		memRepo, err := connectSqlite()
		if err != nil {
			return err
		}

		repo = memRepo
	}

	return nil
}

// ============================================================================

func CreatePost(title, content string) (*Post, error) {
	md := newPost(title, content)

	if err := repo.SavePost(md); err != nil {
		return nil, err
	}

	return md, nil
}

func GetPost(id int) (*Post, error) {
	return repo.GetPost(id)
}

func UpdatePost(id int, title, content string) error {
	foundMd, err := repo.GetPost(id)
	if err != nil {
		return err
	}

	foundMd.Title = title
	foundMd.Content = content
	foundMd.UpdatedAt = time.Now()

	return repo.UpdatePost(foundMd)
}

func DeletePost(id int) error {
	return repo.DeletePost(id)
}

func ListAllPosts() ([]*Post, error) {
	return repo.ListPosts()
}
