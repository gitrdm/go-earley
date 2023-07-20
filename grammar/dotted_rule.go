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
