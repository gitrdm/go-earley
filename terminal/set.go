package terminal

import (
	"strings"

	"github.com/patrickhuber/go-earley/grammar"
)

type Set struct {
	grammar.SymbolImpl
	Terminals []grammar.Terminal
}

func (s *Set) IsMatch(ch rune) bool {
	for _, t := range s.Terminals {
		if t.IsMatch(ch) {
			return true
		}
	}
	return false
}

func NewSet(terminals []grammar.Terminal) grammar.Terminal {
	return &Set{
		Terminals: terminals,
	}
}

func (s *Set) String() string {
	var builder strings.Builder
	builder.WriteRune('[')
	for _, terminal := range s.Terminals {
		builder.WriteString(terminal.String())
	}
	builder.WriteRune(']')
	return builder.String()
}
