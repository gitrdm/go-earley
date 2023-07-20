package lexrule

import "github.com/patrickhuber/go-earley/grammar"

const (
	TerminalType = "terminal"
)

type Terminal interface {
	grammar.LexerRule
	Terminal() grammar.Terminal
}

type terminal struct {
	term grammar.Terminal
	grammar.SymbolImpl
}

// CanApply implements Terminal.
func (t *terminal) CanApply(ch rune) bool {
	return t.term.IsMatch(ch)
}

// ID implements Terminal.
func (t *terminal) Type() string {
	return TerminalType
}

// Terminal implements Terminal.
func (t *terminal) Terminal() grammar.Terminal {
	return t.term
}

func NewTerminal(t grammar.Terminal) Terminal {
	return &terminal{
		term: t,
	}
}
