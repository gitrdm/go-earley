package dfa

import "github.com/patrickhuber/go-earley/grammar"

const (
	DfaType string = "dfa"
)

type Dfa struct {
	State *State
	grammar.SymbolImpl
}

func (d *Dfa) CanApply(ch rune) bool {
	return d.State.IsMatch(ch)
}

func (d *Dfa) TokenType() string {
	// TODO: this should be the string representaiton of the dfa
	return DfaType
}

func (d *Dfa) Type() string {
	return DfaType
}

func New(s *State) Dfa {
	return Dfa{
		State: s,
	}
}

func (d *Dfa) String() string {
	return ""
}
