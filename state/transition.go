package state

import (
	"fmt"

	"github.com/patrickhuber/go-earley/grammar"
)

const (
	TransitionType Type = 1
)

type Transition struct {
	// Origin is the origin of the cached item, not the origin of the transition state
	Origin int
	// DottedRule is the dotted rule of the cached item, not the dotted rule of the transition state
	DottedRule *grammar.DottedRule
	// Symbol is the transition symbol
	Symbol grammar.Symbol
}

func (*Transition) Type() Type { return TransitionType }

func (t *Transition) String() string {
	return fmt.Sprintf("%s : %s, %d",
		t.Symbol, t.DottedRule, t.Origin)
}
