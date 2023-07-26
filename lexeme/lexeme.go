package lexeme

import (
	"github.com/patrickhuber/go-earley/capture"
	"github.com/patrickhuber/go-earley/token"
)

// Lexeme is a mutable token
type Lexeme interface {
	token.Token
	Reset(offset int)
	Scan() bool
	Accepted() bool
}

type lexeme struct {
	capture   capture.Capture
	position  int
	tokenType string
}

// Capture implements Token.
func (t *lexeme) Capture() capture.Capture {
	return t.capture
}

// Position implements Token.
func (t *lexeme) Position() int {
	return t.position
}

// Type implements Token.
func (t *lexeme) Type() string {
	return t.tokenType
}
