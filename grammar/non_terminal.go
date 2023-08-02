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

func (nt *nonTerminal) Equal(other Symbol) bool {
	otherNonTerminal, ok := other.(NonTerminal)
	if !ok {
		return false
	}
	return nt.name == otherNonTerminal.Name()
}
