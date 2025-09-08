package router

import (
	"blog-go/auth"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type RequestUser struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func setupAuthRoutes(r *gin.Engine) {
	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/register", handleRegisterUser)
		authGroup.POST("/login", handleLogin)
		authGroup.POST("/logout", handleLogout)
	}
}

func handleRegisterUser(ctx *gin.Context) {
	var reqUser RequestUser
	if err := ctx.ShouldBindJSON(&reqUser); err != nil {
		sendError(ctx, http.StatusBadRequest, err)
		return
	}

	if err := auth.RegisterUser(auth.NewUser(reqUser.UserName, reqUser.Password)); err != nil {
		sendError(ctx, http.StatusUnauthorized, err)
		return
	}

	sendSuccess(ctx, http.StatusCreated, "user registered", nil)
}

func handleLogin(ctx *gin.Context) {
	var reqUser RequestUser
	if err := ctx.ShouldBindJSON(&reqUser); err != nil {
		sendError(ctx, http.StatusBadRequest, err)
		return
	}

	token, err := auth.Login(auth.NewUser(reqUser.UserName, reqUser.Password))
	if err != nil {
		sendError(ctx, http.StatusBadRequest, err)
		return
	}

	sendSuccess(ctx, http.StatusOK, "login success", gin.H{"token": token.String})
}

func handleLogout(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		sendError(ctx, http.StatusUnauthorized, fmt.Errorf("authorization header required"))
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		sendError(ctx, http.StatusUnauthorized, fmt.Errorf("failed to parse header"))
		return
	}

	if err := auth.Logout(token); err != nil {
		sendError(ctx, http.StatusBadRequest, err)
		return
	}

	sendSuccess(ctx, http.StatusOK, "logout success", nil)
}
