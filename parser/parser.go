package parser

import (
	"github.com/patrickhuber/go-earley/chart"
	"github.com/patrickhuber/go-earley/grammar"
	"github.com/patrickhuber/go-earley/state"
	"github.com/patrickhuber/go-earley/token"
)

type Parser interface {
	Expected() []grammar.LexerRule
	Accepted() bool
	Location() int
	Pulse(tok token.Token) (bool, error)
}

type parser struct {
	location int
	grammar  *grammar.Grammar
	chart    *chart.Chart
}

func New(g *grammar.Grammar) Parser {
	p := &parser{
		grammar: g,
		chart:   chart.New(),
	}
	p.initialize()
	return p
}

func (p *parser) initialize() {
	p.location = 0
	p.chart = chart.New()
	start := p.grammar.StartProductions()

	for s := 0; s < len(start); s += 1 {
		production := start[s]
		state := p.NewState(production, 0, 0)
		p.chart.Enqueue(0, state)
	}
	p.reductionPass(p.location)
}

func (p *parser) NewState(production *grammar.Production, position int, origin int) state.State {
	dr, ok := p.grammar.Rules.Get(production, position)
	if !ok {
		panic("invalid state")
	}
	return p.chart.GetOrCreate(origin, state.NormalType, dr, origin)
}

func (p *parser) Pulse(tok token.Token) (bool, error) {
	p.scanPass(p.Location(), tok)

	tokenRecognized := len(p.chart.Sets) > p.Location()+1
	if !tokenRecognized {
		return false, nil
	}

	p.location++
	p.reductionPass(p.location)

	return true, nil
}

func (p *parser) scanPass(location int, tok token.Token) {
	set := p.chart.Sets[location]
	for _, s := range set.Scans {
		p.scan(s, location, tok)
	}
}

func (p *parser) scan(s *state.Normal, j int, tok token.Token) {

	sym, ok := s.DottedRule.PostDotSymbol().Deconstruct()
	if !ok {
		return
	}

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
	rule, ok := p.grammar.Rules.Next(s.DottedRule)
	if !ok {
		return
	}

	i := s.Origin
	if p.chart.Contains(j+1, state.NormalType, rule, i) {
		return
	}

	// create a next from the dotted rule
	next := p.NewState(rule.Production, rule.Position, s.Origin)
	p.chart.Enqueue(j+1, next)
}

func (parser *parser) reductionPass(location int) {
	set := parser.chart.Sets[location]
	resume := true

	p := 0
	c := 0

	for resume {
		if c < len(set.Completions) {
			completion := set.Completions[c]
			parser.complete(completion, location)
			c++
			continue
		} else if p < len(set.Predictions) {
			evidence := set.Predictions[p]
			parser.predict(evidence, location)
			p++
			continue
		}
		resume = false
	}
	parser.memoize(location)
}

func (p *parser) complete(completed *state.Normal, location int) {
	set := p.chart.Sets[completed.Origin]
	sym := completed.DottedRule.Production.LeftHandSide

	trans, ok := set.FindTransition(sym)
	if ok {
		p.leoComplete(trans, location)
	} else {
		p.earleyComplete(completed, location)
	}
}

func (p *parser) leoComplete(trans *state.Transition, location int) {
	dr := trans.DottedRule
	origin := trans.Origin

	// don't create memory for something that already exists
	if p.chart.Contains(location, state.NormalType, dr, origin) {
		return
	}

	// jump to the set pointed to by the transition item
	set := p.chart.Sets[trans.Origin]

	// find the top most item (the one that has trans.Sym as its postdot symbol)
	for _, c := range set.Completions {
		sym, ok := c.DottedRule.PostDotSymbol().Deconstruct()
		if !ok {
			continue
		}
		if sym != trans.Symbol {
			continue
		}

		// this is the top most item
		topMostItem := p.NewState(c.DottedRule.Production, c.DottedRule.Position, c.Origin)
		p.chart.Enqueue(location, topMostItem)

		// there will only be one of these
		break
	}
}

func (par *parser) earleyComplete(completed *state.Normal, location int) {

	// get the origin set for the completed state
	completedOrigin := completed.Origin
	set := par.chart.Sets[completedOrigin]

	sources := set.FindSourceStates(completed.DottedRule.Production.LeftHandSide)
	count := len(sources)

	for p := 0; p < count; p++ {
		prediction := sources[p]
		rule, ok := par.grammar.Rules.Next(prediction.DottedRule)
		if !ok {
			continue
		}
		origin := prediction.Origin

		if par.chart.Contains(location, state.NormalType, rule, origin) {
			continue
		}

		state := par.NewState(rule.Production, rule.Position, origin)
		par.chart.Enqueue(location, state)
	}
}

func (parser *parser) memoize(location int) {
	set := parser.chart.Sets[location]

	counts := map[grammar.Symbol]int{}
	states := map[grammar.Symbol]state.State{}

	// count the symbols on the right hand side of each rule
	for _, p := range set.Predictions {
		rule := p.DottedRule
		postDotSymbol, ok := rule.PostDotSymbol().Deconstruct()
		if !ok {
			continue
		}
		_, ok = counts[postDotSymbol]
		if ok {
			counts[postDotSymbol] += 1
		} else {
			states[postDotSymbol] = p
			counts[postDotSymbol] = 1
		}
	}

	// find leo eligible items and memoize them
	for postDot, count := range counts {
		if count != 1 {
			continue
		}
		prediction, ok := states[postDot]
		if !ok {
			continue
		}
		normal, ok := prediction.(*state.Normal)
		if !ok {
			continue
		}
		if !parser.grammar.IsRightRecursive(normal.DottedRule.Production) {
			continue
		}
		next, ok := parser.grammar.Rules.Next(normal.DottedRule)
		if !ok {
			continue
		}
		if !parser.isQuasiComplete(next) {
			continue
		}

		// find the set where this state originated
		set := parser.chart.Sets[normal.Origin]

		// is there a transition?
		trans, ok := set.FindTransition(postDot)

		if ok {
			// if so, copy it here
			clone := &state.Transition{
				Origin:     trans.Origin,
				DottedRule: trans.DottedRule,
				Symbol:     trans.Symbol,
			}
			trans = clone
		} else {
			// otherwise create it
			trans = &state.Transition{
				Origin:     location,
				DottedRule: next,
				Symbol:     postDot,
			}
		}
		parser.chart.Enqueue(location, trans)
	}
}

func (parser *parser) isQuasiComplete(rule *grammar.DottedRule) bool {
	if rule.Complete() {
		return true
	}
	// all postdot symbols are nullable
	for i := rule.Position; i < len(rule.Production.RightHandSide); i++ {
		sym := rule.Production.RightHandSide[i]
		nt, ok := sym.(grammar.NonTerminal)
		if !ok {
			return false
		}
		if !parser.grammar.IsTransativeNullable(nt) {
			return false
		}
		// page 4 leo paper
		// check if S can derive S
		if rule.Production.LeftHandSide == parser.grammar.Start && parser.grammar.Start == nt {
			return false
		}
	}
	return true
}

func (par *parser) predict(evidence *state.Normal, location int) {
	rule := evidence.DottedRule
	sym, ok := rule.PostDotSymbol().Deconstruct()
	if !ok {
		return
	}
	nonTerminal, ok := sym.(grammar.NonTerminal)
	if !ok {
		return
	}
	productions := par.grammar.RulesFor(nonTerminal)

	count := len(productions)
	for p := 0; p < count; p++ {
		production := productions[p]
		par.predictProduction(location, production)
	}

	isNullable := par.grammar.IsTransativeNullable(nonTerminal)
	if isNullable {
		par.predictAycockHorspool(evidence, location)
	}
}

func (p *parser) predictProduction(location int, production *grammar.Production) {
	rule, ok := p.grammar.Rules.Get(production, 0)
	if !ok {
		return
	}
	if p.chart.Contains(location, state.NormalType, rule, location) {
		return
	}
	s := p.NewState(rule.Production, rule.Position, location)
	p.chart.Enqueue(location, s)
}

func (p *parser) predictAycockHorspool(evidence *state.Normal, location int) {
	next, ok := p.grammar.Rules.Next(evidence.DottedRule)
	if !ok {
		return
	}
	if p.chart.Contains(location, evidence.Type(), next, evidence.Origin) {
		return
	}
	state := p.NewState(next.Production, next.Position, evidence.Origin)
	p.chart.Enqueue(location, state)
}

func (p *parser) Location() int {
	return p.location
}

// Accepted implements Parser.
func (p *parser) Accepted() bool {
	_, ok := p.findAcceptedCompletion(p.location)
	return ok
}

func (p *parser) findAcceptedCompletion(location int) (*state.Normal, bool) {
	set := p.chart.Sets[p.location]
	start := p.grammar.Start
	reductions := set.FindReductions(start)
	for c := 0; c < len(reductions); c++ {
		completion := reductions[c]
		if completion.Origin == 0 && completion.DottedRule.Production.LeftHandSide == start {
			return completion, true
		}
	}
	return nil, false
}

// Expected implements Parser.
func (p *parser) Expected() []grammar.LexerRule {
	set := p.chart.Sets[p.location]

	var expected []grammar.LexerRule
	for _, s := range set.Scans {
		postDot, ok := s.DottedRule.PostDotSymbol().Deconstruct()
		if !ok {
			continue
		}
		lexRule, ok := postDot.(grammar.LexerRule)
		if !ok {
			continue
		}
		expected = append(expected, lexRule)
	}
	return expected
}
