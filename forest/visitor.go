package forest

import (
	"fmt"
	"io"
)

type Visitor interface {
	VisitToken(t *Token) bool
	VisitIntermediate(i *Intermediate) bool
	VisitSymbol(s *Symbol) bool
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

func (p *Printer) VisitToken(t *Token) bool {
	if _, ok := p.cache[t]; ok {
		return false
	}
	p.cache[t] = struct{}{}
	return true
}

func (p *Printer) VisitIntermediate(i *Intermediate) bool {
	if _, ok := p.cache[i]; ok {
		return false
	}
	p.PrintIntermediate(i)
	p.PrintInternal(i)
	p.cache[i] = struct{}{}
	fmt.Fprintln(p.writer)
	return true
}

func (p *Printer) PrintIntermediate(i *Intermediate) {
	fmt.Fprintf(p.writer, "(%s, %d, %d)", i.Rule.String(), i.Origin(), i.Location())
}

func (p *Printer) VisitSymbol(s *Symbol) bool {
	if _, ok := p.cache[s]; ok {
		return false
	}
	p.PrintSymbol(s)
	p.PrintInternal(s)
	p.cache[s] = struct{}{}
	fmt.Fprintln(p.writer)
	return true
}

func (p *Printer) PrintSymbol(s *Symbol) {
	fmt.Fprintf(p.writer, "(%s, %d, %d)", s.Symbol.String(), s.Origin(), s.Location())
}

func (p *Printer) PrintInternal(i Internal) {
	fmt.Fprintf(p.writer, " -> ")
	for j, alt := range i.Alternatives() {
		if j > 0 {
			fmt.Fprintln(p.writer)
			fmt.Fprint(p.writer, "\t| ")
		}
		for _, child := range alt.Children() {
			switch c := child.(type) {
			case *Symbol:
				p.PrintSymbol(c)
			case *Intermediate:
				p.PrintIntermediate(c)
			case *Token:
				p.PrintToken(c)
			}
			fmt.Fprintf(p.writer, " ")
		}
	}
}

func (p *Printer) PrintToken(t *Token) {
	fmt.Fprintf(p.writer, "(%s, %d, %d)", t.Token.Type(), t.Origin(), t.Location())
}
