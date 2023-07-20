package grammar

type Grammar interface {
	Start() NonTerminal
	Productions() []Production
}

type grammar struct {
	start       NonTerminal
	productions []Production
}

func (g *grammar) Start() NonTerminal {
	return g.start
}

func (g *grammar) Productions() []Production {
	return g.productions
}

func New(start NonTerminal, productions ...Production) Grammar {
	return &grammar{
		start:       start,
		productions: productions,
	}
}
