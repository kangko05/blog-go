package main

import (
	"blog-go/logger"
	"blog-go/router"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	router.Init(logger.NewConsoleLogger())

	if err := router.New().Run(":8000"); err != nil {
		panic(err)
	}
}
