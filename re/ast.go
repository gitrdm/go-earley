package re

type Definition struct {
	Start      bool
	Expression Expression
	End        bool
}

type Expression interface {
	expression()
}

type ExpressionTerm struct {
	Term Term
}

func (ExpressionTerm) expression() {}

type ExpressionTermExpression struct {
	Term       Term
	Expression Expression
}

func (ExpressionTermExpression) expression() {}

type Term interface {
	term()
}

type TermFactor struct {
	Factor Factor
}

func (TermFactor) term() {}

type TermFactorTerm struct {
	Fator Factor
	Term  Term
}

func (TermFactorTerm) term() {}

type Factor interface {
	factor()
}

type FactorAtom struct {
	Atom Atom
}

func (FactorAtom) factor() {}

type FactorAtomIterator struct {
	Atom     Atom
	Iterator Iterator
}

func (FactorAtomIterator) factor() {}

type Iterator string

const (
	ZeroOrOne  Iterator = "?"
	OneOrMany  Iterator = "+"
	ZeroOrMany Iterator = "*"
)

type Atom interface {
	atom()
}

type AtomAny struct{}

func (AtomAny) atom() {}

type AtomCharacter struct{}

func (AtomCharacter) atom() {}

type AtomExpression struct{}

func (AtomExpression) atom() {}

type AtomSet struct{}

func (AtomSet) atom() {}

type Set interface {
	set()
}
type NegativeSet struct {
	Characterclass CharacterClass
}

func (NegativeSet) set() {}

type PositiveSet struct {
	CharacterClass CharacterClass
}

func (PositiveSet) set() {}

type CharacterClass interface {
	characterClass()
}

type CharacterRange interface {
	characterRange()
}

type CharacterRangeCharacterClassCharacter struct {
	Begin CharacterClassCharacter
}

type CharacterRangeCharacterClassCharacterRange struct {
	Begin CharacterClassCharacter
	End   CharacterClassCharacter
}

type Character interface {
	character()
}

type CharacterClassCharacter interface {
	characterClassCharacter()
}

// NotMetaCharacter represents a non meta character .^$()[]+*?\/
// /[^.^$()[\]+*?\\\/]/;
type NotMetaCharacter struct {
	Char rune
}

func (NotMetaCharacter) character()               {}
func (NotMetaCharacter) characterClassCharacter() {}

// EscapeSequence is an backslash followed by any character
type EscapeSequence struct {
	Char rune
}

func (EscapeSequence) character()               {}
func (EscapeSequence) characterClassCharacter() {}

// NotCloseBracketCharacter is a non close bracket ] character
type NotCloseBracketCharacter struct {
	Char rune
}

func (NotCloseBracketCharacter) characterClassCharacter() {}
