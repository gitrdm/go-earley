package state

import (
	"fmt"
	"strconv"

	"github.com/patrickhuber/go-earley/grammar"
)

const (
	TransitionType Type = 1
)

type Transition struct {
	Origin     int
	DottedRule *grammar.DottedRule
	Symbol     grammar.Symbol
}

func (*Transition) Type() Type { return TransitionType }

func (t *Transition) String() string {
	return fmt.Sprintf("%s : %s, %s",
		t.Symbol, t.DottedRule, strconv.Itoa(t.Origin))
}
