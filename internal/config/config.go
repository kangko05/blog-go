package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	JWTSECRET  = "JWTSECRET"
	DBPATH     = "DBPATH"
	SERVERPORT = "SERVERPORT"
)

type Config struct {
	ServerPort string
	JwtSecret  string
	DbPath     string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	jwtSecret := os.Getenv(JWTSECRET)
	if len(jwtSecret) == 0 {
		log.Println("[warn] using empty jwt secret")
	}

	dbpath := os.Getenv(DBPATH)
	if len(dbpath) == 0 {
		return nil, fmt.Errorf("empty dbpath provided")
	}

	srvPort := os.Getenv(SERVERPORT)
	if len(dbpath) == 0 {
		return nil, fmt.Errorf("empty server port provided")
	}

	return &Config{
		ServerPort: srvPort,
		JwtSecret:  jwtSecret,
		DbPath:     dbpath,
	}, nil
}
