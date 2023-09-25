package forest

import (
	"fmt"

	"github.com/patrickhuber/go-earley/grammar"
)

type Intermediate struct {
	Rule     *grammar.DottedRule
	internal *internal
	origin   int
	location int
}

func NewIntermediate(rule *grammar.DottedRule, origin int, location int, alternatives ...Group) *Intermediate {
	return &Intermediate{
		Rule:     rule,
		origin:   origin,
		location: location,
		internal: &internal{
			alternatives: alternatives,
		},
	}
}

func (Intermediate) node()           {}
func (i Intermediate) Origin() int   { return i.origin }
func (i Intermediate) Location() int { return i.location }
func (i Intermediate) Alternatives() []Group {
	return i.internal.alternatives
}

func (i *Intermediate) AddUniqueFamily(w, v Node) {
	i.internal.AddUniqueFamily(w, v)
}

func (i Intermediate) String() string {
	return fmt.Sprintf("(%s, %d, %d)", i.Rule.String(), i.origin, i.location)
}

func (i *Intermediate) Accept(v Visitor) {
	if !v.VisitIntermediate(i) {
		// already visited
		return
	}
	for _, alt := range i.internal.alternatives {
		for _, child := range alt.Children() {
			if acceptor, ok := child.(Acceptor); ok {
				acceptor.Accept(v)
			}
		}
	}
}
