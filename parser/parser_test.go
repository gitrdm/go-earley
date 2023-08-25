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
		// S -> S S
		grammar.NewProduction(S, S, S),
		// S -> S
		grammar.NewProduction(S, S),
		// S -> 's'
		grammar.NewProduction(S, s),
	)

	p := parser.New(g)

	for i := 0; i < 10; i++ {
		tok := token.FromString("s", i, s.TokenType())
		ok, err := p.Pulse(tok)
		require.NoError(t, err, "loop %d", i)
		require.True(t, ok, "loop %d", i)
	}
	require.True(t, p.Accepted())
}

func TestAycockHorspool(t *testing.T) {
	/*
		S' -> S
		S  -> A A A A
		A  -> a | E
		E  ->
	*/
	SPrime := grammar.NewNonTerminal("S'")
	S := grammar.NewNonTerminal("S")
	A := grammar.NewNonTerminal("A")
	E := grammar.NewNonTerminal("E")
	a := lexrule.NewString("a")

	g := grammar.New(SPrime,
		grammar.NewProduction(SPrime, S),
		grammar.NewProduction(S, A, A, A, A),
		grammar.NewProduction(A, a),
		grammar.NewProduction(A, E),
		grammar.NewProduction(E),
	)

	p := parser.New(g)
	tok := token.FromString("a", 0, a.TokenType())
	ok, err := p.Pulse(tok)

	require.NoError(t, err)
	require.True(t, ok)
	require.True(t, p.Accepted())
}

func TestLeo(t *testing.T) {
	t.Run("leo 1", func(t *testing.T) {
		A := grammar.NewNonTerminal("A")
		a := lexrule.NewString("a")

		// A -> A 'a'
		// A ->
		g := grammar.New(A,
			grammar.NewProduction(A, a, A),
			grammar.NewProduction(A),
		)

		p := parser.New(g)
		for i := 0; i < 10; i++ {
			tok := token.FromString("a", i, a.TokenType())
			ok, err := p.Pulse(tok)
			require.NoError(t, err, "loop %d", i)
			require.True(t, ok, "loop %d", i)
		}
		require.True(t, p.Accepted())
	})

	t.Run("leo 2", func(t *testing.T) {
		// S -> 'a' S
		// S -> C
		// C -> 'a' C 'b'
		// C ->
		S := grammar.NewNonTerminal("S")
		C := grammar.NewNonTerminal("C")
		a := lexrule.NewString("a")
		b := lexrule.NewString("b")

		g := grammar.New(S,
			grammar.NewProduction(S, a, S),
			grammar.NewProduction(S, C),
			grammar.NewProduction(C, a, C, b),
			grammar.NewProduction(C))

		p := parser.New(g)
		for i := 0; i < 10; i++ {
			lr := a
			if i > 5 {
				lr = b
			}
			tok := token.FromString(lr.Value, i, lr.TokenType())
			ok, err := p.Pulse(tok)
			require.NoError(t, err)
			require.True(t, ok)
		}
		require.True(t, p.Accepted())
	})
}
