package forest

import (
	"fmt"

	"github.com/patrickhuber/go-earley/grammar"
)

type Symbol struct {
	Symbol   grammar.Symbol
	internal *internal
	origin   int
	location int
	expanded bool
	paths    map[Path]Node
}

type Path interface {
	Next() Path
	Node() Node
}

func NewSymbol(sym grammar.Symbol, origin int, location int, alternatives ...Group) *Symbol {
	return &Symbol{
		Symbol:   sym,
		origin:   origin,
		location: location,
		internal: &internal{
			alternatives: alternatives,
		},
	}
}

func (Symbol) node()           {}
func (s Symbol) Origin() int   { return s.origin }
func (s Symbol) Location() int { return s.location }
func (s *Symbol) Alternatives() []Group {
	if !s.expanded {
		s.expand()
		s.expanded = true
	}
	return s.internal.alternatives
}

func (s *Symbol) expand() {
	for path, node := range s.paths {
		leftTree := path.Node()
		rightSubTree := node
		next := path.Next()

		if next.Node() == nil || next.Node().Location() == rightSubTree.Location() {
			s.internal.alternatives = append(s.internal.alternatives, group{
				children: []Node{
					leftTree,
					rightSubTree,
				},
			})
			return
		}
		rightTree := NewSymbol(s.Symbol, next.Node().Origin(), s.location)
		rightTree.AddPath(path.Next(), node)
		s.internal.alternatives = append(s.internal.alternatives, &group{
			children: []Node{
				leftTree,
				rightTree,
			},
		})
	}
}

func (s *Symbol) AddUniqueFamily(w, v Node) {
	s.internal.AddUniqueFamily(w, v)
}

func (s *Symbol) AddPath(path Path, node Node) {
	if s.paths == nil {
		s.paths = make(map[Path]Node)
	}
	s.paths[path] = node
}

func (s Symbol) String() string {
	return fmt.Sprintf("(%s, %d, %d)", s.Symbol.String(), s.origin, s.location)
}

func (s *Symbol) Accept(v Visitor) {
	v.VisitSymbol(s)
	for _, alt := range s.internal.alternatives {
		for _, child := range alt.Children() {
			if acceptor, ok := child.(Acceptor); ok {
				acceptor.Accept(v)
			}
		}
	}
}
