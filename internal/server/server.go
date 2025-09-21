package server

import (
	"blog-go/internal/auth"
	"blog-go/internal/post"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func SetupRouter(authService *auth.Service, postService *post.Service) *gin.Engine {
	r := gin.Default()

	config := cors.DefaultConfig()

	if gin.Mode() == gin.DebugMode {
		config.AllowOrigins = []string{
			"http://localhost:3000",
			"http://localhost:5173",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:5173",
		}
	} else {
		// config.AllowOrigins = []string{
		// 	"https://mydomain.com",
		// 	"https://www.mydomain.com",
		// }
	}

	config.AllowMethods = []string{"GET"}

	config.AllowHeaders = []string{
		"Origin", "Content-Type", "Accept", "Authorization",
		"X-Requested-With", "X-CSRF-Token",
	}

	config.ExposeHeaders = []string{
		"Content-Length", "Content-Type",
	}

	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	r.Use(cors.New(config))

	// middlewares
	limiter := &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		lastUsed: make(map[string]time.Time),
		mu:       sync.RWMutex{},
	}

	limiter.cleanup()

	r.Use(
		securityHeadersMiddleware(),
		accessLogMiddleware(),
		rateLimitMiddleware(limiter),
	)

	r.GET("/checkhealth", handleCheckHealth())
	r.GET("/notes", handleListNotes(postService))
	r.GET("/projects", handleListProjects(postService))
	r.GET("/posts/:id", handleGetPost(postService))

	// post := r.Group("/posts")
	// {
	// 	post.GET("", handleListAllPosts(postService))
	// 	post.GET("/:id", handleGetPost(postService))
	//
	// 	post.POST("/", authMiddleware(authService), handleCreatePost(postService))
	// 	post.PUT("/:id", authMiddleware(authService), handleUpdatePost(postService))
	// 	post.DELETE("/:id", authMiddleware(authService), handleDeletePost(postService))
	// }

	// auth := r.Group("/auth")
	// {
	// 	auth.POST("/register", handleRegister(authService))
	// 	auth.POST("/login", handleLogin(authService))
	// 	auth.POST("/logout", handleLogout(authService))
	// }

	return r
}

func handleCheckHealth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	}
}
