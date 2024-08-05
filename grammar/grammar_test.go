package grammar_test

import (
	"testing"

	"github.com/patrickhuber/go-earley/grammar"
	"github.com/stretchr/testify/require"
)

func TestGrammar(t *testing.T) {
	t.Run("transitive null", func(t *testing.T) {
		S := grammar.NewNonTerminal("S")
		A := grammar.NewNonTerminal("A")
		E := grammar.NewNonTerminal("E")

		g := grammar.New(S, grammar.NewProduction(S, A),
			grammar.NewProduction(A, E),
			grammar.NewProduction(E))
		require.True(t, g.IsTransativeNullable(S))
		require.True(t, g.IsTransativeNullable(A))
		require.True(t, g.IsTransativeNullable(E))
	})
	
	t.Run("right recursive", func(t *testing.T) {
		A := grammar.NewNonTerminal("A")
		a := grammar.NewStringLexerRule("a")

		// A -> A 'a'
		A_aA := grammar.NewProduction(A, a, A)
		// A ->
		A_ := grammar.NewProduction(A)

		g := grammar.New(A,
			A_aA,
			A_)

		require.True(t, g.IsRightRecursive(A_aA))
		require.False(t, g.IsRightRecursive(A_))

		A_Aa := grammar.NewProduction(A, A, a)
		g = grammar.New(A,
			A_Aa,
			A_)
		require.False(t, g.IsRightRecursive(A_Aa))
		require.False(t, g.IsRightRecursive(A_))
	})
}
