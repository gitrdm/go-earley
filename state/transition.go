package state

import "github.com/patrickhuber/go-earley/grammar"

const (
	TransitionType Type = 1
)

type Transition struct {
	Origin     int
	DottedRule *grammar.DottedRule
}

func (*Transition) Type() Type { return TransitionType }
