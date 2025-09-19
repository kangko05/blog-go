package main

import (
	"blog-go/internal/auth"
	"blog-go/internal/config"
	"blog-go/internal/post"
	"blog-go/internal/repo"
	"blog-go/internal/server"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, db, router, err := prepare()
	if err != nil {
		panic(err)
	}

	srv := &http.Server{
		Addr:    cfg.ServerPort,
		Handler: router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cleanup(ctx, srv, db)
}

func cleanup(ctx context.Context, srv *http.Server, db *repo.Database) {
	log.Println("shutting down...")

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("srv shutdown: %v\n", err)
	}

	if err := db.Close(); err != nil {
		log.Printf("db close: %v\n", err)
	}

	log.Println("gracefully shutdown the app")
}

func prepare() (*config.Config, *repo.Database, *gin.Engine, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, nil, nil, err
	}

	db, err := repo.ConnectDatabase(cfg)
	if err != nil {
		return nil, nil, nil, err
	}

	authRepo, err := repo.NewAuthRepository(db)
	if err != nil {
		return nil, nil, nil, err
	}

	postRepo, err := repo.NewPostRepository(db)
	if err != nil {
		return nil, nil, nil, err
	}

	authService, err := auth.NewService(cfg, authRepo)
	if err != nil {
		return nil, nil, nil, err
	}

	postService, err := post.NewService(postRepo)
	if err != nil {
		return nil, nil, nil, err
	}

	router := server.SetupRouter(authService, postService)

	return cfg, db, router, nil
}
