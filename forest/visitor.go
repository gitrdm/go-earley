package forest

import (
	"fmt"
	"io"
)

type Visitor interface {
	VisitToken(t *Token)
	VisitIntermediate(i *Intermediate)
	VisitSymbol(s *Symbol)
}

type Acceptor interface {
	Accept(v Visitor)
}

type Printer struct {
	writer io.Writer
	cache  map[Node]struct{}
}

func NewPrinter(writer io.Writer) *Printer {
	return &Printer{
		writer: writer,
		cache:  map[Node]struct{}{},
	}
}

func (p *Printer) VisitToken(t *Token) {
	if _, ok := p.cache[t]; ok {
		return
	}
	fmt.Fprintf(p.writer, "(%s, %d, %d)", t.Token.Type(), t.Origin(), t.Location())
	p.cache[t] = struct{}{}
}

func (p *Printer) VisitIntermediate(i *Intermediate) {
	if _, ok := p.cache[i]; ok {
		return
	}
	p.PrintIntermediate(i)
	p.PrintInternal(i)
	p.cache[i] = struct{}{}
}

func (p *Printer) PrintIntermediate(i *Intermediate) {
	fmt.Fprintf(p.writer, "(%s, %d, %d)", i.Rule.String(), i.Origin(), i.Location())
}

func (p *Printer) VisitSymbol(s *Symbol) {
	if _, ok := p.cache[s]; ok {
		return
	}
	p.PrintSymbol(s)
	p.PrintInternal(s)
	p.cache[s] = struct{}{}
}

func (p *Printer) PrintSymbol(s *Symbol) {
	fmt.Fprintf(p.writer, "(%s, %d, %d)", s.Symbol.String(), s.Origin(), s.Location())
}

func (p *Printer) PrintInternal(i Internal) {
	fmt.Fprintf(p.writer, " -> ")
	for i, alt := range i.Alternatives() {
		if i > 0 {
			fmt.Fprint(p.writer, "\t| ")
		}
		for _, child := range alt.Children() {
			switch c := child.(type) {
			case *Symbol:
				p.PrintSymbol(c)
			case *Intermediate:
				p.PrintIntermediate(c)
			case *Token:
				p.VisitToken(c)
			}
			fmt.Fprintf(p.writer, " ")
		}
		fmt.Fprintln(p.writer)
	}
}
