package lexeme

// Lexeme is a mutable token
type Lexeme interface {
	Scan(ch rune) bool
	Accepted() bool
}

type lexeme struct {
	span      *Span
	tokenType string
}

// Position implements Token.
func (t *lexeme) Position() int {
	return t.span.Offset
}

// Type implements Token.
func (t *lexeme) Type() string {
	return t.tokenType
}
