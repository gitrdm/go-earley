package lexeme

import (
	"unicode/utf8"

	"github.com/patrickhuber/go-earley/lexrule"
)

type String struct {
	lexeme
	index int
	rule  *lexrule.String
}

// Accepted implements Lexeme.
func (s *String) Accepted() bool {
	return s.index == len(s.rule.Value)
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
func (s *String) Scan(ch rune) bool {
	if s.index >= len(s.rule.Value) {
		return false
	}
	r, n := utf8.DecodeRuneInString(s.rule.Value[s.index:])
	if ch != r {
		return false
	}
	s.index += n
	return true
}

// Type implements Lexeme.
func (*String) Type() string {
	panic("unimplemented")
}

func NewString(lexerRule *lexrule.String, offset int) *String {
	return &String{
		rule:   lexerRule,
		lexeme: lexeme{},
	}
}
