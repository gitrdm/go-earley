package forest

import (
	"github.com/patrickhuber/go-earley/grammar"
	"github.com/patrickhuber/go-earley/token"
)

type Node interface {
	node()
}

type Symbol struct {
	Symbol   grammar.Symbol
	Internal *Internal
	Origin   int
	Location int
}

func (Symbol) node() {}

type Intermediate struct {
	Rule     *grammar.DottedRule
	Internal *Internal
	Origin   int
	Location int
}

func (Intermediate) node() {}

type Token struct {
	Token    token.Token
	Origin   int
	Location int
}

func (Token) node() {}

type Group struct {
	Children []Node
}

type Internal struct {
	Alternatives []*Group
}

func (i *Internal) AddUniqueFamily(w, v Node) {
	childCount := 1
	if v != nil {
		childCount += 1
	}
	for _, group := range i.Alternatives {

		if len(group.Children) != childCount {
			continue
		}
		if i.isMatchedSubtree(w, v, group) {
			return
		}
	}

	group := &Group{}
	group.Children = append(group.Children, w)
	if childCount > 1 {
		group.Children = append(group.Children, v)
	}
	i.Alternatives = append(i.Alternatives, group)
}

func (i *Internal) isMatchedSubtree(first, second Node, group *Group) bool {

	firstCompare := group.Children[0]

	if first != firstCompare {
		return false
	}

	if second == nil {
		return true
	}

	secondCompare := group.Children[1]

	return secondCompare == second
}
