package parser

import "github.com/patrickhuber/go-earley/grammar"

type Parser interface {
	Expected() []grammar.LexerRule
	Accepted() bool
}

func New(g grammar.Grammar) Parser {
	return nil
}
