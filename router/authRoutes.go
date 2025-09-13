package router

import (
	"blog-go/auth"
	"blog-go/logger"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type RequestUser struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// helper - get handler reulst from context

func getHandlerResult(ctx *gin.Context) *HandlerResult {
	v, exists := ctx.Get(HandlerResultKey)
	if !exists {
		v = &HandlerResult{
			LogLevel: logger.LogWarn,
			Message:  "failed to find handler result from context, setting default result",
		}

		ctx.Set(HandlerResultKey, v)
	}

	hr, ok := v.(*HandlerResult)
	if !ok {
		hr = &HandlerResult{
			LogLevel: logger.LogWarn,
			Message:  "wrong handler result type, setting default result",
		}

		ctx.Set(HandlerResultKey, hr)
	}

	return hr
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
	hr := getHandlerResult(ctx)

	var reqUser RequestUser
	if err := ctx.ShouldBindJSON(&reqUser); err != nil {

		hr.SetHttpStatus(
			logger.LogError,
			fmt.Sprintf("failed to bind request user to json: %v", err),
			http.StatusBadRequest,
		)

		sendError(ctx, hr.Status, err)

		return
	}

	if err := auth.RegisterUser(auth.NewUser(reqUser.UserName, reqUser.Password)); err != nil {
		hr.SetHttpStatus(
			logger.LogError,
			fmt.Sprintf("failed to register user: %v", err),
			http.StatusUnauthorized,
		)

		sendError(ctx, hr.Status, err)

		return
	}

	hr.SetHttpStatus(
		logger.LogInfo,
		fmt.Sprintf("user registered: %s", reqUser.UserName),
		http.StatusCreated,
	)

	sendSuccess(ctx, hr.Status, hr.Message, nil)
}

func handleLogin(ctx *gin.Context) {
	hr := getHandlerResult(ctx)

	var reqUser RequestUser
	if err := ctx.ShouldBindJSON(&reqUser); err != nil {
		hr.SetHttpStatus(
			logger.LogError,
			fmt.Sprintf("failed binding request user to json: %v", err),
			http.StatusBadRequest,
		)

		sendError(ctx, hr.Status, err)

		return
	}

	token, err := auth.Login(auth.NewUser(reqUser.UserName, reqUser.Password))
	if err != nil {

		hr.SetHttpStatus(
			logger.LogError,
			fmt.Sprintf("failed to login: %v", err),
			http.StatusBadRequest,
		)

		sendError(ctx, hr.Status, err)
		return
	}

	hr.SetHttpStatus(logger.LogInfo, "login success", http.StatusOK)

	sendSuccess(ctx, hr.Status, hr.Message, gin.H{"token": token.String})
}

func handleLogout(ctx *gin.Context) {
	hr := getHandlerResult(ctx)

	authHeader := ctx.GetHeader("Authorization")

	if authHeader == "" {
		errMsg := "authorization header required: "

		hr.SetHttpStatus(
			logger.LogError,
			errMsg+"failed to find auth header from request",
			http.StatusUnauthorized,
		)

		sendError(ctx, hr.Status, fmt.Errorf("%s", errMsg))

		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		errMsg := "failed to parse header: "

		hr.SetHttpStatus(
			logger.LogError,
			errMsg+"failed to parse auth header - no 'Bearer '",
			http.StatusUnauthorized,
		)

		sendError(ctx, hr.Status, fmt.Errorf("%s", errMsg))

		return
	}

	if err := auth.Logout(token); err != nil {
		hr.SetHttpStatus(
			logger.LogError,
			"failed to logout",
			http.StatusBadRequest,
		)

		sendError(ctx, hr.Status, err)

		return
	}

	hr.SetHttpStatus(
		logger.LogInfo,
		"logout success",
		http.StatusOK,
	)

	sendSuccess(ctx, hr.Status, hr.Message, nil)

}
