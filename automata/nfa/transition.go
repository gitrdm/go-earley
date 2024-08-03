package nfa

import "github.com/patrickhuber/go-earley/grammar"

type Transition interface {
	transition()
	Target() *State
}

func NewTerminal(terminal grammar.Terminal, target *State) *Terminal {
	return &Terminal{
		terminal: terminal,
		target:   target,
	}
}

type Terminal struct {
	terminal grammar.Terminal
	target   *State
}

func (Terminal) transition() {}

func (t Terminal) Target() *State {
	return t.target
}

func NewNull(target *State) *Null {
	return &Null{
		target: target,
	}
}

type Null struct {
	target *State
}

func (Null) transition() {}

func (n Null) Target() *State {
	return n.target
}
