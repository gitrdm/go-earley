package chart

import (
	"github.com/patrickhuber/go-earley/grammar"
	"github.com/patrickhuber/go-earley/state"
)

type Set struct {
	Predictions []state.State
	Scans       []state.State
	Completions []state.State
	Transitions []state.Transition
	Location    int
}

func NewSet() *Set {
	return &Set{}
}

func (s *Set) Contains(ty state.Type, dr *grammar.DottedRule, origin int) bool {
	if ty != state.NormalType {
		return false
	}
	if dr.Complete() {
		_, ok := s.FindCompletion(dr, origin)
		if ok {
			return ok
		}
	}
	current := dr.PostDotSymbol()
	_, ok := current.(grammar.NonTerminal)
	if ok {
		_, ok := s.FindPrediction(dr, origin)
		if ok {
			return ok
		}
	}

	_, ok = s.FindScan(dr, origin)
	return ok
}

func (s *Set) FindCompletion(dr *grammar.DottedRule, origin int) (state.State, bool) {
	return nil, false
}

func (s *Set) FindPrediction(dr *grammar.DottedRule, origin int) (state.State, bool) {
	return nil, false
}

func (s *Set) FindScan(dr *grammar.DottedRule, origin int) (state.State, bool) {
	return nil, false
}

// Enqueue implements Set.
func (s *Set) Enqueue(st state.State) bool {
	panic("unimplemented")
}
