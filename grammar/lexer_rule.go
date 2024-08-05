package grammar

type LexerRule interface {
	Symbol
	CanApply(ch rune) bool
	LexerRuleType() string
	TokenType() string
}
