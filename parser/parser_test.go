package parser_test

import (
	"testing"

	"github.com/patrickhuber/go-earley/grammar"
	"github.com/patrickhuber/go-earley/lexrule"
	"github.com/patrickhuber/go-earley/parser"
	"github.com/patrickhuber/go-earley/token"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	S := grammar.NewNonTerminal("S")
	s := lexrule.NewString("s")

	// a series of S's
	// S -> S S | S | 's'
	g := grammar.New(
		S,
		grammar.NewProduction(S, S, S),
		grammar.NewProduction(S, S),
		grammar.NewProduction(S, s),
	)

	p := parser.New(g)

	for i := 0; i < 10; i++ {
		tok := token.FromString("s", i, s.Type())
		ok, err := p.Pulse(tok)
		require.NoError(t, err, "loop %d", i)
		require.True(t, ok, "loop %d", i)
	}
}
