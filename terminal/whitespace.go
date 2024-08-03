package terminal

import (
	"unicode"

	"github.com/patrickhuber/go-earley/grammar"
)

type Whitespace struct {
	grammar.SymbolImpl
}

// IsMatch implements grammar.Terminal.
func (*Whitespace) IsMatch(ch rune) bool {
	return unicode.IsSpace(ch)
}

func NewWhitespace() *Whitespace {
	return &Whitespace{}
}

func (*Whitespace) String() string {
	return "\\s"
}
