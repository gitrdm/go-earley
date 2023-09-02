package forest

import (
	"github.com/patrickhuber/go-earley/grammar"
	"github.com/patrickhuber/go-earley/token"
)

type Set struct {
	Symbols       []*Symbol
	Intermediates []*Intermediate
	Tokens        []*Token
}

func (s *Set) AddOrGetExistingSymbolNode(
	sym grammar.Symbol,
	origin int,
	location int) *Symbol {
	for _, symbol := range s.Symbols {
		if symbol.Symbol != sym {
			continue
		}
		if symbol.Origin != origin {
			continue
		}
		if symbol.Location != location {
			continue
		}
		return symbol
	}
	// not found, so create it
	symbol := &Symbol{
		Symbol:   sym,
		Origin:   origin,
		Location: location,
		Internal: &Internal{},
	}
	s.Symbols = append(s.Symbols, symbol)
	return symbol
}

func (s *Set) AddOrGetExistingIntermediateNode(
	rule *grammar.DottedRule,
	origin,
	location int) *Intermediate {
	for _, intermediate := range s.Intermediates {
		if intermediate.Rule != rule {
			continue
		}
		if intermediate.Origin != origin {
			continue
		}
		if intermediate.Location != location {
			continue
		}
		return intermediate
	}
	intermediate := &Intermediate{
		Rule:     rule,
		Origin:   origin,
		Location: location,
		Internal: &Internal{},
	}
	s.Intermediates = append(s.Intermediates, intermediate)
	return intermediate
}

func (s *Set) AddOrGetExistingTokenNode(tok token.Token, location int) *Token {
	for _, token := range s.Tokens {
		if token.Token != tok {
			continue
		}
		if token.Location != location {
			continue
		}
		return token
	}
	token := &Token{
		Token:    tok,
		Origin:   tok.Position(),
		Location: location,
	}
	s.Tokens = append(s.Tokens, token)
	return token
}

func (s *Set) Clear() {
	s.Intermediates = s.Intermediates[:0]
	s.Symbols = s.Symbols[:0]
	s.Tokens = s.Tokens[:0]
}
