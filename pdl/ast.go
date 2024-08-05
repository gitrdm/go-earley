package pdl

import "github.com/patrickhuber/go-earley/re"

type Definition interface {
	definition()
}

type DefinitionBlock struct {
	Block Block
}

func (DefinitionBlock) definition() {}

type DefinitionBlockDefinition struct {
	Block      Block
	Definition Definition
}

func (DefinitionBlockDefinition) definition() {}

type Block interface {
	block()
}

type Rule struct {
	QualifiedIdentifier QualifiedIdentifier
	Expression          Expression
}

func (Rule) block() {}

type Setting struct {
	SettingIdentifier   SettingIdentifier
	QualifiedIdentifier QualifiedIdentifier
}

func (Setting) block() {}

type LexerRule struct {
	QualifiedIdentifier QualifiedIdentifier
	LexerRuleExpression LexerRuleExpression
}

func (LexerRule) block() {}

type Expression interface {
	expression()
}

type ExpressionTerm struct {
	Term Term
}

type ExpressionTermExpression struct {
	Term       Term
	Expression Expression
}

type Term interface {
	term()
}

type TermFactor struct {
	Factor Factor
}

type TermFactorTerm struct {
	Factor Factor
	Term   Term
}

type Factor interface {
	factor()
}

type Literal interface {
	literal()
	factor()
}

type Repetition struct {
	Expression Expression
}

func (Repetition) factor() {}

type Optional struct {
	Expression Expression
}

func (Optional) factor() {}

type Grouping struct {
	Expression Expression
}

func (Grouping) factor() {}

type QualifiedIdentifier interface {
	qualifiedIdentifier()
	factor()
}

type QualifiedIdentifierIdentifier struct{}

func (QualifiedIdentifierIdentifier) qualifiedIdentifier() {}

func (QualifiedIdentifierIdentifier) factor() {}

type QualifiedIdentifierIdentifierLetter struct{}

func (QualifiedIdentifierIdentifierLetter) qualifiedIdentifier() {}

func (QualifiedIdentifierIdentifierLetter) factor() {}

type SettingIdentifier struct {
}

func (SettingIdentifier) factor() {}

type LexerRuleExpression interface {
	lexerRuleExpression()
}

type LexerRuleExpressionTerm struct {
	LexerRuleTerm LexerRuleTerm
}

func (LexerRuleExpressionTerm) lexerRuleExpression() {}

type LexerRuleExpressionTermExpression struct {
	LexerRuleTerm       LexerRuleTerm
	LexerRuleExpression LexerRuleExpression
}

func (LexerRuleExpressionTermExpression) lexerRuleExpression() {}

type LexerRuleTerm interface {
	lexerRuleTerm()
}

type LexerRuleTermFactor struct{}

func (LexerRuleTermFactor) lexerRuleTerm() {}

type LexerRuleTermFactorTerm struct{}

func (LexerRuleTermFactorTerm) lexerRuleTerm() {}

type RegularExpression struct {
	Definition re.Definition
}
