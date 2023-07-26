package lexeme

import (
	"github.com/patrickhuber/go-earley/capture"
	"github.com/patrickhuber/go-earley/grammar"
)

type Factory interface {
	Type() string
	Create(lexerRule grammar.LexerRule, cap capture.Capture, offset int) (Lexeme, error)
	Free(lexeme Lexeme) error
}
