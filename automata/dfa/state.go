package dfa

import "github.com/patrickhuber/go-earley/grammar"

type State struct {
	Final       bool
	Transitions []Transition
}

type Transition struct {
	Target   *State
	Terminal grammar.Terminal
}
