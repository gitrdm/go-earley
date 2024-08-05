package grammar

const (
	StringLexerRuleType = "string"
)

type StringLexerRule struct {
	Value string
	SymbolImpl
}

// CanApply implements String.
func (s *StringLexerRule) CanApply(ch rune) bool {
	for _, r := range s.Value {
		return r == ch
	}
	// empty case
	return false
}

func (t *StringLexerRule) TokenType() string {
	return t.Value
}

// ID implements String.
func (s *StringLexerRule) LexerRuleType() string {
	return StringLexerRuleType
}

func NewStringLexerRule(str string) *StringLexerRule {
	return &StringLexerRule{
		Value: str,
	}
}

func (s *StringLexerRule) String() string {
	return s.Value
}
