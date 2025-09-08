package router

import "github.com/gin-gonic/gin"

func New() *gin.Engine {
	r := gin.Default()

	setupAuthRoutes(r)

	return r
}
