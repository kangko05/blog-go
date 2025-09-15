package session

import (
	"blog-go/internal/user"
)

type AppContext interface {
	GetUserRepo() user.Repository
	GetSessionRepo() Repository
}

type Repository interface{}

type Session struct {
	userRepo    user.Repository
	sessionRepo Repository
}

func New(ctx AppContext) *Session {
	return &Session{
		userRepo:    ctx.GetUserRepo(),
		sessionRepo: ctx.GetSessionRepo(),
	}
}

func (s *Session) Login(username, password string) error {
	return user.Validate(s.userRepo, username, password)
}

