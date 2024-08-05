package token

import (
	"fmt"

	"github.com/patrickhuber/go-collections/generic/queue"
	"github.com/patrickhuber/go-earley/grammar"
)

type stringFactory struct {
	queue queue.Queue[*String]
}

// Create implements Factory.
func (f *stringFactory) Create(lexerRule grammar.LexerRule, str string, position int) (Lexeme, error) {
	rule, ok := lexerRule.(*grammar.StringLexerRule)
	if !ok || lexerRule.LexerRuleType() != grammar.StringLexerRuleType {
		return nil, fmt.Errorf("string factory expected lexer rule of type %s but found %s", grammar.StringLexerRuleType, lexerRule.LexerRuleType())
	}
	if f.queue.Length() == 0 {
		return NewString(rule, position), nil
	}
	reused := f.queue.Dequeue()
	reused.Reset(position)
	return reused, nil
}

func (f *stringFactory) Free(lexeme Lexeme) error {
	s, ok := lexeme.(*String)
	if !ok {
		return fmt.Errorf("Free expected *lexeme.String but found %T", lexeme)
	}
	f.queue.Enqueue(s)
	return nil
}

func (f *stringFactory) Type() string {
	return grammar.StringLexerRuleType
}

func NewStringFactory() Factory {
	return &stringFactory{
		queue: queue.New[*String](),
	}
}
