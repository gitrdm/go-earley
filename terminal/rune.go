package terminal

import (
	"github.com/patrickhuber/go-earley/grammar"
)

type character struct {
	ch rune
	grammar.SymbolImpl
}

// IsMatch implements grammar.Terminal.
func (r *character) IsMatch(ch rune) bool {
	return r.ch == ch
}

func NewRune(ch rune) grammar.Terminal {
	return &character{
		ch: ch,
	}
}

func (c *character) String() string {
	return string(c.ch)
}
