package token

import (
	"fmt"

	"github.com/patrickhuber/go-collections/generic/queue"
	"github.com/patrickhuber/go-earley/grammar"
)

type TerminalFactory struct {
	queue queue.Queue[*Terminal]
}

func NewTerminalFactory() *TerminalFactory {
	return &TerminalFactory{
		queue: queue.New[*Terminal](),
	}
}

// Create implements Factory.
func (f *TerminalFactory) Create(lexerRule grammar.LexerRule, str string, position int) (Lexeme, error) {
	rule, ok := lexerRule.(*grammar.TerminalLexerRule)
	if !ok || lexerRule.LexerRuleType() != grammar.TerminalLexerRuleType {
		return nil, fmt.Errorf("terminal factory expected lexer rule of type %s but found %s", grammar.TerminalLexerRuleType, lexerRule.LexerRuleType())
	}
	if f.queue.Length() == 0 {
		return NewTerminal(rule, position), nil
	}
	reused := f.queue.Dequeue()
	reused.Reset(position)
	return reused, nil
}

// Free implements Factory.
func (f *TerminalFactory) Free(lexeme Lexeme) error {
	t, ok := lexeme.(*Terminal)
	if !ok {
		return fmt.Errorf("Free expected *lexeme.Terminal but found %T", lexeme)
	}
	f.queue.Enqueue(t)
	return nil
}

// Type implements Factory.
func (f *TerminalFactory) Type() string {
	return grammar.TerminalLexerRuleType
}

func NewFactory() Factory {
	return &TerminalFactory{
		queue: queue.New[*Terminal](),
	}
}
