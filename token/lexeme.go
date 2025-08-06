package token

import "github.com/patrickhuber/go-earley/grammar"

// Lexeme is a mutable token
type Lexeme interface {
	Token
	Scan(ch rune) bool
	Accepted() bool
	LexerRule() grammar.LexerRule
}
