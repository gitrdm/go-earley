package dfa_test

import (
	"testing"

	"github.com/patrickhuber/go-earley/automata/dfa"
	"github.com/patrickhuber/go-earley/terminal"
)

func TestLexRule(t *testing.T) {
	t.Run("can apply single letter", func(t *testing.T) {
		zero := &dfa.State{
			Final: false,
		}
		one := &dfa.State{
			Final: true,
		}
		zero.Transitions = append(zero.Transitions, dfa.Transition{
			Terminal: terminal.NewLetter(),
			Target:   one,
		})
		lexRule := dfa.Dfa{
			Start: zero,
		}
		if !lexRule.CanApply('a') {
			t.Fatal("should apply on 'a'")
		}
	})
}
