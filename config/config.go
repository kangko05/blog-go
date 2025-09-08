package config

import (
	"fmt"
	"os"
)

type Config struct {
	JwtSecret    string // for auth package
	DatabasePath string // for database package
}

func Load() (*Config, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("failed to read JWT_SECRET from env")
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		return nil, fmt.Errorf("failed to read DB_PATH from env")
	}

	if _, err := os.Stat(dbPath); err != nil {
		return nil, err
	}

	return &Config{
		JwtSecret:    jwtSecret,
		DatabasePath: dbPath,
	}, nil
}
