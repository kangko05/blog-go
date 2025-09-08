package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"blog-go/auth"
	"blog-go/config"
	"blog-go/database"
	"blog-go/logger"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	cfg := &config.Config{
		JwtSecret:    "test-secret",
		DatabasePath: ":memory:",
	}

	Init(logger.NewConsoleLogger())

	database.Init(cfg.DatabasePath)

	authRepo, _ := database.NewAuthRepository()
	auth.Init(authRepo)

	code := m.Run()

	database.Close()

	os.Exit(code)
}

func TestAuthAPI(t *testing.T) {

	router := New()

	t.Run("register user success", func(t *testing.T) {
		reqBody := RequestUser{
			UserName: "testuser",
			Password: "testpass123",
		}
		jsonData, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, "user registered", response.Message)
	})

	t.Run("register duplicate user", func(t *testing.T) {
		reqBody := RequestUser{
			UserName: "testuser",
			Password: "testpass123",
		}
		jsonData, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response Response
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.False(t, response.Success)
		assert.Contains(t, response.Error, "already exists")
	})

	t.Run("register invalid request", func(t *testing.T) {
		reqBody := RequestUser{
			UserName: "",
			Password: "testpass123",
		}
		jsonData, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response Response
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.False(t, response.Success)
		assert.NotEmpty(t, response.Error)
	})

	t.Run("login success", func(t *testing.T) {
		reqBody := RequestUser{
			UserName: "testuser",
			Password: "testpass123",
		}
		jsonData, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response Response
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response.Success)
		assert.Equal(t, "login success", response.Message)

		dataMap, ok := response.Data.(map[string]any)
		assert.True(t, ok)
		token, exists := dataMap["token"]
		assert.True(t, exists)
		assert.NotEmpty(t, token)
	})

	t.Run("login invalid credentials", func(t *testing.T) {
		reqBody := RequestUser{
			UserName: "testuser",
			Password: "wrongpass",
		}
		jsonData, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response Response
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.False(t, response.Success)
		assert.NotEmpty(t, response.Error)
	})

	t.Run("logout success", func(t *testing.T) {
		reqBody := RequestUser{
			UserName: "testuser",
			Password: "testpass123",
		}
		jsonData, _ := json.Marshal(reqBody)

		loginReq := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonData))
		loginReq.Header.Set("Content-Type", "application/json")

		loginW := httptest.NewRecorder()
		router.ServeHTTP(loginW, loginReq)

		var loginResponse Response
		json.Unmarshal(loginW.Body.Bytes(), &loginResponse)
		dataMap := loginResponse.Data.(map[string]any)
		token := dataMap["token"].(string)

		logoutReq := httptest.NewRequest("POST", "/api/auth/logout", nil)
		logoutReq.Header.Set("Authorization", "Bearer "+token)

		logoutW := httptest.NewRecorder()
		router.ServeHTTP(logoutW, logoutReq)

		assert.Equal(t, http.StatusOK, logoutW.Code)

		var logoutResponse Response
		json.Unmarshal(logoutW.Body.Bytes(), &logoutResponse)
		assert.True(t, logoutResponse.Success)
		assert.Equal(t, "logout success", logoutResponse.Message)
	})

	t.Run("logout missing auth header", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/auth/logout", nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response Response
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.False(t, response.Success)
		assert.Contains(t, response.Error, "authorization header required")
	})

	t.Run("logout invalid auth header", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/auth/logout", nil)
		req.Header.Set("Authorization", "InvalidToken")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response Response
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.False(t, response.Success)
		assert.Contains(t, response.Error, "failed to parse header")
	})
}
