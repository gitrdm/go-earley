package nfa_test

import (
	"testing"

	"github.com/patrickhuber/go-earley/automata/nfa"
	"github.com/patrickhuber/go-earley/terminal"
	"github.com/stretchr/testify/require"
)

func TestNfa(t *testing.T) {
	start := &nfa.State{}
	end := &nfa.State{}
	start.Transitions =
		[]nfa.Transition{
			nfa.NewNull(start),
			nfa.NewTerminal(terminal.NewCharacter('a'), end),
		}
	n := nfa.Nfa{
		Start: start,
	}
	require.Equal(t, 2, len(n.Start.Transitions))
}
