package grammar

type LexerRule interface {
	Symbol
	CanApply(ch rune) bool
	Type() string
	TokenType() string
}
