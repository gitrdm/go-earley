package lexeme

import (
	"github.com/patrickhuber/go-earley/capture"
	"github.com/patrickhuber/go-earley/lexrule"
)

type String struct {
	lexeme
	str   string
	index int
}

// Accepted implements Lexeme.
func (s *String) Accepted() bool {
	return s.index == len(s.str)
}

// Position implements Lexeme.
func (s *String) Position() int {
	return s.index
}

// Reset implements Lexeme.
func (s *String) Reset(offset int) {
	s.index = -1
}

// Scan implements Lexeme.
func (*String) Scan() bool {
	panic("unimplemented")
}

// Type implements Lexeme.
func (*String) Type() string {
	panic("unimplemented")
}

func NewString(lexerRule lexrule.String, capture capture.Capture, offset int) *String {
	return &String{
		lexeme: lexeme{
			capture: capture,
		},
	}
}
