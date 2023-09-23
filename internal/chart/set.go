package chart

import (
	"github.com/patrickhuber/go-earley/grammar"
	"github.com/patrickhuber/go-earley/internal/state"
)

type Set struct {
	Predictions []*state.Normal
	Scans       []*state.Normal
	Completions []*state.Normal
	Transitions map[grammar.Symbol]*state.Transition
	Location    int

	reductions map[grammar.Symbol][]*state.Normal
}

func NewSet() *Set {
	return &Set{}
}

func (s *Set) Contains(ty state.Type, dr *grammar.DottedRule, origin int) bool {
	_, ok := s.find(dr, origin)
	return ok
}

func (s *Set) GetOrCreate(dr *grammar.DottedRule, origin int) *state.Normal {
	st, ok := s.find(dr, origin)
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

func (s *Set) find(dr *grammar.DottedRule, origin int) (*state.Normal, bool) {
	if dr.Complete() {
		return s.FindCompletion(dr, origin)
	}
	postDot, ok := dr.PostDotSymbol().Deconstruct()
	if ok {
		if _, ok := postDot.(grammar.NonTerminal); ok {
			return s.FindPrediction(dr, origin)
		}
	}
	return s.FindScan(dr, origin)
}

func (s *Set) findIn(dr *grammar.DottedRule, origin int, states []*state.Normal) (*state.Normal, bool) {
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

func (s *Set) FindCompletion(dr *grammar.DottedRule, origin int) (*state.Normal, bool) {
	return s.findIn(dr, origin, s.Completions)
}

func (s *Set) FindPrediction(dr *grammar.DottedRule, origin int) (*state.Normal, bool) {
	return s.findIn(dr, origin, s.Predictions)
}

func (s *Set) FindScan(dr *grammar.DottedRule, origin int) (*state.Normal, bool) {
	return s.findIn(dr, origin, s.Scans)
}

func (s *Set) FindTransition(sym grammar.Symbol) (*state.Transition, bool) {
	trans, ok := s.Transitions[sym]
	return trans, ok
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
		postDot, ok := p.DottedRule.PostDotSymbol().Deconstruct()
		if !ok {
			continue
		}
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
	postDot, ok := rule.PostDotSymbol().Deconstruct()
	if ok {
		if _, ok := postDot.(grammar.NonTerminal); ok {
			return s.addUniquePrediction(st)
		}
	}
	return s.addUniqueScan(st)
}

func (s *Set) enqueueTransition(st *state.Transition) bool {
	_, ok := s.FindTransition(st.Symbol)
	if ok {
		return false
	}
	// create if not exists
	if s.Transitions == nil {
		s.Transitions = map[grammar.Symbol]*state.Transition{}
	}
	s.Transitions[st.Symbol] = st
	return true
}

func (s *Set) addUniqueCompletion(completion *state.Normal) bool {
	completions, ok := addUnique(s.Completions, completion)
	if !ok {
		return false
	}
	s.Completions = completions

	sym := completion.DottedRule.Production.LeftHandSide
	if s.reductions == nil {
		s.reductions = make(map[grammar.Symbol][]*state.Normal)
	}
	list, ok := s.reductions[sym]
	if !ok {
		list = []*state.Normal{}
	}
	list = append(list, completion)
	s.reductions[sym] = list
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

func addUnique(states []*state.Normal, item *state.Normal) ([]*state.Normal, bool) {
	if contains(states, item) {
		return states, false
	}
	states = append(states, item)
	return states, true
}

func contains(states []*state.Normal, item *state.Normal) bool {
	found := false
	for _, s := range states {
		if equal(s, item) {
			found = true
			break
		}
	}
	return found
}

func equal(n1 *state.Normal, n2 *state.Normal) bool {
	if n1.Origin != n2.Origin {
		return false
	}
	return n1.DottedRule.String() == n2.DottedRule.String()
}
