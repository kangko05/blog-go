package server

import (
	"blog-go/internal/auth"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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

type AccessLog struct {
	TimeStamp time.Time `json:"timestamp"`
	Ip        string    `json:"ip"`
	Method    string    `json:"method"`
	Path      string    `json:"path"`
	Status    int       `json:"status"`
	Duration  int64     `json:"duration"`
	UserAgent string    `json:"user_agent"`
	Referer   string    `json:"referer"`
}

func accessLogMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		now := time.Now()

		al := &AccessLog{
			TimeStamp: now,
			Ip:        ctx.ClientIP(),
			Method:    ctx.Request.Method,
			Path:      ctx.Request.URL.Path,
			UserAgent: ctx.Request.UserAgent(),
			Referer:   ctx.Request.Referer(),
		}

		ctx.Next()

		al.Status = ctx.Writer.Status()
		al.Duration = time.Since(now).Milliseconds()

		// write json file
		wb, err := json.Marshal(al)
		if err != nil {
			log.Println(err)
			return
		}

		file, err := os.OpenFile("logs/access.logs", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Println(err)
			return
		}
		defer file.Close()

		if _, err := file.Write(wb); err != nil {
			log.Println(err)
		}

		if _, err := file.WriteString("\n"); err != nil {
			log.Println(err)
		}
	}
}
