package router

import (
	"blog-go/logger"

	"github.com/gin-gonic/gin"
)

var rlogger logger.Logger

func Init(l logger.Logger) {
	rlogger = l
}

func New() *gin.Engine {
	r := gin.Default()

	setupAuthRoutes(r)

	return r
}
