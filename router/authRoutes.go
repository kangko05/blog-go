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
		ctx.Error(fmt.Errorf("failed to bind request user to json: %v", err))
		sendError(ctx, http.StatusBadRequest, err)

		return
	}

	rlogger.Info("attempt to register", map[string]any{
		"username": reqUser.UserName,
	})

	if err := auth.RegisterUser(auth.NewUser(reqUser.UserName, reqUser.Password)); err != nil {
		ctx.Error(fmt.Errorf("failed to register user: %v", err))
		sendError(ctx, http.StatusUnauthorized, err)

		return
	}

	sendSuccess(ctx, http.StatusCreated, "user registered", nil)

	rlogger.Info("user registered", map[string]any{
		"username": reqUser.UserName,
		"ip":       ctx.ClientIP(),
	})
}

func handleLogin(ctx *gin.Context) {
	var reqUser RequestUser
	if err := ctx.ShouldBindJSON(&reqUser); err != nil {
		rlogger.Error("failed binding request user to json", map[string]any{
			"error": err.Error(),
			"ip":    ctx.ClientIP(),
		})

		sendError(ctx, http.StatusBadRequest, err)

		return
	}

	rlogger.Info("attempt to login", map[string]any{
		"username": reqUser.UserName,
	})

	token, err := auth.Login(auth.NewUser(reqUser.UserName, reqUser.Password))
	if err != nil {
		rlogger.Error("failed to login", map[string]any{
			"username": reqUser.UserName,
			"error":    err.Error(),
			"ip":       ctx.ClientIP(),
		})

		sendError(ctx, http.StatusBadRequest, err)
		return
	}

	sendSuccess(ctx, http.StatusOK, "login success", gin.H{"token": token.String})

	rlogger.Info("user logged in", map[string]any{
		"username": reqUser.UserName,
		"ip":       ctx.ClientIP(),
	})
}

func handleLogout(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")

	if authHeader == "" {
		errMsg := "authorization header required"

		rlogger.Error("failed to find auth header from request", map[string]any{
			"error": errMsg,
			"ip":    ctx.ClientIP(),
		})

		sendError(ctx, http.StatusUnauthorized, fmt.Errorf("%s", errMsg))

		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		errMsg := "failed to parse header"

		rlogger.Error("failed to parse auth header - no 'Bearer '", map[string]any{
			"error": errMsg,
			"ip":    ctx.ClientIP(),
		})

		sendError(ctx, http.StatusUnauthorized, fmt.Errorf("%s", errMsg))

		return
	}

	if err := auth.Logout(token); err != nil {
		rlogger.Error("failed to logout", map[string]any{
			"error": err.Error(),
			"ip":    ctx.ClientIP(),
		})

		sendError(ctx, http.StatusBadRequest, err)

		return
	}

	sendSuccess(ctx, http.StatusOK, "logout success", nil)

	rlogger.Info("user logged out", map[string]any{
		"ip": ctx.ClientIP(),
	})
}
