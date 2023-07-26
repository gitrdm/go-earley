package grammar

type Grammar struct {
	Start       []NonTerminal
	Productions []Production
}

func New(start []NonTerminal, productions []Production) *Grammar {
	return &Grammar{
		Start:       start,
		Productions: productions,
	}
}
