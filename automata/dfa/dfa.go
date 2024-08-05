package dfa

import "github.com/patrickhuber/go-earley/grammar"

const (
	LexerRuleType = "dfa"
)

func NewDfa(start *State, tokenType string) *Dfa {
	return &Dfa{
		Start:     start,
		tokenType: tokenType,
	}
}

type Dfa struct {
	grammar.SymbolImpl
	Start     *State
	tokenType string
}

func (d *Dfa) CanApply(ch rune) bool {
	// check for matches of the first state transitions
	for _, trans := range d.Start.Transitions {
		if trans.Terminal.IsMatch(ch) {
			return true
		}
	}
	return false
}

func (d Dfa) TokenType() string {
	return d.tokenType
}

func (Dfa) LexerRuleType() string {
	return LexerRuleType
}

func (d Dfa) String() string {
	return d.tokenType
}
