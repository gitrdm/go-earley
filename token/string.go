package token

import (
	"unicode/utf8"

	"github.com/patrickhuber/go-earley/grammar"
)

type String struct {
	position int
	rule     *grammar.StringLexerRule
}

func NewString(lexerRule *grammar.StringLexerRule, position int) *String {
	return &String{
		rule:     lexerRule,
		position: position,
	}
}

// Accepted implements Lexeme.
func (s *String) Accepted() bool {
	return s.position == len(s.rule.Value)
}

// Position implements Lexeme.
func (s *String) Position() int {
	return s.position
}

// Reset implements Lexeme.
func (s *String) Reset(offset int) {
	s.position = -1
}

// Scan implements Lexeme.
func (s *String) Scan(ch rune) bool {
	if s.position >= len(s.rule.Value) {
		return false
	}
	r, n := utf8.DecodeRuneInString(s.rule.Value[s.position:])
	if ch != r {
		return false
	}
	s.position += n
	return true
}

// Type implements Token.
func (s *String) TokenType() string {
	return s.rule.TokenType()
}

func (s *String) LexerRule() grammar.LexerRule {
	return s.rule
}
