package grammar_test

import (
	"testing"

	"github.com/patrickhuber/go-earley/grammar"
	"github.com/stretchr/testify/require"
)

func TestProduction(t *testing.T) {

	t.Run("left_hand_side", func(t *testing.T) {
		RunProductionTest(func(production *grammar.Production) {
			symbol := production.LeftHandSide
			require.NotNil(t, symbol)
		})
	})

	t.Run("right_hand_side", func(t *testing.T) {
		RunProductionTest(func(production *grammar.Production) {
			rhs := production.RightHandSide
			require.NotNil(t, rhs)
			require.Equal(t, 2, len(rhs))
			for _, sym := range rhs {
				require.NotNil(t, sym)
			}
		})
	})
}

func RunProductionTest(action func(production *grammar.Production)) {
	leftHandSide := grammar.NewNonTerminal("S")
	rightHandSide := []grammar.Symbol{
		grammar.NewNonTerminal("S"),
		grammar.NewStringLexerRule("S"),
	}
	production := grammar.NewProduction(leftHandSide, rightHandSide...)
	action(production)
}
