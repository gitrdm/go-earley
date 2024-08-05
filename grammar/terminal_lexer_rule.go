package grammar

const (
	TerminalLexerRuleType = "terminal"
)

type TerminalLexerRule struct {
	Terminal Terminal
	SymbolImpl
}

// CanApply implements Terminal.
func (t *TerminalLexerRule) CanApply(ch rune) bool {
	return t.Terminal.IsMatch(ch)
}

// ID implements Terminal.
func (t *TerminalLexerRule) LexerRuleType() string {
	return TerminalLexerRuleType
}

func (t *TerminalLexerRule) TokenType() string {
	return t.Terminal.String()
}

func NewTerminalLexerRule(t Terminal) LexerRule {
	return &TerminalLexerRule{
		Terminal: t,
	}
}

func (t *TerminalLexerRule) String() string {
	return t.Terminal.String()
}
