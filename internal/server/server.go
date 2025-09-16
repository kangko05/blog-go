package server

import (
	"blog-go/internal/auth"
	"blog-go/internal/post"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(authService *auth.Service, postService *post.Service) *gin.Engine {
	r := gin.Default()

	r.Use(accessLogMiddleware())

	r.GET("/checkhealth", handleCheckHealth())

	auth := r.Group("/auth")
	{
		auth.POST("/register", handleRegister(authService))
		auth.POST("/login", handleLogin(authService))
		auth.POST("/logout", handleLogout(authService))
	}

	post := r.Group("/posts")
	{
		post.GET("/", handleListAllPosts(postService))
		post.GET("/:id", handleGetPost(postService))

		post.POST("/", authMiddleware(authService), handleCreatePost(postService))
		post.PUT("/:id", authMiddleware(authService), handleUpdatePost(postService))
		post.DELETE("/:id", authMiddleware(authService), handleDeletePost(postService))
	}

	return r
}

func handleCheckHealth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	}
}
