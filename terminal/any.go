package terminal

import "github.com/patrickhuber/go-earley/grammar"

type Any interface {
}

type _any struct {
	grammar.SymbolImpl
}

// IsMatch implements grammar.Terminal.
func (*_any) IsMatch(ch rune) bool {
	return true
}

func NewAny() grammar.Terminal {
	return &_any{}
}
