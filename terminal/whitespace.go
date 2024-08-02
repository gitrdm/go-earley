package terminal

import (
	"unicode"

	"github.com/patrickhuber/go-earley/grammar"
)

type whitespace struct {
	grammar.SymbolImpl
}

// IsMatch implements grammar.Terminal.
func (*whitespace) IsMatch(ch rune) bool {
	return unicode.IsSpace(ch)
}

func NewWhitespace() grammar.Terminal {
	return &whitespace{}
}

func (*whitespace) String() string {
	return "\\s"
}
