package router

import (
	"blog-go/logger"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const HandlerResultKey string = "HandlerResult"

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

		hr := &HandlerResult{
			RequestIP:     ctx.ClientIP(),
			RequestMethod: ctx.Request.Method,
			RequestURL:    ctx.Request.URL.RequestURI(),
		}

		ctx.Set(HandlerResultKey, hr)

		ctx.Next()

		duration := time.Since(start)

		hr.Duration = duration

		// page not found
		if hr.Status == 0 {
			hr.LogLevel = logger.LogError
			hr.Status = http.StatusNotFound
			hr.Message = "page not found"
		}

		rlogger.Write(hr.LogLevel, hr.Message, hr.MapResult())
	}
}

func New() *gin.Engine {
	r := gin.Default()

	r.Use(routerLogger())

	setupAuthRoutes(r)

	return r
}
