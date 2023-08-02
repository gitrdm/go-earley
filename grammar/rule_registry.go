package grammar

type RuleRegistry interface {
	Register(*DottedRule)
	Next(*DottedRule) (*DottedRule, bool)
	Get(production *Production, position int) (*DottedRule, bool)
}

type registry struct {
	productionRules map[*Production]map[int]*DottedRule
}

func NewRegistry() RuleRegistry {
	return &registry{
		productionRules: make(map[*Production]map[int]*DottedRule),
	}
}

func (r *registry) Register(rule *DottedRule) {
	rules, ok := r.productionRules[rule.Production]
	if !ok {
		rules = map[int]*DottedRule{}
		r.productionRules[rule.Production] = rules

	}
	rules[rule.Position] = rule
}

// Get implements RuleRegistry.
func (r *registry) Get(production *Production, position int) (*DottedRule, bool) {
	dottedRules, ok := r.productionRules[production]
	if !ok {
		return nil, false
	}
	dr, ok := dottedRules[position]
	return dr, ok
}

// Next implements RuleRegistry.
func (r *registry) Next(dr *DottedRule) (*DottedRule, bool) {
	rules, ok := r.productionRules[dr.Production]
	if !ok {
		return nil, false
	}
	next, ok := rules[dr.Position+1]
	if !ok {
		return nil, false
	}
	return next, true
}
