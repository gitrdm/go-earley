package dfa

type Lexeme struct {
	d       *Dfa
	current *State
}

func NewLexeme(d *Dfa) *Lexeme {
	return &Lexeme{
		d:       d,
		current: d.Start,
	}
}

func (l *Lexeme) Accepted() bool {
	return l.current.Final
}

func (l *Lexeme) Scan(ch rune) bool {
	for _, trans := range l.current.Transitions {
		if trans.Terminal.IsMatch(ch) {
			l.current = trans.Target
			return true
		}
	}
	return false
}
