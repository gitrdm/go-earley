package grammar

type DottedRule interface {
	Production() Production
	Position() int
}

type dottedRule struct {
	production Production
	position   int
}

func (d *dottedRule) Production() Production {
	return d.production
}

func (d *dottedRule) Position() int {
	return d.position
}

func NewDottedRule(production Production, position int) DottedRule {
	return &dottedRule{
		production: production,
		position:   position,
	}
}

func PreDotSymbol(rule DottedRule) Symbol {
	if rule.Position() == 0 || len(rule.Production().RightHandSide()) == 0 {
		return nil
	}
	return rule.Production().RightHandSide()[rule.Position()]
}

func PostDotSymbol(rule DottedRule) Symbol {
	rhs := rule.Production().RightHandSide()
	if rule.Position() >= len(rhs) {
		return nil
	}
	return rhs[rule.Position()]
}
