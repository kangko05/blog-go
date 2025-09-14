package router

import (
	"blog-go/logger"

	"github.com/gin-gonic/gin"
)

const HandlerResultKey string = "HandlerResult"

var rlogger logger.Logger

func Init(l logger.Logger) {
	rlogger = l
}

func New() *gin.Engine {
	r := gin.Default()

	// r.Use(routerLogger())

	setupAuthRoutes(r)

	return r
}
