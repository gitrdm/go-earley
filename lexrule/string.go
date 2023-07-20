package lexrule

import (
	"github.com/patrickhuber/go-earley/grammar"
)

const (
	StringType = "string"
)

type String interface {
	grammar.LexerRule
	Value() []rune
}

type stringLiteral struct {
	runes []rune
	grammar.SymbolImpl
}

// CanApply implements String.
func (s *stringLiteral) CanApply(ch rune) bool {
	for _, r := range s.runes {
		return r == ch
	}
	return false
}

// ID implements String.
func (s *stringLiteral) Type() string {
	return StringType
}

// Value implements String.
func (s *stringLiteral) Value() []rune {
	return s.runes
}

func NewString(str string) String {
	return &stringLiteral{
		runes: []rune(str),
	}
}
