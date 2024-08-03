package lexeme

import (
	"fmt"

	"github.com/patrickhuber/go-collections/generic/queue"
	"github.com/patrickhuber/go-earley/grammar"
	"github.com/patrickhuber/go-earley/lexrule"
)

type stringFactory struct {
	queue queue.Queue[*String]
}

// Create implements Factory.
func (f *stringFactory) Create(lexerRule grammar.LexerRule, str string, offset int) (Lexeme, error) {
	rule, ok := lexerRule.(*lexrule.String)
	if !ok || lexerRule.Type() != lexrule.StringType {
		return nil, fmt.Errorf("string factory expected lexer rule of type %s but found %s", lexrule.StringType, lexerRule.Type())
	}
	if f.queue.Length() == 0 {
		return NewString(rule, offset), nil
	}
	reused := f.queue.Dequeue()
	reused.Reset(offset)
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
	return lexrule.StringType
}

func NewStringFactory() Factory {
	return &stringFactory{
		queue: queue.New[*String](),
	}
}
