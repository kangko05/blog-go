package main

import (
	"blog-go/auth"
	"blog-go/database"
	"blog-go/logger"
	"blog-go/router"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	// init app
	if err := database.Init("dev.db"); err != nil {
		panic(err)
	}
	defer database.Close()

	ar, err := database.NewAuthRepository()
	if err != nil {
		panic(err)
	}

	auth.Init(ar)

	router.Init(logger.NewConsoleLogger())

	//

	if err := router.New().Run(":8000"); err != nil {
		panic(err)
	}
}
