package cli

import (
	"blog-go/internal/config"
	"blog-go/internal/post"
	"blog-go/internal/repo"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	postService *post.Service
	cdb         *repo.Database
)

func init() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	db, err := repo.ConnectDatabase(cfg)
	if err != nil {
		fmt.Println(err)
		return
	}

	cdb = db

	postRepo, err := repo.NewPostRepository(db)
	if err != nil {
		fmt.Println(err)
		return
	}

	ps, err := post.NewService(postRepo)
	if err != nil {
		fmt.Println(err)
		return
	}

	postService = ps
}

func Init() error {

	return nil
}

var rootCmd = &cobra.Command{
	Use:   "blog",
	Short: "Personal blog post management CLI tool",
	Long: `A command line interface for managing personal blog posts.
Create, edit, and delete markdown posts that are stored in the database.`,
}

func Execute() {
	defer func() {
		if cdb != nil {
			cdb.Close()
		}
	}()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
