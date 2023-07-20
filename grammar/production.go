package grammar

type Production interface {
	LeftHandSide() NonTerminal
	RightHandSide() []Symbol
}

type production struct {
	leftHandSide  NonTerminal
	rightHandSide []Symbol
}

func (p *production) LeftHandSide() NonTerminal {
	return p.leftHandSide
}

func (p *production) RightHandSide() []Symbol {
	return p.rightHandSide
}

func NewProduction(lhs NonTerminal, rhs ...Symbol) Production {
	return &production{
		leftHandSide:  lhs,
		rightHandSide: rhs,
	}
}
