package token

type Token interface {
	Position() int
	TokenType() string
}
