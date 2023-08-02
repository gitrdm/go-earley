package token

type Token interface {
	Capture() string
	Position() int
	Type() string
}

type token struct {
	capture   string
	position  int
	tokenType string
}

// Capture implements Token.
func (t *token) Capture() string {
	return t.capture
}

// Position implements Token.
func (t *token) Position() int {
	return t.position
}

// Type implements Token.
func (t *token) Type() string {
	return t.tokenType
}

func FromString(value string, position int, tokenType string) Token {
	return &token{
		capture:   value,
		position:  position,
		tokenType: tokenType,
	}
}
