package dfa_test

import (
	"testing"

	"github.com/patrickhuber/go-earley/automata/dfa"
	"github.com/patrickhuber/go-earley/terminal"
)

func TestState(t *testing.T) {
	t.Run("is match", func(t *testing.T) {
		zero := &dfa.State{
			Final: false,
		}
		one := &dfa.State{
			Final: true,
		}
		zero.Transitions = append(zero.Transitions, dfa.Transition{
			Target:   one,
			Terminal: terminal.NewNumber(),
		})
		if !zero.IsMatch('1') {
			t.Fatalf("expected number to match")
		}
		if zero.IsMatch('a') {
			t.Fatalf("expected letter not to match")
		}
	})
}
