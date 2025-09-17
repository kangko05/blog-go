package server

import (
	"blog-go/internal/auth"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
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

// access log =================================================================

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

// rate limiter================================================================

type RateLimiter struct {
	limiters map[string]*rate.Limiter
	lastUsed map[string]time.Time
	mu       sync.RWMutex
}

func (m *RateLimiter) cleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	go func() {
		for range ticker.C {
			m.mu.Lock()
			now := time.Now()
			for ip, lastTime := range m.lastUsed {
				if now.Sub(lastTime) > 30*time.Minute {
					delete(m.limiters, ip)
					delete(m.lastUsed, ip)
				}
			}
			m.mu.Unlock()
		}
	}()
}

func (m *RateLimiter) getLimiter(ip string) *rate.Limiter {
	m.mu.Lock()
	limiter, exists := m.limiters[ip]
	defer m.mu.Unlock()

	m.lastUsed[ip] = time.Now()

	if !exists {
		limiter = rate.NewLimiter(rate.Limit(10), 30)
		m.limiters[ip] = limiter
	}

	return limiter
}

func rateLimitMiddleware(limiter *RateLimiter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rlimiter := limiter.getLimiter(ctx.ClientIP())

		if !rlimiter.Allow() {
			ctx.Status(http.StatusTooManyRequests)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// security header ============================================================

func securityHeadersMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("X-Content-Type-Options", "nosniff")
		ctx.Header("X-Frame-Options", "DENY")
		ctx.Header("X-XSS-Protection", "1; mode=block")
		ctx.Header("Strict-Transport-Security", "max-age=31536000")
		ctx.Next()
	}
}
