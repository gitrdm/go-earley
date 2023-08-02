package grammar

type Grammar struct {
	Start       NonTerminal
	Productions []*Production
}

func New(start NonTerminal, productions ...*Production) *Grammar {
	return &Grammar{
		Start:       start,
		Productions: productions,
	}
}

func (g *Grammar) RulesFor(nt NonTerminal) []*Production {
	// TODO: optimize this if it becomes a memory hog
	var productions []*Production
	for p := range g.Productions {
		production := g.Productions[p]
		if production.LeftHandSide == nt {
			productions = append(productions, production)
		}
	}
	return productions
}

func (g *Grammar) StartProductions() []*Production {
	return g.RulesFor(g.Start)
}

func (g *Grammar) IsTransativeNullable(nt NonTerminal) bool {
	panic("not implemented")
}
