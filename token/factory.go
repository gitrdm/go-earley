package token

import (
	"github.com/patrickhuber/go-earley/grammar"
)

type Factory interface {
	Type() string
	Create(lexerRule grammar.LexerRule, span AccumulatorSpan) (Token, error)
	Free(token Token) error
}
