package re

import "github.com/patrickhuber/go-types"

type Expression struct {
	Term       types.Option[Term]
	Expression types.Option[Expression]
}

type Term struct {
	Term   types.Option[Term]
	Factor Factor
}

type Factor struct {
	Atom     Atom
	Iterator types.Option[Iterator]
}

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

func (AtomAny) atom()

type AtomCharacter struct{}

func (AtomCharacter) atom()

type AtomExpression struct{}

func (AtomExpression) atom()

type AtomSet struct{}

func (AtomSet) atom()

type Set struct {
	Negate         bool
	CharacterClass CharacterClass
}

type CharacterClass interface {
	characterclass()
}

type CharacterClassRange struct {
	Range Range
}

func (CharacterClassRange) characterclass() {}

type CharacterClassAlteration struct {
	Range      Range
	Alteration CharacterClass
}

func (CharacterClassAlteration) characterclass() {}

type Range struct{}
