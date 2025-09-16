package post

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

<<<<<<< HEAD
func TestMain(m *testing.M) {
	if err := Init(nil); err != nil {
=======
var postService *Service

func TestMain(m *testing.M) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "PostRepository", nil)

	var err error
	postService, err = NewService(ctx)
	if err != nil {
>>>>>>> 6bd65ff (router done for now)
		panic(err)
	}

	exitCode := m.Run()
	repo.Close()
	os.Exit(exitCode)
}

func TestCreatePost(t *testing.T) {
	t.Run("create post success", func(t *testing.T) {
		title, content := "Test Title", "# Test Content\n\nThis is a test post."

		post, err := CreatePost(title, content)

		assert.Nil(t, err)
		assert.Equal(t, title, post.Title)
		assert.Equal(t, content, post.Content)
		assert.False(t, post.CreatedAt.IsZero())
		assert.False(t, post.UpdatedAt.IsZero())
	})
}

func TestGetPost(t *testing.T) {
	title, content := "Get Test", "Content for get test"
	createdPost, err := CreatePost(title, content)
	assert.Nil(t, err)

	t.Run("get post success", func(t *testing.T) {
		post, err := GetPost(createdPost.Id)
		assert.Nil(t, err)
		assert.Equal(t, title, post.Title)
		assert.Equal(t, content, post.Content)
		assert.Equal(t, createdPost.Id, post.Id)
	})

	t.Run("get non-existent post", func(t *testing.T) {
		_, err := GetPost(99999)
		assert.NotNil(t, err)
	})
}

// func TestUpdatePost(t *testing.T) {
// 	originalTitle, originalContent := "Original Title", "Original content"
// 	createdPost, _ := CreatePost(originalTitle, originalContent)
//
// 	t.Run("update post success", func(t *testing.T) {
// 		newTitle, newContent := "Updated Title", "Updated content"
//
// 		err := UpdatePost(createdPost.Id, newTitle, newContent)
// 		assert.Nil(t, err)
//
// 		updatedPost, err := GetPost(createdPost.Id)
// 		assert.Nil(t, err)
// 		assert.Equal(t, newTitle, updatedPost.Title)
// 		assert.Equal(t, newContent, updatedPost.Content)
// 		assert.Equal(t, createdPost.CreatedAt, updatedPost.CreatedAt)
// 		assert.True(t, updatedPost.UpdatedAt.After(createdPost.UpdatedAt))
// 	})
//
// 	t.Run("update non-existent post", func(t *testing.T) {
// 		err := UpdatePost(99999, "title", "content")
// 		assert.NotNil(t, err)
// 	})
// }

func TestDeletePost(t *testing.T) {
	title, content := "Delete Test", "Content to be deleted"
	createdPost, _ := CreatePost(title, content)

	t.Run("delete post success", func(t *testing.T) {
		err := DeletePost(createdPost.Id)
		assert.Nil(t, err)

		_, err = GetPost(createdPost.Id)
		assert.NotNil(t, err)
	})

	t.Run("delete non-existent post", func(t *testing.T) {
		err := DeletePost(99999)
		assert.Nil(t, err)
	})
}

func TestListAllPosts(t *testing.T) {
	posts := []struct{ title, content string }{
		{"Post 1", "Content 1"},
		{"Post 2", "Content 2"},
		{"Post 3", "Content 3"},
	}

	var createdIds []int
	for _, p := range posts {
		created, _ := CreatePost(p.title, p.content)
		createdIds = append(createdIds, created.Id)
	}

	t.Run("list all posts", func(t *testing.T) {
		allPosts, err := ListAllPosts()
		assert.Nil(t, err)
		assert.GreaterOrEqual(t, len(allPosts), 3)

		found := 0
		for _, post := range allPosts {
			for _, id := range createdIds {
				if post.Id == id {
					found++
				}
			}
		}
		assert.Equal(t, 3, found)
	})
}

func TestMarkdownToHTML(t *testing.T) {
	t.Run("markdown conversion", func(t *testing.T) {
		title := "Markdown Test"
		content := "# Header\n\n**Bold text** and *italic text*\n\n- List item 1\n- List item 2"

		post, err := CreatePost(title, content)
		assert.Nil(t, err)

		html := post.Html()
		assert.Contains(t, html, "<h1")
		assert.Contains(t, html, "<strong>")
		assert.Contains(t, html, "<em>")
		assert.Contains(t, html, "<ul>")
		assert.Contains(t, html, "<li>")
	})
}
