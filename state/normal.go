package state

import (
	"go/ast"

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
	DottedRule *grammar.DottedRule
	Origin     int
	Node       ast.Node
}

func (*Normal) Type() Type {
	return NormalType
}
