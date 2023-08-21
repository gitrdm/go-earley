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
}
