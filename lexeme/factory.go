package lexeme

import (
	"github.com/patrickhuber/go-earley/grammar"
)

type Factory interface {
	Type() string
	Create(lexerRule grammar.LexerRule, str string, offset int) (Lexeme, error)
	Free(lexeme Lexeme) error
}
