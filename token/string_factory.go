package token

import (
	"fmt"

	"github.com/patrickhuber/go-collections/generic/queue"
	"github.com/patrickhuber/go-earley/grammar"
	"github.com/patrickhuber/go-earley/lexrule"
)

type stringFactory struct {
	queue queue.Queue[*stringToken]
}

// Create implements Factory.
func (f *stringFactory) Create(lexerRule grammar.LexerRule, span AccumulatorSpan) (Token, error) {
	rule, ok := lexerRule.(lexrule.String)
	if !ok || lexerRule.Type() != lexrule.StringType {
		return nil, fmt.Errorf("string factory expected lexer rule of type %s but found %s", lexrule.StringType, lexerRule.Type())
	}
	if f.queue.Length() == 0 {
		return NewString(rule, span), nil
	}
	reused := f.queue.Dequeue()
	reused.Reset()
	return reused, nil
}

func (f *stringFactory) Free(token Token) error {
	return nil
}

func (f *stringFactory) Type() string {
	return lexrule.StringType
}

func NewStringFactory() Factory {
	return &stringFactory{
		queue: queue.New[*stringToken](),
	}
}
