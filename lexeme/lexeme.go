package lexeme

import (
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
	capture   string
	span      *Span
	tokenType string
}

// Capture implements Token.
func (t *lexeme) Capture() string {
	return t.capture[t.span.Offset:t.span.Length]
}

// Position implements Token.
func (t *lexeme) Position() int {
	return t.span.Offset
}

// Type implements Token.
func (t *lexeme) Type() string {
	return t.tokenType
}
