package grammar

import (
	"strings"

	"github.com/patrickhuber/go-types"
	"github.com/patrickhuber/go-types/option"
)

type DottedRule struct {
	Production *Production
	Position   int
	preDot     Symbol
	postDot    Symbol
	str        string
}

func NewDottedRule(production *Production, position int) *DottedRule {
	return &DottedRule{
		Production: production,
		Position:   position,
	}
}

func (dr *DottedRule) Complete() bool {
	return dr.Position >= len(dr.Production.RightHandSide)
}

func (dr *DottedRule) PreDotSymbol() types.Option[Symbol] {
	if dr.preDot != nil {
		return option.Some(dr.preDot)
	}

	rhs := dr.Production.RightHandSide
	if dr.Position == 0 || len(rhs) == 0 {
		return option.None[Symbol]()
	}
	dr.preDot = rhs[dr.Position-1]
	return option.Some(dr.preDot)
}

func (dr *DottedRule) PostDotSymbol() types.Option[Symbol] {
	if dr.postDot != nil {
		return option.Some(dr.postDot)
	}
	rhs := dr.Production.RightHandSide
	if dr.Position >= len(rhs) {
		return option.None[Symbol]()
	}
	dr.postDot = rhs[dr.Position]
	return option.Some(dr.postDot)
}

func (dr *DottedRule) String() string {
	if len(dr.str) > 0 {
		return dr.str
	}
	sb := &strings.Builder{}
	sb.WriteString(dr.Production.LeftHandSide.Name())
	sb.WriteString(" ->")
	for i, s := range dr.Production.RightHandSide {
		if i == dr.Position {
			sb.WriteString("•")
		} else {
			sb.WriteString(" ")
		}
		sb.WriteString(s.String())
	}
	if len(dr.Production.RightHandSide) == dr.Position {
		sb.WriteString("•")
	}
	dr.str = sb.String()
	return dr.str
}
