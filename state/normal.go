package state

import (
	"fmt"
	"strconv"

	"github.com/patrickhuber/go-earley/forest"
	"github.com/patrickhuber/go-earley/grammar"
)

const (
	NormalType Type = 0
)

func NewNormal(rule *grammar.DottedRule, origin int) *Normal {
	return &Normal{
		DottedRule: rule,
		Origin:     origin,
	}
}

type Normal struct {
	Origin     int
	DottedRule *grammar.DottedRule
	Node       forest.Node
}

func (*Normal) Type() Type {
	return NormalType
}
func (n *Normal) String() string {
	return fmt.Sprintf("%s, %s",
		n.DottedRule, strconv.Itoa(n.Origin))
}
