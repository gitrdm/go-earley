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
