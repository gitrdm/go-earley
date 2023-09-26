package parser

import (
	"fmt"

	"github.com/patrickhuber/go-earley/forest"
	"github.com/patrickhuber/go-earley/grammar"
	"github.com/patrickhuber/go-earley/internal/chart"
	"github.com/patrickhuber/go-earley/internal/state"
	"github.com/patrickhuber/go-earley/token"
)

type Parser interface {
	Expected() []grammar.LexerRule
	Accepted() bool
	Location() int
	Pulse(tok token.Token) (bool, error)
	GetForestRoot() (forest.Node, bool)
}

type parser struct {
	location               int
	grammar                *grammar.Grammar
	chart                  *chart.Chart
	nodes                  *forest.Set
	optimizeRightRecursion bool
}

type Option func(*parser)

// OptimizeRightRecursion optimizes right recursion
// the default is true
func OptimizeRightRecursion(ok bool) Option {
	return func(p *parser) {
		p.optimizeRightRecursion = ok
	}
}

func New(g *grammar.Grammar, options ...Option) Parser {
	p := &parser{
		grammar:                g,
		chart:                  chart.New(),
		nodes:                  &forest.Set{},
		optimizeRightRecursion: true,
	}
	for _, option := range options {
		option(p)
	}
	p.initialize()
	return p
}

func (p *parser) initialize() {
	fmt.Printf("--------- %d ---------", p.location)
	fmt.Println()
	p.location = 0
	p.chart = chart.New()
	start := p.grammar.StartProductions()

	for s := 0; s < len(start); s += 1 {
		production := start[s]
		state := p.newState(production, 0, 0)
		p.chart.Enqueue(0, state)
		fmt.Printf("%s : Init", state)
		fmt.Println()
	}
	p.reductionPass(p.location)
}

func (p *parser) newState(production *grammar.Production, position int, origin int) *state.Normal {
	rule, ok := p.grammar.Rules.Get(production, position)
	if !ok {
		panic("invalid state")
	}
	return state.NewNormal(rule, origin)
}

func (p *parser) Pulse(tok token.Token) (bool, error) {
	fmt.Printf("--------- %d ---------", p.location+1)
	fmt.Println()
	p.scanPass(p.Location(), tok)

	tokenRecognized := len(p.chart.Sets) > p.Location()+1
	if !tokenRecognized {
		return false, nil
	}

	p.location++

	p.reductionPass(p.location)
	p.nodes.Clear()

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
	if lexRule.TokenType() != tok.Type() {
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

	// create the parse node
	tokenNode := p.nodes.AddOrGetExistingTokenNode(tok, j+1)
	parseNode := p.createParseNode(rule, s.Origin, s.Node, tokenNode, j+1)

	// create a next from the dotted rule
	next := p.newState(rule.Production, rule.Position, s.Origin)
	next.Node = parseNode
	p.chart.Enqueue(j+1, next)
	fmt.Printf("%s : Scan", next)
	fmt.Println()
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
		} else if p < len(set.Predictions) {
			evidence := set.Predictions[p]
			parser.predict(evidence, location)
			p++
		} else {
			resume = false
		}
	}
	if parser.optimizeRightRecursion {
		parser.memoize(location)
	}
}

func (p *parser) complete(completed *state.Normal, location int) {
	set := p.chart.Sets[completed.Origin]
	sym := completed.DottedRule.Production.LeftHandSide

	if completed.Node == nil {
		completed.Node = p.nodes.AddOrGetExistingSymbolNode(
			completed.DottedRule.Production.LeftHandSide,
			completed.Origin,
			location)
	}

	trans, ok := set.FindTransition(sym)
	if ok {
		p.leoComplete(completed, trans, location)
	} else {
		p.earleyComplete(completed, location)
	}
}

func (p *parser) leoComplete(completed *state.Normal, trans *state.Transition, location int) {

	dottedRule := trans.DottedRule
	origin := trans.Origin

	// check if the item exists
	if p.chart.Contains(location, state.NormalType, dottedRule, origin) {
		return
	}

	// this is the top most item
	top := p.newState(dottedRule.Production, dottedRule.Position, origin)
	node := p.nodes.AddOrGetExistingSymbolNode(dottedRule.Production.LeftHandSide, origin, location)
	top.Node = node

	root, ok := p.chart.Sets[trans.Root].FindTransition(trans.Symbol)
	if !ok {
		root = trans
	}
	node.AddPath(root, completed.Node)

	p.chart.Enqueue(location, top)

	fmt.Printf("%s : Leo Complete", top)
	fmt.Println()
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

		// create a parse node before the existence check
		// this is done on purpose
		node := par.createParseNode(rule, origin, prediction.Node, completed.Node, location)

		if par.chart.Contains(location, state.NormalType, rule, origin) {
			continue
		}

		state := par.newState(rule.Production, rule.Position, origin)
		state.Node = node

		par.chart.Enqueue(location, state)

		fmt.Printf("%s : Earley Complete", state)
		fmt.Println()
	}
}

// memoize implements the memoization algorithm in the marpa paper
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

		// if postdot is terminal, skip
		_, ok = postDotSymbol.(grammar.NonTerminal)
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

		// create the transition
		trans, ok := parser.newTransition(normal, normal.Origin, location)
		if !ok {
			continue
		}

		// add the transition
		parser.chart.Enqueue(location, trans)

		// log the transistion creation
		fmt.Printf("%s : Transition", trans.String())
		fmt.Println()
	}
}

func (parser *parser) newTransition(predict *state.Normal, origin int, location int) (*state.Transition, bool) {
	sym, ok := predict.DottedRule.PostDotSymbol().Deconstruct()
	if !ok {
		return nil, ok
	}
	trans, ok := parser.chart.Sets[origin].FindTransition(sym)
	if ok {
		// if so, copy it here
		clone := &state.Transition{
			DottedRule: trans.DottedRule,
			Origin:     trans.Origin,
			Symbol:     trans.Symbol,
			Predict:    predict.Node,
			Root:       trans.Root,
		}
		trans.SetNext(clone)
		trans = clone

	} else {
		next, ok := parser.grammar.Rules.Next(predict.DottedRule)
		if !ok {
			return nil, ok
		}
		// otherwise create it
		trans = &state.Transition{
			DottedRule: next,
			Origin:     predict.Origin,
			Symbol:     sym,
			Predict:    predict.Node,
			Root:       location,
		}
	}
	return trans, true
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
		par.predictAycockHorspool(evidence, nonTerminal, location)
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
	s := p.newState(rule.Production, rule.Position, location)
	p.chart.Enqueue(location, s)
	fmt.Printf("%s : Predict", s)
	fmt.Println()
}

func (p *parser) predictAycockHorspool(evidence *state.Normal, nullableSymbol grammar.Symbol, location int) {
	next, ok := p.grammar.Rules.Next(evidence.DottedRule)
	if !ok {
		return
	}
	if p.chart.Contains(location, evidence.Type(), next, evidence.Origin) {
		return
	}
	state := p.newState(next.Production, next.Position, evidence.Origin)

	emptyNode := p.nodes.AddOrGetExistingSymbolNode(nullableSymbol, location, location)

	// create the node for the completed item
	node := p.createParseNode(next, evidence.Origin, evidence.Node, emptyNode, location)
	state.Node = node

	p.chart.Enqueue(location, state)
	fmt.Printf("%s : Predict AH", state)
	fmt.Println()
}

func (p *parser) Location() int {
	return p.location
}

// Accepted implements Parser.
func (p *parser) Accepted() bool {
	_, ok := p.findAcceptedCompletion(p.location)
	return ok
}

func (p *parser) GetForestRoot() (forest.Node, bool) {
	s, ok := p.findAcceptedCompletion(p.location)
	if !ok {
		return nil, false
	}
	return s.Node, true
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

func (p *parser) createParseNode(
	rule *grammar.DottedRule,
	origin int,
	w,
	v forest.Node,
	location int) forest.Node {

	/*
		B -> ax*D, j, i, w, v, V
		if D is nil {
			let s = B
		}
		else {
			let s = (B -> ax*D)
		}
		if a is nil and D is not nil{
			let y = v
		}
		else {
			if no node labeled (s,j,i) in V, create one add to V
			if w == null and y has no family of children (v), add one
			if w != null and y has no family of children (w,v), add one
		}
		return y
	*/
	var internal forest.Internal
	var node forest.Node

	if rule.Complete() {
		symbol := p.nodes.AddOrGetExistingSymbolNode(
			rule.Production.LeftHandSide,
			origin,
			location,
		)
		node = symbol
		internal = symbol
	} else {
		intermediate := p.nodes.AddOrGetExistingIntermediateNode(
			rule,
			origin,
			location,
		)
		node = intermediate
		internal = intermediate
	}

	// this will merge parent and child nodes
	// it is placed after to support caching of the original node
	if !rule.Complete() &&
		rule.Position == 1 &&
		v != nil {
		return v
	}

	if w == nil && v != nil {
		internal.AddUniqueFamily(v, nil)
	} else if w != nil && v != nil {
		internal.AddUniqueFamily(w, v)
	}
	return node
}
