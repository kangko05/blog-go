package router

import (
	"blog-go/logger"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

var rlogger logger.Logger

func Init(l logger.Logger) {
	rlogger = l
}

func routerLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		log := fmt.Sprintf("request received %s %s", ctx.Request.Method, ctx.Request.URL.RequestURI())
		rlogger.Info(log, map[string]any{
			"requestIP": ctx.ClientIP(),
		})

		ctx.Next()

		duration := time.Since(start)

		result := map[string]any{
			"requestIP": ctx.ClientIP(),
			"status":    ctx.Writer.Status(),
			"duration":  duration,
		}

		if len(ctx.Errors) > 0 {
			result["error"] = ctx.Errors.Last().Error()
			rlogger.Error("request failed", result)
		} else {
			rlogger.Info("request completed", result)
		}
	}
}

func New() *gin.Engine {
	r := gin.Default()

	r.Use(routerLogger())

	setupAuthRoutes(r)

	return r
}
