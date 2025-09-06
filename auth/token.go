package auth

type Token struct{}

func NewToken() *Token {
	return &Token{}
}
