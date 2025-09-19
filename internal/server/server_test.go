package server

// func setupTestServer(t *testing.T) (*gin.Engine, func()) {
// 	gin.SetMode(gin.TestMode)
//
// 	cfg := &config.Config{
// 		JwtSecret: "test-secret",
// 		DbPath:    ":memory:",
// 	}
//
// 	db, err := repo.ConnectDatabase(cfg)
// 	assert.Nil(t, err)
//
// 	authRepo, err := repo.NewAuthRepository(db)
// 	assert.Nil(t, err)
//
// 	postRepo, err := repo.NewPostRepository(db)
// 	assert.Nil(t, err)
//
// 	authService, err := auth.NewService(cfg, authRepo)
// 	assert.Nil(t, err)
//
// 	postService, err := post.NewService(postRepo)
// 	assert.Nil(t, err)
//
// 	router := SetupRouter(authService, postService)
//
// 	cleanup := func() {
// 		db.Close()
// 	}
//
// 	return router, cleanup
// }

// func TestAuthFlow(t *testing.T) {
// 	router, cleanup := setupTestServer(t)
// 	defer cleanup()
//
// 	t.Run("register and login", func(t *testing.T) {
// 		regReq := map[string]string{
// 			"username": "testuser",
// 			"password": "testpass",
// 		}
// 		reqBody, _ := json.Marshal(regReq)
//
// 		req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqBody))
// 		req.Header.Set("Content-Type", "application/json")
// 		w := httptest.NewRecorder()
//
// 		router.ServeHTTP(w, req)
// 		assert.Equal(t, http.StatusOK, w.Code)
//
// 		req = httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))
// 		req.Header.Set("Content-Type", "application/json")
// 		w = httptest.NewRecorder()
//
// 		router.ServeHTTP(w, req)
// 		assert.Equal(t, http.StatusOK, w.Code)
// 		assert.NotEmpty(t, w.Body.String())
// 	})
// }
//
// func TestPostFlow(t *testing.T) {
// 	router, cleanup := setupTestServer(t)
// 	defer cleanup()
//
// 	regReq := map[string]string{
// 		"username": "blogger",
// 		"password": "blogpass",
// 	}
// 	reqBody, _ := json.Marshal(regReq)
//
// 	req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqBody))
// 	req.Header.Set("Content-Type", "application/json")
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)
//
// 	req = httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))
// 	req.Header.Set("Content-Type", "application/json")
// 	w = httptest.NewRecorder()
// 	router.ServeHTTP(w, req)
// 	token := w.Body.String()
//
// 	t.Run("create post without auth should fail", func(t *testing.T) {
// 		postReq := map[string]string{
// 			"title":   "Test Post",
// 			"content": "Test Content",
// 		}
// 		reqBody, _ := json.Marshal(postReq)
//
// 		req := httptest.NewRequest("POST", "/posts/", bytes.NewBuffer(reqBody))
// 		req.Header.Set("Content-Type", "application/json")
// 		w := httptest.NewRecorder()
//
// 		router.ServeHTTP(w, req)
// 		assert.Equal(t, http.StatusUnauthorized, w.Code)
// 	})
//
// 	t.Run("create post with auth should succeed", func(t *testing.T) {
// 		postReq := map[string]string{
// 			"title":   "My Blog Post",
// 			"content": "# Hello\n\nThis is my first post!",
// 		}
// 		reqBody, _ := json.Marshal(postReq)
//
// 		req := httptest.NewRequest("POST", "/posts/", bytes.NewBuffer(reqBody))
// 		req.Header.Set("Content-Type", "application/json")
// 		req.Header.Set("Authorization", "Bearer "+token)
// 		w := httptest.NewRecorder()
//
// 		router.ServeHTTP(w, req)
// 		assert.Equal(t, http.StatusOK, w.Code)
//
// 		var createdPost post.Post
// 		err := json.Unmarshal(w.Body.Bytes(), &createdPost)
// 		assert.Nil(t, err)
// 		assert.Equal(t, "My Blog Post", createdPost.Title)
// 		assert.Greater(t, createdPost.Id, 0)
// 	})
//
// 	t.Run("get posts should work without auth", func(t *testing.T) {
// 		req := httptest.NewRequest("GET", "/posts/", nil)
// 		w := httptest.NewRecorder()
//
// 		router.ServeHTTP(w, req)
// 		assert.Equal(t, http.StatusOK, w.Code)
//
// 		var posts []*post.Post
// 		err := json.Unmarshal(w.Body.Bytes(), &posts)
// 		assert.Nil(t, err)
// 		assert.GreaterOrEqual(t, len(posts), 1)
// 	})
// }
//
// func TestFullFlow(t *testing.T) {
// 	router, cleanup := setupTestServer(t)
// 	defer cleanup()
//
// 	regReq := map[string]string{
// 		"username": "fulltest",
// 		"password": "fullpass",
// 	}
// 	reqBody, _ := json.Marshal(regReq)
//
// 	req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqBody))
// 	req.Header.Set("Content-Type", "application/json")
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)
// 	assert.Equal(t, http.StatusOK, w.Code)
//
// 	req = httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))
// 	req.Header.Set("Content-Type", "application/json")
// 	w = httptest.NewRecorder()
// 	router.ServeHTTP(w, req)
// 	assert.Equal(t, http.StatusOK, w.Code)
// 	token := w.Body.String()
//
// 	postReq := map[string]string{
// 		"title":   "Integration Test",
// 		"content": "Full flow test content",
// 	}
// 	reqBody, _ = json.Marshal(postReq)
//
// 	req = httptest.NewRequest("POST", "/posts/", bytes.NewBuffer(reqBody))
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Authorization", "Bearer "+token)
// 	w = httptest.NewRecorder()
// 	router.ServeHTTP(w, req)
// 	assert.Equal(t, http.StatusOK, w.Code)
//
// 	updateReq := map[string]string{
// 		"title":   "Updated Title",
// 		"content": "Updated content",
// 	}
// 	reqBody, _ = json.Marshal(updateReq)
//
// 	req = httptest.NewRequest("PUT", "/posts/1", bytes.NewBuffer(reqBody))
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Authorization", "Bearer "+token)
// 	w = httptest.NewRecorder()
// 	router.ServeHTTP(w, req)
// 	assert.Equal(t, http.StatusOK, w.Code)
//
// 	req = httptest.NewRequest("POST", "/auth/logout", nil)
// 	req.Header.Set("Authorization", "Bearer "+token)
// 	w = httptest.NewRecorder()
// 	router.ServeHTTP(w, req)
// 	assert.Equal(t, http.StatusOK, w.Code)
//
// 	req = httptest.NewRequest("DELETE", "/posts/1", nil)
// 	req.Header.Set("Authorization", "Bearer "+token)
// 	w = httptest.NewRecorder()
// 	router.ServeHTTP(w, req)
// 	assert.Equal(t, http.StatusUnauthorized, w.Code)
// }
