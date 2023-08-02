package lexeme

import (
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
func (s *String) Scan() bool {
	if s.index >= len(s.capture) {
		return false
	}
	if s.index >= len(s.rule.Value) {
		return false
	}
	if s.capture[s.index] != s.rule.Value[s.index] {
		return false
	}
	s.index++
	return true
}

// Type implements Lexeme.
func (*String) Type() string {
	panic("unimplemented")
}

func NewString(lexerRule *lexrule.String, str string, offset int) *String {
	return &String{
		rule: lexerRule,
		lexeme: lexeme{
			capture: str,
		},
	}
}
