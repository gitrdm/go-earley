package dfa

import "github.com/patrickhuber/go-earley/grammar"

type Lexeme struct {
	dfa      *Dfa
	current  *State
	position int
}

// NewLexeme creates a new Lexeme for the given DFA and position.
func NewLexeme(dfa *Dfa, position int) *Lexeme {
	return &Lexeme{
		dfa:      dfa,
		current:  dfa.Start,
		position: position,
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

func (l *Lexeme) LexerRule() grammar.LexerRule {
	return l.dfa
}

func (l *Lexeme) Position() int {
	return l.position
}

func (l *Lexeme) TokenType() string {
	return l.dfa.TokenType()
}
