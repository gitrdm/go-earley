package terminal

import (
	"fmt"

	"github.com/patrickhuber/go-earley/grammar"
)

type Negate struct {
	grammar.SymbolImpl
	terminal grammar.Terminal
}

func NewNegate(terminal grammar.Terminal) *Negate {
	return &Negate{
		terminal: terminal,
	}
}

func (n *Negate) IsMatch(ch rune) bool {
	return !n.terminal.IsMatch(ch)
}

func (n *Negate) String() string {
	return fmt.Sprintf("[^%v]", n.terminal)
}
