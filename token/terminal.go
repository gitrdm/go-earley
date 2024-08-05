package token

import "github.com/patrickhuber/go-earley/grammar"

type Terminal struct {
	rule     *grammar.TerminalLexerRule
	accepted bool
}

func NewTerminal(lexerRule *grammar.TerminalLexerRule, position int) *Terminal {
	return &Terminal{
		rule:     lexerRule,
		accepted: false,
	}
}
func (t *Terminal) Accepted() bool {
	return t.accepted
}

func (t *Terminal) Position() int {
	if t.accepted {
		return 1
	}
	return 0
}

func (t *Terminal) Reset(offset int) {
	t.accepted = false
}

func (t *Terminal) Scan(ch rune) bool {
	if t.accepted {
		return false
	}
	if !t.rule.CanApply(ch) {
		return false
	}
	t.accepted = true
	return true
}
