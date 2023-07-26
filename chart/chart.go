package chart

import (
	"github.com/patrickhuber/go-earley/ast"
	"github.com/patrickhuber/go-earley/grammar"
)

type Chart interface {
	Sets() []Set
	Contains(index int, ty StateType, rule grammar.DottedRule, origin int) bool
}

type Set interface {
	Predictions() []State
	Scans() []State
	Completions() []State
	Transitions() []Transition
	Location() int
	Enqueue(state State) bool
}

type State interface {
	DottedRule() grammar.DottedRule
	Origin() int
	Node() ast.Node
	Type() StateType
}

type Transition interface {
}

type StateType int

const (
	NormalStateType     StateType = 0
	TransitionStateType StateType = 1
)
