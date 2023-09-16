package forest

import (
	"fmt"

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

func (s Symbol) String() string {
	return fmt.Sprintf("(%s, %d, %d)", s.Symbol.String(), s.Origin, s.Location)
}

func (s *Symbol) Accept(v Visitor) {
	v.VisitSymbol(s)
	for _, alt := range s.Internal.Alternatives {
		for _, child := range alt.Children {
			if acceptor, ok := child.(Acceptor); ok {
				acceptor.Accept(v)
			}
		}
	}
}

type Intermediate struct {
	Rule     *grammar.DottedRule
	Internal *Internal
	Origin   int
	Location int
}

func (Intermediate) node() {}

func (i Intermediate) String() string {
	return fmt.Sprintf("(%s, %d, %d)", i.Rule.String(), i.Origin, i.Location)
}

func (i *Intermediate) Accept(v Visitor) {
	v.VisitIntermediate(i)
	for _, alt := range i.Internal.Alternatives {
		for _, child := range alt.Children {
			if acceptor, ok := child.(Acceptor); ok {
				acceptor.Accept(v)
			}
		}
	}
}

type Token struct {
	Token    token.Token
	Origin   int
	Location int
}

func (Token) node() {}

func (t Token) String() string {
	return fmt.Sprintf("(%s, %d, %d)", t.Token.Type(), t.Origin, t.Location)
}

func (t *Token) Accept(v Visitor) {
	v.VisitToken(t)
}

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
