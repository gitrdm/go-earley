package dfa

import "github.com/patrickhuber/go-earley/grammar"

type State struct {
	Final       bool
	Transitions []Transition
}

func (s *State) IsMatch(ch rune) bool {
	for _, t := range s.Transitions {
		if t.Terminal.IsMatch(ch) {
			return true
		}
	}
	return false
}

type Transition struct {
	Target   *State
	Terminal grammar.Terminal
}
