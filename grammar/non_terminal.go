package grammar

type NonTerminal interface {
	Symbol
	Name() string
}

type nonTerminal struct {
	name string
}

func (nt *nonTerminal) Name() string {
	return nt.name
}

func (nt *nonTerminal) symbol() {}

func NewNonTerminal(name string) NonTerminal {
	return &nonTerminal{
		name: name,
	}
}
