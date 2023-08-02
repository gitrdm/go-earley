package chart

import (
	"github.com/patrickhuber/go-earley/grammar"
	"github.com/patrickhuber/go-earley/state"
)

type Set struct {
	Predictions []*state.Normal
	Scans       []*state.Normal
	Completions []*state.Normal
	Transitions []state.Transition
	Location    int

	reductions map[grammar.Symbol][]*state.Normal
}

func NewSet() *Set {
	return &Set{}
}

func (s *Set) Contains(ty state.Type, dr *grammar.DottedRule, origin int) bool {
	_, ok := s.find(ty, dr, origin)
	return ok
}

func (s *Set) GetOrCreate(ty state.Type, dr *grammar.DottedRule, origin int) state.State {
	st, ok := s.find(ty, dr, origin)
	if ok {
		return st
	}
	normal := &state.Normal{
		Origin:     origin,
		DottedRule: dr,
	}
	s.enqueueNormal(normal)
	return normal
}

func (s *Set) find(ty state.Type, dr *grammar.DottedRule, origin int) (state.State, bool) {
	if ty != state.NormalType {
		return nil, false
	}
	if dr.Complete() {
		return s.FindCompletion(dr, origin)
	}
	current := dr.PostDotSymbol()
	_, ok := current.(grammar.NonTerminal)
	if ok {
		return s.FindPrediction(dr, origin)
	}
	return s.FindScan(dr, origin)
}

func (s *Set) findIn(dr *grammar.DottedRule, origin int, states []*state.Normal) (state.State, bool) {
	for _, state := range states {
		if state.Origin != origin {
			continue
		}
		if state.DottedRule != dr {
			continue
		}
		return state, true
	}
	return nil, false
}

func (s *Set) FindCompletion(dr *grammar.DottedRule, origin int) (state.State, bool) {
	return s.findIn(dr, origin, s.Completions)
}

func (s *Set) FindPrediction(dr *grammar.DottedRule, origin int) (state.State, bool) {
	return s.findIn(dr, origin, s.Predictions)
}

func (s *Set) FindScan(dr *grammar.DottedRule, origin int) (state.State, bool) {
	return s.findIn(dr, origin, s.Scans)
}

func (s *Set) FindTransition(sym grammar.Symbol) (*state.Transition, bool) {
	return nil, false
}

// FindSourceStates finds any state in the set where the post dot symbol is equal to the given symbol sym
// ex:
// A ->  B.C
// A ->  B B.
// B -> .C
// C -> .c
//
// Given C, B -> .C and A-> B.C are returned
func (s *Set) FindSourceStates(sym grammar.Symbol) []*state.Normal {
	var states []*state.Normal
	if sym == nil {
		return states
	}
	for _, p := range s.Predictions {
		postDot := p.DottedRule.PostDotSymbol()
		if postDot == sym {
			states = append(states, p)
		}
	}
	return states
}

// FindReductions returns all completions where the left hand symbol is the same as the search symbol
func (s *Set) FindReductions(sym grammar.Symbol) []*state.Normal {
	var states []*state.Normal
	if sym == nil {
		return states
	}
	reductions, ok := s.reductions[sym]
	if !ok {
		return states
	}
	return reductions
}

// Enqueue implements Set.
func (s *Set) Enqueue(st state.State) bool {
	if st.Type() == state.NormalType {
		return s.enqueueNormal(st.(*state.Normal))
	} else {
		return s.enqueueTransition(st.(*state.Transition))
	}
}

func (s *Set) enqueueNormal(st *state.Normal) bool {
	rule := st.DottedRule
	if rule.Complete() {
		return s.addUniqueCompletion(st)
	}
	postDot := rule.PostDotSymbol()
	if _, ok := postDot.(grammar.NonTerminal); ok {
		return s.addUniquePrediction(st)
	}
	return s.addUniqueScan(st)
}

func (s *Set) enqueueTransition(st *state.Transition) bool {
	return false
}

func (s *Set) addUniqueCompletion(completion *state.Normal) bool {
	completions, ok := addUnique(s.Completions, completion)
	if !ok {
		return false
	}
	s.Completions = completions
	return true
}

func (s *Set) addUniquePrediction(prediction *state.Normal) bool {
	predictions, ok := addUnique(s.Predictions, prediction)
	if !ok {
		return false
	}
	s.Predictions = predictions
	return true
}

func (s *Set) addUniqueScan(scan *state.Normal) bool {
	scans, ok := addUnique(s.Scans, scan)
	if !ok {
		return false
	}
	s.Scans = scans
	return true
}

func addUnique(states []*state.Normal, st *state.Normal) ([]*state.Normal, bool) {
	for _, s := range states {
		if s == st {
			return states, false
		}
	}
	states = append(states, st)
	return states, true
}
