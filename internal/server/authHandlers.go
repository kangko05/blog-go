package server

import (
	"blog-go/internal/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func handleRegister(authService *auth.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}

		if err := authService.Register(req.Username, req.Password); err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}

		ctx.Status(http.StatusOK)
	}
}

func handleLogin(authService *auth.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}

		tokenString, err := authService.Login(req.Username, req.Password)
		if err != nil {
			ctx.Status(http.StatusUnauthorized)
			return
		}

		ctx.String(http.StatusOK, tokenString)
	}
}

func handleLogout(authService *auth.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.Status(http.StatusBadRequest)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if err := authService.Logout(tokenString); err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}

		ctx.Status(http.StatusOK)
	}
}
