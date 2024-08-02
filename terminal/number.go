package terminal

import (
	"unicode"

	"github.com/patrickhuber/go-earley/grammar"
)

type number struct {
	grammar.SymbolImpl
}

// IsMatch implements grammar.Terminal.
func (*number) IsMatch(ch rune) bool {
	return unicode.IsNumber(ch)
}

func NewNumber() grammar.Terminal {
	return &number{}
}

func (*number) String() string {
	return "\\d"
}
