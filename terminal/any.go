package terminal

import "github.com/patrickhuber/go-earley/grammar"

type Any struct {
	grammar.SymbolImpl
}

// IsMatch implements grammar.Terminal.
func (Any) IsMatch(ch rune) bool {
	return true
}

func NewAny() *Any {
	return &Any{}
}

func (Any) String() string {
	return ".*"
}
