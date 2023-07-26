package parser

import (
	"github.com/patrickhuber/go-earley/chart"
	"github.com/patrickhuber/go-earley/grammar"
	"github.com/patrickhuber/go-earley/token"
)

type Parser interface {
	Expected() []grammar.LexerRule
	Accepted() bool
	Location() int
	Pulse(tok token.Token) (bool, error)
}

type parser struct {
	location    int
	dottedRules grammar.RuleRegistry
	grammar     grammar.Grammar
	chart       chart.Chart
}

func New(g grammar.Grammar) Parser {
	return &parser{
		grammar: g,
	}
}

func (p *parser) Pulse(tok token.Token) (bool, error) {
	p.scanPass(p.Location(), tok)

	tokenRecognized := len(p.chart.Sets()) > p.Location()+1
	if !tokenRecognized {
		return false, nil
	}

	p.location++
	p.reductionPass(p.location)

	return true, nil
}

func (p *parser) scanPass(location int, tok token.Token) {
	set := p.chart.Sets()[location]
	for _, s := range set.Scans() {
		p.scan(s, location, tok)
	}
}

func (p *parser) scan(state chart.State, j int, tok token.Token) {

	sym := grammar.PostDotSymbol(state.DottedRule())

	// process lexer rules
	lexRule, ok := sym.(grammar.LexerRule)
	if !ok {
		return
	}

	// skip scanning if the token type doesn't match
	if lexRule.Type() != tok.Type() {
		return
	}

	// grab the next dotted rule from the registry
	rule := p.dottedRules.Next(state.DottedRule())
	i := state.Origin()
	if p.chart.Contains(j+1, chart.NormalStateType, rule, i) {
		return
	}
}

func (p *parser) reductionPass(location int) {

}

func (p *parser) Location() int {
	return p.location
}

// Accepted implements Parser.
func (*parser) Accepted() bool {
	panic("unimplemented")
}

// Expected implements Parser.
func (*parser) Expected() []grammar.LexerRule {
	panic("unimplemented")
}
