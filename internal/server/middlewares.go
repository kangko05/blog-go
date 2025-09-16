package server

import (
	"blog-go/internal/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func authMiddleware(authService *auth.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.Status(http.StatusUnauthorized)
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if err := authService.VerifyToken(tokenString); err != nil {
			ctx.Status(http.StatusUnauthorized)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
