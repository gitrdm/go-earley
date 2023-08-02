package grammar

type Production struct {
	LeftHandSide  NonTerminal
	RightHandSide []Symbol
}

func NewProduction(lhs NonTerminal, rhs ...Symbol) *Production {
	return &Production{
		LeftHandSide:  lhs,
		RightHandSide: rhs,
	}
}