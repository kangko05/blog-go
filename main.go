package main

import (
	"blog-go/internal/auth"
	"blog-go/internal/config"
	"blog-go/internal/post"
	"blog-go/internal/repo"
	"blog-go/internal/server"
)

func main() {
	cfg := config.Load()

	db, err := repo.ConnectDatabase(*cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	authRepo, err := repo.NewAuthRepository(db)
	if err != nil {
		panic(err)
	}

	postRepo, err := repo.NewPostRepository(db)
	if err != nil {
		panic(err)
	}

	authService, err := auth.NewService(cfg, authRepo)
	if err != nil {
		panic(err)
	}

	postService, err := post.NewService(postRepo)
	if err != nil {
		panic(err)
	}

	if err := server.SetupRouter(authService, postService).Run(":8000"); err != nil {
		panic(err)
	}
}
