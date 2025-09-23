package post

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestService(t *testing.T) (*Service, func()) {
	service, err := NewService(nil)
	assert.NoError(t, err)

	cleanup := func() {}

	return service, cleanup
}

func TestCreatePost(t *testing.T) {
	service, cleanup := setupTestService(t)
	defer cleanup()

	t.Run("create post success", func(t *testing.T) {
		title, content := "Test Title", "# Test Content\n\nThis is a test post."
		tags := []string{"this", "is", "test"}

		post, err := service.CreatePost(NOTES, title, content, tags)

		assert.Nil(t, err)
		assert.Equal(t, title, post.Title)
		assert.Equal(t, content, post.Content)
		assert.False(t, post.CreatedAt.IsZero())
		assert.False(t, post.UpdatedAt.IsZero())
	})
}

func TestGetPost(t *testing.T) {
	service, cleanup := setupTestService(t)
	tags := []string{"this", "is", "test"}
	defer cleanup()

	title, content := "Get Test", "Content for get test"
	createdPost, err := service.CreatePost(NOTES, title, content, tags)
	assert.Nil(t, err)

	t.Run("get post success", func(t *testing.T) {
		post, err := service.GetPost(createdPost.Id)
		assert.Nil(t, err)
		assert.Equal(t, title, post.Title)
		assert.Equal(t, content, post.Content)
		assert.Equal(t, createdPost.Id, post.Id)
	})

	t.Run("get non-existent post", func(t *testing.T) {
		_, err := service.GetPost(99999)
		assert.NotNil(t, err)
	})
}

func TestUpdatePost(t *testing.T) {
	service, cleanup := setupTestService(t)
	tags := []string{"this", "is", "test"}
	defer cleanup()

	originalTitle, originalContent := "Original Title", "Original content"
	createdPost, _ := service.CreatePost(NOTES, originalTitle, originalContent, tags)

	t.Run("update post success", func(t *testing.T) {
		newTitle, newContent := "Updated Title", "Updated content"
		tags := []string{"this", "is", "test"}

		err := service.UpdatePost(createdPost.Id, newTitle, newContent, tags)
		assert.Nil(t, err)

		updatedPost, err := service.GetPost(createdPost.Id)
		assert.Nil(t, err)
		assert.Equal(t, newTitle, updatedPost.Title)
		assert.Equal(t, newContent, updatedPost.Content)
		assert.True(t, createdPost.CreatedAt.Equal(updatedPost.CreatedAt))
		assert.True(t, updatedPost.UpdatedAt.After(createdPost.UpdatedAt))
	})

	t.Run("update non-existent post", func(t *testing.T) {
		err := service.UpdatePost(99999, "title", "content", tags)
		assert.NotNil(t, err)
	})
}

func TestDeletePost(t *testing.T) {
	service, cleanup := setupTestService(t)
	tags := []string{"this", "is", "test"}
	defer cleanup()

	title, content := "Delete Test", "Content to be deleted"
	createdPost, _ := service.CreatePost(NOTES, title, content, tags)

	t.Run("delete post success", func(t *testing.T) {
		err := service.DeletePost(createdPost.Id)
		assert.Nil(t, err)

		_, err = service.GetPost(createdPost.Id)
		assert.NotNil(t, err)
	})

	t.Run("delete non-existent post", func(t *testing.T) {
		err := service.DeletePost(99999)
		assert.Nil(t, err)
	})
}

func TestListAllPosts(t *testing.T) {
	service, cleanup := setupTestService(t)
	defer cleanup()

	posts := []struct{ title, content string }{
		{"Post 1", "Content 1"},
		{"Post 2", "Content 2"},
		{"Post 3", "Content 3"},
	}
	tags := []string{"this", "is", "test"}

	var createdIds []int
	for _, p := range posts {
		created, _ := service.CreatePost(NOTES, p.title, p.content, tags)
		createdIds = append(createdIds, created.Id)
	}

	t.Run("list all posts", func(t *testing.T) {
		allPosts, err := service.ListAllPosts()
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
	service, cleanup := setupTestService(t)
	tags := []string{"this", "is", "test"}
	defer cleanup()

	t.Run("markdown conversion", func(t *testing.T) {
		title := "Markdown Test"
		content := "# Header\n\n**Bold text** and *italic text*\n\n- List item 1\n- List item 2"

		post, err := service.CreatePost(NOTES, title, content, tags)
		assert.Nil(t, err)

		html := post.Html()
		assert.Contains(t, html, "<h1")
		assert.Contains(t, html, "<strong>")
		assert.Contains(t, html, "<em>")
		assert.Contains(t, html, "<ul>")
		assert.Contains(t, html, "<li>")
	})
}
