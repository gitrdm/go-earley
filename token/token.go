package token

import "github.com/patrickhuber/go-earley/capture"

type Token interface {
	Capture() capture.Capture
	Position() int
	Type() string
}

type token struct {
	capture   capture.Capture
	position  int
	tokenType string
}

// Capture implements Token.
func (t *token) Capture() capture.Capture {
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
		capture:   capture.FromString(value),
		position:  position,
		tokenType: tokenType,
	}
}

func FromCapture(capture capture.Capture, position int, tokenType string) Token {
	return &token{
		capture:   capture,
		position:  position,
		tokenType: tokenType,
	}
}
