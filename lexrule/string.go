package lexrule

import (
	"github.com/patrickhuber/go-earley/grammar"
)

const (
	StringType = "string"
)

type String struct {
	Value string
	grammar.SymbolImpl
}

// CanApply implements String.
func (s *String) CanApply(ch rune) bool {
	for _, r := range s.Value {
		return r == ch
	}
	// empty case
	return false
}

// ID implements String.
func (s *String) Type() string {
	return StringType
}

func NewString(str string) *String {
	return &String{
		Value: str,
	}
}

func (s *String) String() string {
	return s.Value
}
