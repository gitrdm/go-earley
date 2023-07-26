package grammar

type RuleRegistry interface {
	Register(*DottedRule)
	Next(*DottedRule) *DottedRule
}
