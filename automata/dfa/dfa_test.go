package dfa_test

import (
	"testing"

	"github.com/patrickhuber/go-earley/automata/dfa"
	"github.com/patrickhuber/go-earley/terminal"
	"github.com/stretchr/testify/require"
)

func TestDfa(t *testing.T) {
	start := &dfa.State{}
	end := &dfa.State{Final: true}
	start.Transitions = []dfa.Transition{
		{
			Target:   end,
			Terminal: terminal.NewCharacter('a'),
		},
	}
	d := dfa.NewDfa(start, "test")
	require.Equal(t, 1, len(d.Start.Transitions))
	require.True(t, start.Transitions[0].Terminal.IsMatch('a'))
}
