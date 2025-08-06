package dfa

import (
	"fmt"

	"github.com/patrickhuber/go-collections/generic/queue"
	"github.com/patrickhuber/go-earley/grammar"
	"github.com/patrickhuber/go-earley/token"
)

type Factory struct {
	queue queue.Queue[*Lexeme]
}

// Create implements token.Factory.
func (f *Factory) Create(lexerRule grammar.LexerRule, str string, offset int) (token.Lexeme, error) {
	rule, ok := lexerRule.(*Dfa)
	if !ok || lexerRule.LexerRuleType() != LexerRuleType {
		return nil, fmt.Errorf("dfa factory expected lexer rule of type %s but found %s", LexerRuleType, lexerRule.LexerRuleType())
	}
	if f.queue.Length() == 0 {
		return NewLexeme(rule, offset), nil
	}
	panic("unimplemented")
}

// Free implements token.Factory.
func (f *Factory) Free(lexeme token.Lexeme) error {
	panic("unimplemented")
}

// Type implements token.Factory.
func (f *Factory) Type() string {
	panic("unimplemented")
}

func NewFactory() token.Factory {
	return &Factory{}
}
