package router

import "github.com/gin-gonic/gin"

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func sendSuccess(ctx *gin.Context, code int, msg string, data any) {
	ctx.JSON(code, &Response{
		Success: true,
		Message: msg,
		Data:    data,
	})
}

func sendError(ctx *gin.Context, code int, err error) {
	ctx.JSON(code, &Response{
		Success: false,
		Error:   err.Error(),
	})
}
