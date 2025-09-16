package auth

import (
	"blog-go/internal/config"
	"log"
)

type Service struct {
	repo      Repository
	jwtSecret string
}

func NewService(cfg *config.Config, repo Repository) (*Service, error) {
	jwtSecret := cfg.JwtSecret

	if repo == nil {
		log.Println("[warn] got nil repository, using in-memory repo")
		memRepo, err := connectSqlite()
		if err != nil {
			return nil, err
		}

		repo = memRepo
	}

	return &Service{
		repo:      repo,
		jwtSecret: jwtSecret,
	}, nil
}

func (as *Service) Register(username, password string) error {
	return register(as.repo, username, password)
}

// returns token string & error
func (as *Service) Login(username, password string) (string, error) {
	return login(as.repo, as.jwtSecret, username, password)
}

func (as *Service) Logout(tokenString string) error {
	return logout(as.repo, tokenString)
}

func (as *Service) VerifyToken(tokenString string) error {
	return verifyToken(as.repo, as.jwtSecret, tokenString)
}
