package grammar

type Grammar struct {
	Start          NonTerminal
	Productions    []*Production
	Rules          RuleRegistry
	transitiveNull map[Symbol]struct{}
}

func New(start NonTerminal, productions ...*Production) *Grammar {
	g := &Grammar{
		Start:       start,
		Productions: productions,
	}
	// compute dotted rules registry
	g.Rules = compute(g)
	// compute transitive null
	g.transitiveNull = identifyNullableSymbols(g)

	return g
}

func compute(g *Grammar) RuleRegistry {
	r := NewRegistry()
	for p := range g.Productions {
		production := g.Productions[p]

		// this needs to be len(rhs)+1 because dots are between characters
		for i := 0; i <= len(production.RightHandSide); i++ {
			dr := &DottedRule{
				Production: production,
				Position:   i,
			}
			r.Register(dr)
		}
	}
	return r
}

func identifyNullableSymbols(g *Grammar) map[Symbol]struct{} {

	transitiveNull := make(map[Symbol]struct{})

	work := []*DottedRule{}
	unprocessed := []*DottedRule{}

	for _, p := range g.Productions {
		if len(p.RightHandSide) == 0 {
			transitiveNull[p.LeftHandSide] = struct{}{}
		}
		rule, ok := g.Rules.Get(p, 0)
		if !ok {
			continue
		}
		work = append(work, rule)
	}

	changes := 0

	for len(work) > 0 || len(unprocessed) > 0 {
		if len(work) == 0 {
			if changes == 0 {
				break
			}
			temp := unprocessed
			unprocessed = work
			work = temp
		}
		var rule *DottedRule
		work, rule = dequeue(work)

		if _, ok := transitiveNull[rule.Production.LeftHandSide]; ok {
			changes++
			continue
		}

		if rule.Complete() {
			transitiveNull[rule.Production.LeftHandSide] = struct{}{}
			changes++
			continue
		}

		if _, ok := rule.PostDotSymbol().(NonTerminal); !ok {
			changes++
			continue
		}

		if _, ok := transitiveNull[rule.PostDotSymbol()]; ok {
			next, ok := g.Rules.Next(rule)
			if !ok {
				continue
			}
			for _, u := range unprocessed {
				if u == next {
					continue
				}
			}
			unprocessed = enqueue(unprocessed, next)
			changes++
			continue
		}
		unprocessed = enqueue(unprocessed, rule)
	}
	return transitiveNull
}

func enqueue[T any](queue []T, item T) []T {
	return append(queue, item)
}

func dequeue[T any](queue []T) ([]T, T) {
	return queue[1:], queue[0]
}

func (g *Grammar) RulesFor(nt NonTerminal) []*Production {
	// TODO: optimize this if it becomes a memory hog
	var productions []*Production
	for p := range g.Productions {
		production := g.Productions[p]
		if production.LeftHandSide == nt {
			productions = append(productions, production)
		}
	}
	return productions
}

func (g *Grammar) StartProductions() []*Production {
	return g.RulesFor(g.Start)
}

func (g *Grammar) IsTransativeNullable(nt NonTerminal) bool {
	_, ok := g.transitiveNull[nt]
	return ok
}
