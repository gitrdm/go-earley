package grammar

type DottedRule struct {
	Production *Production
	Position   int
	preDot     Symbol
	postDot    Symbol
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

func (dr *DottedRule) PreDotSymbol() Symbol {
	if dr.preDot != nil {
		return dr.preDot
	}

	if dr.Position == 0 || len(dr.Production.RightHandSide) == 0 {
		return nil
	}
	dr.preDot = dr.Production.RightHandSide[dr.Position]
	return dr.preDot
}

func (dr *DottedRule) PostDotSymbol() Symbol {
	if dr.postDot != nil {
		return dr.postDot
	}
	rhs := dr.Production.RightHandSide
	if dr.Position >= len(rhs) {
		return nil
	}
	dr.postDot = rhs[dr.Position]
	return dr.postDot
}
