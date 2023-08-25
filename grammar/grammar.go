package grammar

import (
	"github.com/patrickhuber/go-earley/bitmatrix"
	"github.com/patrickhuber/go-types"
	"github.com/patrickhuber/go-types/handle"
	"github.com/patrickhuber/go-types/result"
)

type Grammar struct {
	Start          NonTerminal
	Productions    []*Production
	Rules          RuleRegistry
	transitiveNull map[Symbol]struct{}
	rightRecursive map[*Production]struct{}
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
	// compute right recursive
	rightRecursive, err := g.identifyRightRecursiveSymbols().Deconstruct()
	if err == nil {
		g.rightRecursive = rightRecursive
	}
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
			changes = 0

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

		sym, ok := rule.PostDotSymbol().Deconstruct()
		if !ok {
			continue
		}

		if _, ok := sym.(NonTerminal); !ok {
			changes++
			continue
		}

		if _, ok := transitiveNull[sym]; ok {
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

func (g *Grammar) identifyRightRecursiveSymbols() (res types.Result[map[*Production]struct{}]) {
	defer handle.Error(&res)

	// if a rule is right recursive, add it to the rules
	// if a rule has a pre dot symbol that is nullable,
	//     walking back in the rule
	//         skip if reaching the end or a terminal
	//         add if reaching a right recursive predot symbol
	rules := []*DottedRule{}
	for _, p := range g.Productions {

		for s := len(p.RightHandSide); s > 0; s-- {
			rule, ok := g.Rules.Get(p, s)
			if !ok {
				break
			}
			sym, ok := rule.PreDotSymbol().Deconstruct()
			if !ok {
				break
			}
			nt, ok := sym.(NonTerminal)
			if !ok {
				break
			}
			// is the rule right recursive?
			if rule.Production.LeftHandSide == nt {
				rules = append(rules, rule)
				break
			}
			// is the rhs nullable?
			if !g.IsTransativeNullable(nt) {
				break
			}
		}
	}

	adjacency := bitmatrix.New(len(rules))
	for row := 0; row < len(rules); row++ {
		left := rules[row]
		for col := 0; col < len(rules); col++ {
			right := rules[col]
			predot, ok := right.PreDotSymbol().Deconstruct()
			if !ok {
				continue
			}
			if left.Production.LeftHandSide == predot {
				adjacency.Set(row, col, true)
			}
		}
	}

	rightRecursive := map[*Production]struct{}{}
	reachability := result.New(
		bitmatrix.TransitiveClosure(adjacency)).Unwrap()

	for row := 0; row < len(rules); row++ {
		reachable := result.New(
			reachability.Get(row, row)).Unwrap()
		if reachable {
			rightRecursive[rules[row].Production] = struct{}{}
		}
	}
	return result.Ok(rightRecursive)
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

func (g *Grammar) IsRightRecursive(p *Production) bool {
	_, ok := g.rightRecursive[p]
	return ok
}
