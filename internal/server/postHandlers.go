package server

import (
	"blog-go/internal/post"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func handleListAllPosts(postService *post.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		posts, err := postService.ListAllPosts()
		if err != nil {
			ctx.Status(http.StatusInternalServerError)
			return
		}

		ctx.JSON(http.StatusOK, posts)
	}
}

func handleGetPost(postService *post.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idstr := ctx.Param("id")

		id, err := strconv.Atoi(idstr)
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}

		post, err := postService.GetPost(id)
		if err != nil {
			ctx.Status(http.StatusInternalServerError)
			return
		}

		ctx.JSON(http.StatusOK, post)
	}
}

func handleCreatePost(postService *post.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req struct {
			Title    string        `json:"title"`
			Category post.Category `json:"category"`
			Content  string        `json:"content"`
		}

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}

		post, err := postService.CreatePost(req.Category, req.Title, req.Content)
		if err != nil {
			ctx.Status(http.StatusInternalServerError)
			return
		}

		ctx.JSON(http.StatusOK, post)
	}
}

func handleUpdatePost(postService *post.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idstr := ctx.Param("id")

		id, err := strconv.Atoi(idstr)
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}

		var req struct {
			Title   string `json:"title"`
			Content string `json:"content"`
		}

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}

		if err := postService.UpdatePost(id, req.Title, req.Content); err != nil {
			ctx.Status(http.StatusInternalServerError)
			return
		}

		ctx.Status(http.StatusOK)
	}
}

func handleDeletePost(postService *post.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idstr := ctx.Param("id")

		id, err := strconv.Atoi(idstr)
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}

		if err := postService.DeletePost(id); err != nil {
			ctx.Status(http.StatusInternalServerError)
			return
		}

		ctx.Status(http.StatusOK)
	}
}
