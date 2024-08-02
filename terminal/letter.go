package terminal

import (
	"unicode"

	"github.com/patrickhuber/go-earley/grammar"
)

type letter struct {
	grammar.SymbolImpl
}

// IsMatch implements grammar.Terminal.
func (*letter) IsMatch(ch rune) bool {
	return unicode.IsLetter(ch)
}

func NewLetter() grammar.Terminal {
	return &letter{}
}

func (*letter) String() string {
	return "\\w"
}
